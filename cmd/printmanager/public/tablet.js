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
}

function status(text, bad) {
  const s = $("status");
  s.textContent = text;
  s.classList.toggle("bad", !!bad);
}

async function enroll(ev) {
  ev.preventDefault();
  $("enroll-error").textContent = "";
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

function time(s) {
  const d = new Date(s);
  return d.toLocaleString("ja-JP", { month: "numeric", day: "numeric", hour: "2-digit", minute: "2-digit" });
}

function item(s) {
  const li = document.createElement("li");
  li.className = "slip " + s.state;
  const name = document.createElement("span");
  name.className = "name";
  name.textContent = s.name || "伝票 " + s.job;
  const meta = document.createElement("span");
  meta.className = "meta";
  meta.textContent = `#${s.job} · ${time(s.received)} · ${s.pages} ページ`;
  const badge = document.createElement("span");
  badge.className = "badge " + s.state;
  badge.textContent = stateText[s.state] || s.state;
  li.append(badge, name, meta);
  if (s.printed && s.printed.length) {
    const p = document.createElement("span");
    p.className = "badge printed";
    p.textContent = "紙に出した";
    li.append(p);
  }
  if (s.state === "new" || s.state === "opened") li.addEventListener("click", () => openSlip(s));
  return li;
}

function render() {
  $("title").textContent = board.title;
  $("open-count").textContent = board.open.length ? `(${board.open.length})` : "";
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
      if (!first && board.open.some((s) => s.state === "new" && !known.has(s.id))) beep();
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
  } catch (e) {
    status(e.message, true);
  }
});
$("print").addEventListener("click", async () => {
  try {
    await post(`/device/slips/${encodeURIComponent(current.id)}/print`);
    status("刷りました");
  } catch (e) {
    status("刷れませんでした: " + e.message, true);
  }
});
for (const t of ["open", "ended"]) {
  $("tab-" + t).addEventListener("click", () => {
    tab = t;
    $("tab-open").classList.toggle("active", t === "open");
    $("tab-ended").classList.toggle("active", t === "ended");
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
