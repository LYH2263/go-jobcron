const tbody = document.querySelector("#jobs tbody");
const statsEl = document.getElementById("stats");

async function api(path, opts) {
  const res = await fetch(path, opts);
  const data = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(data.error || res.statusText);
  return data;
}

async function refresh() {
  const [jobs, stats] = await Promise.all([
    api("/api/jobs?limit=50"),
    api("/api/stats"),
  ]);
  statsEl.textContent = `pending=${stats.Pending} executed=${stats.Executed} closed=${stats.Closed}`;
  tbody.innerHTML = "";
  for (const j of jobs) {
    const tr = document.createElement("tr");
    tr.innerHTML = `
      <td><code>${j.ID}</code></td>
      <td>${j.Name}</td>
      <td>${j.Kind}</td>
      <td>${j.Status}</td>
      <td>${j.NextRun || "-"}</td>
      <td>
        <button data-pause="${j.ID}">暂停</button>
        <button class="secondary" data-resume="${j.ID}">恢复</button>
      </td>`;
    tbody.appendChild(tr);
  }
}

document.getElementById("refresh").onclick = refresh;
document.getElementById("run-due").onclick = async () => {
  await api("/api/run?max=5", { method: "POST" });
  await refresh();
};

document.getElementById("create-form").onsubmit = async (e) => {
  e.preventDefault();
  const fd = new FormData(e.target);
  const body = {
    name: fd.get("name"),
    kind: fd.get("kind"),
    payload: JSON.parse(fd.get("payload") || "{}"),
    run_at: fd.get("run_at") || "",
    cron_expr: fd.get("cron_expr") || "",
    url: fd.get("url") || "",
  };
  await api("/api/jobs", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  e.target.reset();
  await refresh();
};

tbody.onclick = async (e) => {
  const pause = e.target.getAttribute("data-pause");
  const resume = e.target.getAttribute("data-resume");
  if (pause) {
    await api(`/api/jobs/${pause}?action=pause`, { method: "POST" });
    await refresh();
  }
  if (resume) {
    await api(`/api/jobs/${resume}?action=resume`, { method: "POST" });
    await refresh();
  }
};

refresh();
