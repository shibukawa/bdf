// The administrator's page: printers, locations, devices and their
// enrollment codes, through the API with the admin token (kept in this
// tab's sessionStorage only).

const $ = (id) => document.getElementById(id);
let token = "";
try {
  token = sessionStorage.getItem("bdf-pm-admin") || "";
} catch {}
$("token").value = token;

// status shows a message for a while, as a toast.
function status(text, bad) {
  if (!text) return;
  const t = document.createElement("div");
  t.role = "status";
  t.className = "alert shadow-lg " + (bad ? "alert-error" : "alert-success");
  t.textContent = text;
  $("toasts").append(t);
  setTimeout(() => t.remove(), bad ? 6000 : 3000);
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

function button(text, fn, className = "btn btn-sm") {
  const b = el("button", text, { type: "button", className });
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

function stat(title, value) {
  const s = el("div", undefined, { className: "stat px-4 py-2" });
  s.append(el("div", title, { className: "stat-title" }), el("div", String(value), { className: "stat-value text-2xl" }));
  return s;
}

async function load() {
  if (!token) return;
  try {
    const printers = await api("GET", "/api/printers");
    const rows = printers.map((p) => {
      const tr = el("tr");
      const route = el("span", p.route, { className: "badge badge-ghost" });
      tr.append(el("td", p.name, { className: "font-semibold" }), el("td", p.url, { className: "font-mono text-sm break-all" }), el("td"), el("td"));
      tr.children[2].append(route);
      tr.lastChild.append(button("削除", () => api("DELETE", "/api/printers/" + p.name), "btn btn-sm btn-ghost text-error"));
      return tr;
    });
    if (!rows.length) {
      const tr = el("tr");
      tr.append(el("td", "プリンタはまだありません", { colSpan: 4, className: "text-base-content/50" }));
      rows.push(tr);
    }
    $("printers").replaceChildren(...rows);

    const locations = await api("GET", "/api/locations");
    const cards = [];
    for (const l of locations) {
      const card = el("div", undefined, { className: "card bg-base-100 shadow-sm" });
      const body = el("div", undefined, { className: "card-body gap-3" });
      const head = el("div", undefined, { className: "flex flex-wrap items-center gap-2" });
      head.append(el("h3", l.title || l.name, { className: "card-title" }), el("span", l.name, { className: "badge badge-ghost font-mono" }));
      const fb = l.fallback
        ? el("span", `紙に逃がす: ${l.fallback.printer}${l.fallback.after ? "（" + l.fallback.after + " で）" : "（端末がつながっていないとき）"}`, { className: "badge badge-info badge-soft" })
        : el("span", "紙に逃がさない", { className: "badge badge-ghost" });
      head.append(fb);
      body.append(head, el("p", `IPP: /ipp/print/${l.name}`, { className: "font-mono text-sm text-base-content/60" }));
      const board = await api("GET", `/api/locations/${l.name}/board`);
      const stats = el("div", undefined, { className: "stats stats-horizontal bg-base-200" });
      stats.append(stat("未完了", board.open.length), stat("完了・取り消し", board.ended.length));
      body.append(stats);
      const devices = await api("GET", `/api/locations/${l.name}/devices`);
      const table = el("table", undefined, { className: "table table-sm" });
      const tbody = el("tbody");
      for (const d of devices) {
        const tr = el("tr");
        tr.append(el("td", d.name, { className: "font-semibold" }), el("td", "最後の接続 " + since(d.last_seen), { className: "text-sm text-base-content/60" }), el("td"));
        tr.lastChild.append(button("外す", () => api("DELETE", "/api/devices/" + d.id), "btn btn-xs btn-ghost text-error"));
        tbody.append(tr);
      }
      if (!devices.length) {
        const tr = el("tr");
        tr.append(el("td", "端末はまだありません", { className: "text-base-content/50" }));
        tbody.append(tr);
      }
      table.append(tbody);
      body.append(table);
      const actions = el("div", undefined, { className: "card-actions justify-end" });
      actions.append(
        button("削除", () => api("DELETE", "/api/locations/" + l.name), "btn btn-sm btn-ghost text-error"),
        button(
          "端末を足す（登録コード）",
          async () => {
            const e = await api("POST", `/api/locations/${l.name}/enrollments`);
            $("qr").innerHTML = e.qr; // SVG the manager drew
            $("enroll-code").textContent = e.code.replace(/(.{5})/, "$1-");
            $("enroll-url").textContent = e.url;
            $("enroll-url").href = e.url;
            $("enroll-expires").textContent = `${since(e.expires)} まで、1 回だけ使えます`;
            $("enroll-modal").showModal();
          },
          "btn btn-sm btn-primary",
        ),
      );
      body.append(actions);
      card.append(body);
      cards.push(card);
    }
    if (!cards.length) cards.push(el("p", "場所はまだありません", { className: "text-base-content/50" }));
    $("locations").replaceChildren(...cards);
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
