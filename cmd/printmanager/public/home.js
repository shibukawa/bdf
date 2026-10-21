// Everyone's page: the jobs are drawn by the server; this sends a file to
// a location as the person signed in (their session's cookie).

const $ = (id) => document.getElementById(id);

function status(text, bad) {
  const t = document.createElement("div");
  t.role = "status";
  t.className = "alert shadow-lg " + (bad ? "alert-error" : "alert-success");
  t.textContent = text;
  $("toasts").append(t);
  setTimeout(() => t.remove(), bad ? 6000 : 3000);
}

const form = $("upload");
if (form) {
  form.addEventListener("submit", async (ev) => {
    ev.preventDefault();
    const file = $("upload-file").files[0];
    if (!file) return;
    const button = form.querySelector("button");
    button.disabled = true;
    try {
      const res = await fetch(`/api/locations/${encodeURIComponent($("upload-location").value)}/slips?name=${encodeURIComponent(file.name)}`, {
        method: "POST",
        body: file,
      });
      const out = await res.json();
      if (!res.ok) throw new Error(out.error || res.statusText);
      status(`送りました（#${out.job}）`);
      setTimeout(() => location.reload(), 600);
    } catch (e) {
      status(e.message, true);
    } finally {
      button.disabled = false;
    }
  });
}
