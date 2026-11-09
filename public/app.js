// Frontend cảnh báo sạt lở Bình Định — HTK 0396241674 — 09/11/2026
// Gọi API Go Fiber: GET /api/points, POST /api/risk
"use strict";

let allPoints = [];
let chart = null;

async function loadPoints() {
  const res = await fetch("/api/points");
  const data = await res.json();
  allPoints = data.points || [];
  renderTable("");
  renderChart(allPoints);
}

function renderTable(district) {
  const tbody = document.querySelector("#points-table tbody");
  const rows = allPoints.filter((p) => !district || p.huyen === district);
  tbody.innerHTML = rows
    .map(
      (p) => `<tr>
        <td>${p.id}</td>
        <td>${p.ten}</td>
        <td>${p.huyen}</td>
        <td>${p.mua24h} mm</td>
        <td><strong>${p.chiSo}</strong></td>
        <td><span class="pill" style="background:${p.mau}">${p.mucDo}</span></td>
        <td>${p.hoDan}</td>
        <td>${p.ghiChu}</td>
      </tr>`
    )
    .join("");
}

function renderChart(points) {
  const ctx = document.getElementById("riskChart");
  if (!ctx || typeof Chart === "undefined") return;
  if (chart) chart.destroy();
  chart = new Chart(ctx, {
    type: "bar",
    data: {
      labels: points.map((p) => `#${p.id} ${p.ten.split("(")[0].trim()}`),
      datasets: [
        {
          label: "Chỉ số nguy cơ sạt lở (0–100)",
          data: points.map((p) => p.chiSo),
          backgroundColor: points.map((p) => p.mau),
          borderRadius: 5,
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      scales: {
        y: { min: 0, max: 100, ticks: { color: "#94a3b8" }, grid: { color: "#334155" } },
        x: { ticks: { color: "#94a3b8", maxRotation: 60, minRotation: 45, font: { size: 10 } }, grid: { display: false } },
      },
      plugins: {
        legend: { labels: { color: "#e2e8f0" } },
        tooltip: {
          callbacks: {
            afterLabel: (item) => {
              const p = points[item.dataIndex];
              return [`Mức: ${p.mucDo}`, `Mưa 24h: ${p.mua24h} mm`, `Khuyến nghị: ${p.khuyenNghi}`];
            },
          },
        },
      },
    },
  });
}

async function calcRisk(e) {
  e.preventDefault();
  const payload = {
    rain: Number(document.getElementById("rain").value),
    slope: Number(document.getElementById("slope").value),
    soil: Number(document.getElementById("soil").value),
  };
  const box = document.getElementById("risk-result");
  try {
    const res = await fetch("/api/risk", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    const data = await res.json();
    if (!data.ok) throw new Error(data.error || "API lỗi");
    box.hidden = false;
    box.style.borderColor = data.color;
    document.getElementById("risk-value").textContent = data.riskIndex;
    document.getElementById("risk-value").style.color = data.color;
    const lvl = document.getElementById("risk-level");
    lvl.textContent = data.level;
    lvl.style.color = data.color;
    document.getElementById("risk-advice").textContent = data.advice;
  } catch (err) {
    box.hidden = false;
    document.getElementById("risk-advice").textContent =
      "Không gọi được API: " + err.message + ". Thử lại hoặc gọi HTK 0396241674.";
  }
}

document.addEventListener("DOMContentLoaded", () => {
  document.getElementById("risk-form").addEventListener("submit", calcRisk);
  document.getElementById("district-filter").addEventListener("change", (e) => {
    renderTable(e.target.value);
    renderChart(allPoints.filter((p) => !e.target.value || p.huyen === e.target.value));
  });
  loadPoints().catch(() => {
    document.querySelector("#points-table tbody").innerHTML =
      '<tr><td colspan="8">Không tải được dữ liệu. Gọi hotline HTK 0396241674 nếu cần hỗ trợ.</td></tr>';
  });
});