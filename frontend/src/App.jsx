import { useEffect, useState } from "react";
import "./App.css";

const API = "http://localhost:8080";

function App() {
  const [shipments, setShipments] = useState([]);
  const [selected, setSelected] = useState(null);
  const [scanning, setScanning] = useState(false);

  const [showNewShipment, setShowNewShipment] = useState(false);
  const [form, setForm] = useState({
    ship: "",
    origin: "",
    destination: "",
    contract_id: "",
    receipt_code: "",
  });

  useEffect(() => {
    fetch(`${API}/api/shipments`)
      .then((res) => res.json())
      .then(setShipments);
  }, []);

  async function openShipment(id) {
    const res = await fetch(`${API}/api/shipments/${id}`);
    const shipment = await res.json();
    setSelected(shipment);
  }

  async function runScan(id) {
    setScanning(true);

    try {
      const res = await fetch(`${API}/api/shipments/${id}/scan`, {
        method: "POST",
      });

      const shipment = await res.json();

      setSelected(shipment);

      setShipments((current) =>
        current.map((item) => (item.id === shipment.id ? shipment : item)),
      );
    } finally {
      setScanning(false);
    }
  }

  async function createShipment(e) {
    e.preventDefault();

    const res = await fetch(`${API}/api/shipments`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(form),
    });

    const shipment = await res.json();

    setShipments((current) => [shipment, ...current]);
    setSelected(shipment);
    setShowNewShipment(false);

    setForm({
      ship: "",
      origin: "",
      destination: "",
      contract_id: "",
      receipt_code: "",
    });
  }

  async function verifyPayment(id) {
    const res = await fetch(`${API}/api/shipments/${id}/verify-payment`, {
      method: "POST",
    });

    const shipment = await res.json();

    setSelected(shipment);

    setShipments((current) =>
      current.map((item) => (item.id === shipment.id ? shipment : item)),
    );
  }

  if (selected) {
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
          <button className="back" onClick={() => setSelected(null)}>
            ← Back to shipments
          </button>

          <div className="detail-header">
            <div>
              <h2>Shipment #{selected.id}</h2>
              <p>
                {selected.ship} · {selected.origin} → {selected.destination}
              </p>
            </div>

            <div className={`badge ${selected.status.toLowerCase()}`}>
              {selected.status.replace("_", " ")}
            </div>
          </div>

          <div className="detail-grid">
            <section className="panel">
              <h3>Contract</h3>

              <div className="row">
                <span>Contract</span>
                <strong>{selected.contract.id}</strong>
              </div>

              <div className="row">
                <span>Payment</span>
                <div>
                  <strong>
                    {selected.contract.payment_verified
                      ? "✓ Verified"
                      : "✗ Failed"}
                  </strong>

                  {!selected.contract.payment_verified && (
                    <button
                      onClick={() => verifyPayment(selected.id)}
                      style={{ marginLeft: "12px" }}
                    >
                      Verify Payment
                    </button>
                  )}
                </div>
              </div>

              <div className="row">
                <span>Receipt code</span>
                <strong>{selected.contract.receipt_code}</strong>
              </div>
            </section>

            <section className="panel">
              <h3>Abyssal modules</h3>

              {!selected.abyssal_modules?.length ? (
                <p className="muted">No abyssal modules.</p>
              ) : (
                selected.abyssal_modules.map((module) => (
                  <div className="module" key={module.name}>
                    <strong>{module.name}</strong>

                    <div className="row">
                      <span>Value</span>
                      <strong>{formatISK(module.value_isk)}</strong>
                    </div>

                    <div className="row">
                      <span>Deductible</span>
                      <strong>{formatISK(module.deductible_isk)}</strong>
                    </div>

                    <div className="row">
                      <span>Status</span>
                      <strong>
                        {module.deductible_paid ? "✓ Paid" : "⚠ Unpaid"}
                      </strong>
                    </div>
                  </div>
                ))
              )}
            </section>

            <section className="panel">
              <h3>Asset scan</h3>

              <button onClick={() => runScan(selected.id)} disabled={scanning}>
                {scanning ? "Scanning..." : "Run Asset Scan"}
              </button>

              <div className="scan-result">
                <div
                  className={selected.asset_scan.passed ? "check" : "failed"}
                >
                  {selected.asset_scan?.passed ? "✓" : "✗"}
                </div>

                <div>
                  <strong>
                    {selected.asset_scan?.passed
                      ? "Scan passed"
                      : "Scan failed"}
                  </strong>
                  <p>
                    {selected.asset_scan?.passed
                      ? "Expected assets are present."
                      : "Expected assets do not match the ship."}
                  </p>
                </div>
              </div>

              <div className="assets">
                <div>
                  <h4>Expected</h4>
                  {(selected.asset_scan?.expected_modules ?? []).map(
                    (module) => (
                      <div key={module}>✓ {module}</div>
                    ),
                  )}
                </div>

                <div>
                  <h4>Actual</h4>
                  {(selected.asset_scan?.actual_modules ?? []).length === 0
                    ? "No modules detected"
                    : selected.asset_scan.actual_modules.map((module) => (
                        <div key={module}>✓ {module}</div>
                      ))}
                </div>
              </div>
            </section>
          </div>
        </main>
      </div>
    );
  }

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

          <button onClick={() => setShowNewShipment(true)}>
            + New Shipment
          </button>
        </div>

        {showNewShipment && (
          <div className="panel new-shipment">
            <h3>New Shipment</h3>

            <form onSubmit={createShipment}>
              <input
                placeholder="Ship"
                value={form.ship}
                onChange={(e) => setForm({ ...form, ship: e.target.value })}
                required
              />

              <input
                placeholder="Origin"
                value={form.origin}
                onChange={(e) => setForm({ ...form, origin: e.target.value })}
                required
              />

              <input
                placeholder="Destination"
                value={form.destination}
                onChange={(e) =>
                  setForm({ ...form, destination: e.target.value })
                }
                required
              />

              <input
                placeholder="Contract ID"
                value={form.contract_id}
                onChange={(e) =>
                  setForm({ ...form, contract_id: e.target.value })
                }
                required
              />

              <input
                placeholder="Receipt code"
                value={form.receipt_code}
                onChange={(e) =>
                  setForm({ ...form, receipt_code: e.target.value })
                }
                required
              />

              <div className="form-actions">
                <button type="submit">Create Shipment</button>

                <button
                  type="button"
                  className="back"
                  onClick={() => setShowNewShipment(false)}
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        )}

        <div className="shipments">
          {shipments.map((s) => (
            <button
              className="shipment-card"
              key={s.id}
              onClick={() => openShipment(s.id)}
            >
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
            </button>
          ))}
        </div>
      </main>
    </div>
  );
}

function formatISK(value) {
  return `${new Intl.NumberFormat("en-US").format(value)} ISK`;
}

export default App;
