// The administrator's page: every job, printers, locations, devices and
// their enrollment codes, through the API as the admin signed in (the
// session's cookie; the server shows this page to admins only).

const $ = (id) => document.getElementById(id);

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
    headers: body ? { "Content-Type": "application/json" } : {},
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

const kinds = { shared: "共有", group: "グループ", personal: "個人" };

// kind is a badge of who may print to a location or printer.
function kind(k, owner) {
  const cls =
    { group: "badge badge-sm badge-soft badge-secondary whitespace-nowrap", personal: "badge badge-sm badge-soft badge-accent whitespace-nowrap" }[k] ||
    "badge badge-sm badge-soft badge-neutral whitespace-nowrap";
  return el("span", kinds[k || "shared"] + (owner ? ": " + owner : ""), { className: cls });
}

const states = {
  new: ["未開封", "badge badge-sm badge-warning whitespace-nowrap"],
  opened: ["開いた", "badge badge-sm badge-info whitespace-nowrap"],
  done: ["完了", "badge badge-sm badge-success whitespace-nowrap"],
  canceled: ["取り消し", "badge badge-sm badge-ghost whitespace-nowrap"],
};
const via = { oidc: "サインイン", "submit-token": "送信のトークン", "admin-token": "管理のトークン", none: "確かめていない" };

async function loadSlips() {
  const slips = await api("GET", "/api/slips");
  const rows = slips.map((s) => {
    const tr = el("tr");
    const [label, cls] = states[s.state] || [s.state, "badge badge-sm"];
    const by = s.submitted_by || { kind: "anonymous", via: "none" };
    const who = el("td");
    who.append(el("div", by.name || by.id || "—", { className: "font-semibold" }), el("div", via[by.via] || by.via, { className: "text-xs text-base-content/60" }));
    const st = el("td");
    st.append(el("span", label, { className: cls }));
    tr.append(
      el("td", String(s.job), { className: "font-mono" }),
      el("td", s.name || "slip-" + s.id, { className: "max-w-48 break-words" }),
      el("td", s.location, { className: "font-mono text-sm" }),
      st,
      who,
      el("td", s.claimed_user || "—", { className: "hidden sm:table-cell" }),
      el("td", s.group || "—", { className: "hidden md:table-cell" }),
    );
    return tr;
  });
  if (!rows.length) {
    const tr = el("tr");
    tr.append(el("td", "ジョブはまだありません", { colSpan: 7, className: "text-base-content/50" }));
    rows.push(tr);
  }
  $("slips").replaceChildren(...rows);
}

async function load() {
  try {
    await loadSlips();
    const printers = await api("GET", "/api/printers");
    const rows = printers.map((p) => {
      const tr = el("tr");
      const route = el("span", p.route, { className: "badge badge-ghost" });
      tr.append(el("td", p.name, { className: "font-semibold" }), el("td", p.url, { className: "font-mono text-sm break-all" }), el("td"), el("td"), el("td"));
      tr.children[2].append(route);
      tr.children[3].append(kind(p.kind, p.owner));
      tr.lastChild.append(button("削除", () => api("DELETE", "/api/printers/" + p.name), "btn btn-sm btn-ghost whitespace-nowrap text-error"));
      return tr;
    });
    if (!rows.length) {
      const tr = el("tr");
      tr.append(el("td", "プリンタはまだありません", { colSpan: 5, className: "text-base-content/50" }));
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
      head.append(fb, kind(l.kind, l.owner));
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
    const body = { url: $("printer-url").value, kind: $("printer-kind").value };
    if ($("printer-owner").value) body.owner = $("printer-owner").value;
    await api("PUT", "/api/printers/" + $("printer-name").value, body);
    ev.target.reset();
    await load();
  } catch (e) {
    status(e.message, true);
  }
});
$("location-form").addEventListener("submit", async (ev) => {
  ev.preventDefault();
  const body = { title: $("location-title").value, kind: $("location-kind").value };
  if ($("location-owner").value) body.owner = $("location-owner").value;
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
