import { useEffect, useState } from "react";
import { Bar } from "react-chartjs-2";
import {
  Chart as ChartJS,
  BarElement,
  CategoryScale,
  LinearScale,
  Legend,
  Tooltip,
} from "chart.js";

ChartJS.register(BarElement, CategoryScale, LinearScale, Legend, Tooltip);

export default function PowerballChart() {
  const [rows, setRows] = useState([]);
  const [err, setErr] = useState("");

  useEffect(() => {
    fetch("/api/powerball/sessions")
      .then((r) => {
        if (!r.ok) throw new Error("API error " + r.status);
        return r.json();
      })
      .then(setRows)
      .catch((e) => setErr(e.message));
  }, []);

  if (err) return <p style={{color: "crimson"}}>Error: {err}</p>;
  if (!rows.length) return <p>Loading…</p>;

  const labels = rows.map((r) => r.date);
  const classical = rows.map((r) => r.classical || 0);
  const quantum = rows.map((r) => r.quantum || 0);

  const data = {
    labels,
    datasets: [
      { label: "classical", data: classical, stack: "stack1" },
      { label: "quantum", data: quantum, stack: "stack1" },
    ],
  };

  const options = {
    responsive: true,
    plugins: { legend: { position: "top" } },
    scales: { x: { stacked: true }, y: { stacked: true, beginAtZero: true } },
  };

  return (
    <div style={{maxWidth: 800}}>
      <Bar data={data} options={options} />
    </div>
  );
}
