// The page of a device (a tablet of a location). It enrolls the device
// with a one-time code: it makes an ECDSA P-256 key that cannot be
// exported (WebCrypto, extractable false), keeps it in IndexedDB, and
// sends the public key. Afterwards it signs the manager's challenge for a
// session, polls the location's board, and shows the slips' pages.

const POLL_MS = 4000;
const $ = (id) => document.getElementById(id);

// The device's registration, kept in IndexedDB: the CryptoKey objects
// are stored as they are (the private key never leaves the browser).
const db = {
  open() {
    return new Promise((resolve, reject) => {
      const req = indexedDB.open("bdf-print-manager", 1);
      req.onupgradeneeded = () => req.result.createObjectStore("device");
      req.onsuccess = () => resolve(req.result);
      req.onerror = () => reject(req.error);
    });
  },
  async run(mode, fn) {
    const d = await db.open();
    return new Promise((resolve, reject) => {
      const tx = d.transaction("device", mode);
      const req = fn(tx.objectStore("device"));
      tx.oncomplete = () => resolve(req.result);
      tx.onerror = () => reject(tx.error);
    });
  },
  get: () => db.run("readonly", (s) => s.get("me")),
  put: (v) => db.run("readwrite", (s) => s.put(v, "me")),
  clear: () => db.run("readwrite", (s) => s.delete("me")),
};

const b64 = (buf) => btoa(String.fromCharCode(...new Uint8Array(buf)));

let me = null; // {device, location, privateKey}
let info = null; // /device/me
let etag = "";
let board = null;
let current = null; // the slip open in the viewer
let known = new Set(); // ids of the open slips seen, to tell new ones
let first = true;
let tab = "open";

function show(section) {
  for (const id of ["enroll", "board", "viewer"]) $(id).hidden = id !== section;
  // while a slip is open, only its bar stays at the top
  $("bar").classList.toggle("sticky", section !== "viewer");
}

function status(text, bad) {
  const s = $("status");
  s.textContent = text;
  s.hidden = !text;
  s.className = "badge " + (bad ? "badge-error" : "badge-success badge-soft");
}

// toast shows a message for a while; big is for a slip that came.
function toast(text, kind = "info", big = false) {
  const t = document.createElement("div");
  t.role = "status";
  // whole class names, for Tailwind to find
  const alert = { info: "alert-info", success: "alert-success", warning: "alert-warning", error: "alert-error" }[kind];
  t.className = `alert ${alert} shadow-lg ` + (big ? "px-8 py-6 text-2xl font-bold sm:text-3xl" : "text-lg");
  t.textContent = text;
  $("toasts").append(t);
  setTimeout(() => t.remove(), big ? 8000 : 4000);
}

async function enroll(ev) {
  ev.preventDefault();
  $("enroll-error").hidden = true;
  try {
    const keys = await crypto.subtle.generateKey({ name: "ECDSA", namedCurve: "P-256" }, false, ["sign", "verify"]);
    const spki = await crypto.subtle.exportKey("spki", keys.publicKey);
    const res = await fetch("/device/enroll", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ code: $("code").value, name: $("device-name").value, key: b64(spki) }),
    });
    const body = await res.json();
    if (!res.ok) throw new Error(body.error || res.statusText);
    me = { device: body.device, location: body.location, privateKey: keys.privateKey };
    await db.put(me);
    history.replaceState(null, "", location.pathname);
    await start();
  } catch (e) {
    $("enroll-error").textContent = "登録できませんでした: " + e.message;
    $("enroll-error").hidden = false;
  }
}

// login signs a challenge for a session. A device the manager no longer
// knows is enrolled again.
async function login() {
  const ch = await (await fetch("/device/challenge")).json();
  const sig = await crypto.subtle.sign({ name: "ECDSA", hash: "SHA-256" }, me.privateKey, new TextEncoder().encode(ch.challenge));
  const res = await fetch("/device/session", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ device: me.device, challenge: ch.challenge, signature: b64(sig) }),
  });
  if (res.status === 401) {
    await db.clear();
    me = null;
    show("enroll");
    throw new Error("この端末は登録されていません");
  }
  if (!res.ok) throw new Error("session " + res.status);
}

// api fetches with the session, signing in again when it ended.
async function api(path, opts = {}) {
  let res = await fetch(path, opts);
  if (res.status === 401 && me) {
    await login();
    res = await fetch(path, opts);
  }
  return res;
}

async function post(path) {
  const res = await api(path, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ worker: $("worker").value }),
  });
  const body = await res.json();
  if (!res.ok) throw new Error(body.error || res.statusText);
  etag = ""; // read the board again
  return body;
}

function beep() {
  try {
    const ctx = new AudioContext();
    const o = ctx.createOscillator();
    o.frequency.value = 880;
    o.connect(ctx.destination);
    o.start();
    o.stop(ctx.currentTime + 0.2);
  } catch {}
}

const stateText = { new: "新着", opened: "開いた", done: "完了", canceled: "取り消し" };
const stateBadge = { new: "badge-warning", opened: "badge-info", done: "badge-success", canceled: "badge-error" };

function time(s) {
  const d = new Date(s);
  return d.toLocaleString("ja-JP", { month: "numeric", day: "numeric", hour: "2-digit", minute: "2-digit" });
}

