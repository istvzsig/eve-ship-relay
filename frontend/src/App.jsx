import { useEffect, useState } from "react";
import "./App.css";

const API = "http://localhost:8080";

function App() {
  const [shipments, setShipments] = useState([]);

  useEffect(() => {
    fetch(`${API}/api/shipments`)
      .then((res) => res.json())
      .then(setShipments);
  }, []);

  return (
    <div className="app">
      <header>
        <div>
          <h1>ShipRelay</h1>
          <p>EVE Online logistics</p>
        </div>
        <div className="status">● ONLINE</div>
      </header>

      <main>
        <div className="page-header">
          <div>
            <h2>Shipments</h2>
            <p>Manage incoming and active shipments.</p>
          </div>
          <button>+ New Shipment</button>
        </div>

        <div className="shipments">
          {shipments.map((s) => (
            <div className="shipment-card" key={s.id}>
              <div>
                <strong>#{s.id}</strong>
                <span className="ship">{s.ship}</span>
              </div>

              <div className="route">
                {s.origin} → {s.destination}
              </div>

              <div className={`badge ${s.status.toLowerCase()}`}>
                {s.status.replace("_", " ")}
              </div>
            </div>
          ))}
        </div>
      </main>
    </div>
  );
}

export default App;
