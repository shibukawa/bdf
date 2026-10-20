// The administrator's page: printers, locations, devices and their
// enrollment codes, through the API with the admin token (kept in this
// tab's sessionStorage only).

const $ = (id) => document.getElementById(id);
let token = "";
try {
  token = sessionStorage.getItem("bdf-pm-admin") || "";
} catch {}
$("token").value = token;

function status(text, bad) {
  $("admin-status").textContent = text;
  $("admin-status").classList.toggle("bad", !!bad);
}

async function api(method, path, body) {
  const res = await fetch(path, {
    method,
    headers: { Authorization: "Bearer " + token, ...(body ? { "Content-Type": "application/json" } : {}) },
    body: body ? JSON.stringify(body) : undefined,
  });
  if (res.status === 204) return null;
  const out = await res.json();
  if (!res.ok) throw new Error(out.error || res.statusText);
  return out;
}

function el(tag, text, attrs = {}) {
  const e = document.createElement(tag);
  if (text !== undefined) e.textContent = text;
  Object.assign(e, attrs);
  return e;
}

function button(text, fn) {
  const b = el("button", text, { type: "button" });
  b.addEventListener("click", async () => {
    try {
      await fn();
      await load();
    } catch (e) {
      status(e.message, true);
    }
  });
  return b;
}

function since(t) {
  if (!t) return "—";
  return new Date(t).toLocaleString("ja-JP");
}

async function load() {
  if (!token) return;
  try {
    const printers = await api("GET", "/api/printers");
    const rows = printers.map((p) => {
      const tr = el("tr");
      tr.append(el("td", p.name), el("td", p.url), el("td", p.route), el("td"));
      tr.lastChild.append(button("削除", () => api("DELETE", "/api/printers/" + p.name)));
      return tr;
    });
    $("printers").replaceChildren(...rows);

    const locations = await api("GET", "/api/locations");
    const cards = [];
    for (const l of locations) {
      const card = el("div", undefined, { className: "card" });
      const fb = l.fallback ? `紙に逃がす: ${l.fallback.printer}${l.fallback.after ? "（" + l.fallback.after + " で）" : "（端末がつながっていないとき）"}` : "紙に逃がさない";
      card.append(el("h3", `${l.title || l.name} (${l.name})`), el("p", `${fb} · /ipp/print/${l.name}`, { className: "note" }));
      const board = await api("GET", `/api/locations/${l.name}/board`);
      card.append(el("p", `未完了 ${board.open.length} 件 · 完了・取り消し ${board.ended.length} 件`));
      const devices = await api("GET", `/api/locations/${l.name}/devices`);
      const ul = el("ul", undefined, { className: "devices" });
      for (const d of devices) {
        const li = el("li", `${d.name} · 最後の接続 ${since(d.last_seen)} `);
        li.append(button("外す", () => api("DELETE", "/api/devices/" + d.id)));
        ul.append(li);
      }
      card.append(ul);
      card.append(
        button("端末を足す（登録コード）", async () => {
          const e = await api("POST", `/api/locations/${l.name}/enrollments`);
          $("enrollment").hidden = false;
          $("qr").innerHTML = e.qr; // SVG the manager drew
          $("enroll-code").textContent = e.code.replace(/(.{5})/, "$1-");
          $("enroll-url").textContent = e.url;
          $("enroll-url").href = e.url;
          $("enroll-expires").textContent = `${since(e.expires)} まで、1 回だけ使えます`;
          $("enrollment").scrollIntoView();
        }),
        button("削除", () => api("DELETE", "/api/locations/" + l.name)),
      );
      cards.push(card);
    }
    $("locations").replaceChildren(...cards);
    status("");
  } catch (e) {
    status(e.message, true);
  }
}

$("token-form").addEventListener("submit", (ev) => {
  ev.preventDefault();
  token = $("token").value;
  try {
    sessionStorage.setItem("bdf-pm-admin", token);
  } catch {}
  load();
});
$("tick").addEventListener("click", async () => {
  try {
    const r = await api("POST", "/api/tick");
    status(`tick: 刷った ${r.printed}、失敗 ${r.failed}、片付けた ${r.removed}`);
    await load();
  } catch (e) {
    status(e.message, true);
  }
});
$("printer-form").addEventListener("submit", async (ev) => {
  ev.preventDefault();
  try {
    await api("PUT", "/api/printers/" + $("printer-name").value, { url: $("printer-url").value });
    ev.target.reset();
    await load();
  } catch (e) {
    status(e.message, true);
  }
});
$("location-form").addEventListener("submit", async (ev) => {
  ev.preventDefault();
  const body = { title: $("location-title").value };
  if ($("location-printer").value) {
    body.fallback = { printer: $("location-printer").value };
    if ($("location-after").value) body.fallback.after = $("location-after").value;
  }
  try {
    await api("PUT", "/api/locations/" + $("location-name").value, body);
    ev.target.reset();
    await load();
  } catch (e) {
    status(e.message, true);
  }
});
load();
