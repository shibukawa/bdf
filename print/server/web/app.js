// The page of a bdf printd's state: it reads the API and follows /events.
// It only reads; jobs and queues are changed with the API and bdf jobs.
"use strict";

const state = { queues: new Map(), printers: new Map(), jobs: new Map() };
let token = "";
try { token = sessionStorage.getItem("bdf-token") || ""; } catch (e) {}
let stream = null;

function headers() {
  return token ? { Authorization: "Bearer " + token } : {};
}

async function api(path) {
  const r = await fetch("/api/v1" + path, { headers: headers() });
  if (r.status === 401) {
    document.getElementById("login").classList.add("shown");
    throw new Error("a token is required");
  }
  if (!r.ok) throw new Error(path + ": " + r.status);
  return r.json();
}

function setStatus(text, cls) {
  const s = document.getElementById("status");
  s.textContent = text;
  s.className = cls || "";
}

async function load() {
  try {
    const [q, p, j] = await Promise.all([api("/queues"), api("/printers"), api("/jobs?state=all&limit=200")]);
    state.queues = new Map(q.queues.map((x) => [x.name, x]));
    state.printers = new Map(p.printers.map((x) => [x.name, x]));
    state.jobs = new Map(j.jobs.map((x) => [x.id, x]));
    render();
    follow();
  } catch (e) {
    setStatus(e.message, "bad");
  }
}

// follow reads /events; it reads everything again when the stream ends.
async function follow() {
  if (stream) stream.abort();
  const ctl = new AbortController();
  stream = ctl;
  try {
    const r = await fetch("/api/v1/events", { headers: headers(), signal: ctl.signal });
    if (!r.ok) throw new Error("events: " + r.status);
    setStatus("live", "ok");
    const reader = r.body.pipeThrough(new TextDecoderStream()).getReader();
    let buf = "";
    for (;;) {
      const { value, done } = await reader.read();
      if (done) break;
      buf += value;
      let i;
      while ((i = buf.indexOf("\n\n")) >= 0) {
        event(buf.slice(0, i));
        buf = buf.slice(i + 2);
      }
    }
  } catch (e) {
    if (ctl.signal.aborted) return;
  }
  if (stream !== ctl) return;
  setStatus("disconnected; reconnecting…", "warn");
  setTimeout(load, 3000);
}

function event(block) {
  let name = "message", data = "";
  for (const line of block.split("\n")) {
    if (line.startsWith("event: ")) name = line.slice(7);
    else if (line.startsWith("data: ")) data += line.slice(6);
  }
  if (!data) return;
  const v = JSON.parse(data);
  if (name === "job") state.jobs.set(v.id, v);
  else if (name === "queue") state.queues.set(v.name, v);
  else if (name === "printer") state.printers.set(v.name, v);
  render();
}

// cell makes a cell of text: what a job is named is never read as HTML.
function cell(text, cls) {
  const td = document.createElement("td");
  td.textContent = text == null ? "" : String(text);
  if (cls) td.className = cls;
  return td;
}

function fill(id, rows) {
  const body = document.querySelector("#" + id + " tbody");
  body.replaceChildren(...rows.map((cells) => {
    const tr = document.createElement("tr");
    tr.append(...cells);
    return tr;
  }));
}

function when(t) {
  return t ? new Date(t).toLocaleString() : "";
}

function jobClass(st) {
  if (st === "completed") return "ok";
  if (st === "aborted" || st === "canceled") return "bad";
  if (st === "pending-held" || st === "processing-stopped") return "warn";
  return "";
}

function render() {
  fill("queues", [...state.queues.values()].sort((a, b) => a.name.localeCompare(b.name)).map((q) => [
    cell(q.name), cell((q.printers || []).join(", ")), cell(q.policy),
    cell([q.paused ? "paused" : "", q.rejecting ? "rejecting" : ""].filter(Boolean).join(", ") || "running",
      q.paused || q.rejecting ? "warn" : "ok"),
    cell(q.jobs["pending"] || 0, "num"), cell(q.jobs["pending-held"] || 0, "num"),
    cell((q.jobs["processing"] || 0) + (q.jobs["processing-stopped"] || 0), "num"),
  ]));
  fill("printers", [...state.printers.values()].sort((a, b) => a.name.localeCompare(b.name)).map((p) => [
    cell(p.name), cell(p.up ? p.state || "up" : "down", !p.up ? "bad" : p.state === "stopped" ? "warn" : "ok"),
    cell((p.reasons || []).filter((r) => r !== "none").join(", ") || p.error || "", "wrap"),
    cell(p.jobs, "num"), cell(p["make-and-model"] || "", "muted"), cell(p.url, "muted"),
  ]));
  fill("jobs", [...state.jobs.values()].sort((a, b) => b.id - a.id).slice(0, 200).map((j) => [
    cell(j.id, "num"), cell(j.state, jobClass(j.state)), cell((j.reasons || []).join(", "), "wrap"),
    cell(j.queue), cell(j.name), cell(j.user), cell(j.printer ? j.printer.name : ""),
    cell(when(j.submitted), "muted"), cell(when(j.finished), "muted"),
  ]));
}

document.getElementById("login").addEventListener("submit", (e) => {
  e.preventDefault();
  token = document.getElementById("token").value.trim();
  try { sessionStorage.setItem("bdf-token", token); } catch (err) {}
  document.getElementById("login").classList.remove("shown");
  load();
});

load();
