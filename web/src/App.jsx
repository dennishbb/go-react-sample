import PowerballChart from "./PowerballChart";

export default function App() {
  return (
    <div style={{ fontFamily: "system-ui, sans-serif", padding: 24 }}>
      <h1>Go + React App 🚀</h1>
      <p>Powerball sessions per day (classical vs quantum)</p>
      <PowerballChart />
    </div>
  );
}