function item(s) {
  const active = s.state === "new" || s.state === "opened";
  const li = document.createElement("li");
  li.className =
    "card bg-base-100 shadow-sm transition " +
    (s.state === "new" ? "border-l-8 border-warning " : "") +
    (active ? "cursor-pointer hover:shadow-md active:scale-[0.99]" : "opacity-70");
  const body = document.createElement("div");
  body.className = "card-body flex-row flex-wrap items-center gap-x-4 gap-y-1 p-5";
  const badge = document.createElement("span");
  badge.className = "badge badge-lg " + (stateBadge[s.state] || "badge-ghost");
  badge.textContent = stateText[s.state] || s.state;
  const name = document.createElement("span");
  name.className = "min-w-0 flex-1 break-words text-xl font-semibold sm:text-2xl";
  name.textContent = s.name || "伝票 " + s.job;
  const meta = document.createElement("span");
  meta.className = "w-full text-base text-base-content/60";
  meta.textContent = `#${s.job} · ${time(s.received)} · ${s.pages} ページ`;
  body.append(badge, name);
  if (s.printed && s.printed.length) {
    const p = document.createElement("span");
    p.className = "badge badge-outline badge-lg";
    p.textContent = "紙に出した";
    body.append(p);
  }
  body.append(meta);
  li.append(body);
  if (active) li.addEventListener("click", () => openSlip(s));
  return li;
}

function render() {
  $("title").textContent = board.title;
  $("open-count").textContent = board.open.length;
  $("open-count").hidden = board.open.length === 0;
  $("open-list").replaceChildren(...board.open.map(item));
  $("ended-list").replaceChildren(...board.ended.map(item));
  $("open-list").hidden = tab !== "open";
  $("ended-list").hidden = tab !== "ended";
  $("empty").hidden = (tab === "open" ? board.open : board.ended).length > 0;
}

async function poll() {
  if (!me) return;
  try {
    const res = await api("/device/board", { headers: etag ? { "If-None-Match": etag } : {} });
    if (res.status === 200) {
      etag = res.headers.get("ETag") || "";
      board = await res.json();
      // a slip that came since the last board rings, not those on the first
      const fresh = board.open.filter((s) => s.state === "new" && !known.has(s.id));
      if (!first && fresh.length) {
        beep();
        toast(fresh.length === 1 ? `新しい伝票: ${fresh[0].name || "#" + fresh[0].job}` : `新しい伝票が ${fresh.length} 件`, "warning", true);
      }
      first = false;
      known = new Set(board.open.map((s) => s.id));
      if (!current) render();
    } else if (res.status !== 304) {
      throw new Error("board " + res.status);
    }
    status("接続中 · " + new Date().toLocaleTimeString("ja-JP", { hour: "2-digit", minute: "2-digit" }));
  } catch (e) {
    status("つながっていません: " + e.message, true);
  }
}

async function openSlip(s) {
  current = s;
  show("viewer");
  $("viewer-title").textContent = `${s.name || "伝票 " + s.job} (#${s.job})`;
  $("print").hidden = !(info && info.can_print);
  const width = Math.min(2400, Math.round(window.innerWidth * (window.devicePixelRatio || 1)));
  const pages = [];
  for (let n = 1; n <= s.pages; n++) {
    const img = document.createElement("img");
    img.alt = `${n} ページ`;
    img.className = "w-full rounded-box bg-white shadow-md";
    img.loading = n <= 2 ? "eager" : "lazy";
    img.src = `/device/slips/${encodeURIComponent(s.id)}/pages/${n}.png?w=${width}`;
    pages.push(img);
  }
  $("pages").replaceChildren(...pages);
  if (s.state === "new") {
    try {
      await post(`/device/slips/${encodeURIComponent(s.id)}/open`);
    } catch (e) {
      status(e.message, true);
    }
  }
}

function closeSlip() {
  current = null;
  $("pages").replaceChildren();
  show("board");
  if (board) render();
  poll();
}

async function start() {
  me = me || (await db.get());
  const code = new URLSearchParams(location.search).get("code");
  if (!me) {
    if (code) $("code").value = code;
    show("enroll");
    return;
  }
  try {
    let res = await api("/device/me");
    if (res.status === 401) {
      await login();
      res = await api("/device/me");
    }
    info = await res.json();
  } catch (e) {
    status(e.message, true);
    if (!me) return;
  }
  show("board");
  await poll();
}

$("enroll-form").addEventListener("submit", enroll);
$("back").addEventListener("click", closeSlip);
$("done").addEventListener("click", async () => {
  try {
    await post(`/device/slips/${encodeURIComponent(current.id)}/done`);
    closeSlip();
    toast("完了にしました", "success");
  } catch (e) {
    toast(e.message, "error");
  }
});
$("print").addEventListener("click", async () => {
  try {
    await post(`/device/slips/${encodeURIComponent(current.id)}/print`);
    toast("刷りました", "success");
  } catch (e) {
    toast("刷れませんでした: " + e.message, "error");
  }
});
for (const t of ["open", "ended"]) {
  $("tab-" + t).addEventListener("click", () => {
    tab = t;
    $("tab-open").classList.toggle("tab-active", t === "open");
    $("tab-ended").classList.toggle("tab-active", t === "ended");
    if (board) render();
  });
}
const worker = $("worker");
try {
  worker.value = localStorage.getItem("bdf-pm-worker") || "";
} catch {}
worker.addEventListener("change", () => {
  try {
    localStorage.setItem("bdf-pm-worker", worker.value);
  } catch {}
});

setInterval(() => {
  if (!document.hidden) poll();
}, POLL_MS);
document.addEventListener("visibilitychange", () => {
  if (!document.hidden) poll();
});
start();
