import { useEffect, useState } from "react";
import "./App.css";

import RouteMap from "./RouteMap";

const API = import.meta.env.VITE_API_URL || "http://localhost:8080";

const carriers = [
  { id: "CARRIER-01", name: "Night Hauler" },
  { id: "CARRIER-02", name: "Red Freighter" },
  { id: "CARRIER-03", name: "Void Runner" },
];

const cynoPilots = [
  { id: "CYNO-01", name: "Dark Angel" },
  { id: "CYNO-02", name: "Nightwatch" },
  { id: "CYNO-03", name: "Black Lantern" },
];

function App() {
  const [shipments, setShipments] = useState([]);
  const [selected, setSelected] = useState(null);
  const [scanning, setScanning] = useState(false);
  const [showNewShipment, setShowNewShipment] = useState(false);
  const [carrierId, setCarrierId] = useState("");
  const [cynoPilotId, setCynoPilotId] = useState("");
  const [newShipment, setNewShipment] = useState({
    ship: "",
    origin: "",
    destination: "",
    contract_id: "",
    receipt_code: "",
  });
  const [dispatching, setDispatching] = useState(false);
  const [delivering, setDelivering] = useState(false);

  useEffect(() => {
    let cancelled = false;

    fetch(`${API}/api/shipments`)
      .then((res) => res.json())
      .then((data) => {
        if (!cancelled) {
          setShipments(data);
        }
      })
      .catch((error) => {
        console.error("Failed to load shipments:", error);
      });

    return () => {
      cancelled = true;
    };
  }, []);

  async function openShipment(id) {
    const res = await fetch(`${API}/api/shipments/${id}`);
    const shipment = await res.json();

    setSelected(shipment);
    setCarrierId(shipment.carrier_id || "");
    setCynoPilotId(shipment.cyno_pilot_id || "");
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

  async function createShipment(event) {
    event.preventDefault();

    const res = await fetch(`${API}/api/shipments`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(newShipment),
    });

    const shipment = await res.json();

    setShipments((current) => [shipment, ...current]);
    setSelected(shipment);
    setShowNewShipment(false);

    setNewShipment({
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

  async function assignCarrier(id) {
    if (!carrierId) {
      return;
    }

    const res = await fetch(`${API}/api/shipments/${id}/assign-carrier`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        carrier_id: carrierId,
      }),
    });

    const shipment = await res.json();

    setSelected(shipment);

    setShipments((current) =>
      current.map((item) => (item.id === shipment.id ? shipment : item)),
    );

    setCarrierId("");
  }

  async function assignCyno(id) {
    if (!cynoPilotId) {
      return;
    }

    const res = await fetch(`${API}/api/shipments/${id}/assign-cyno`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        cyno_pilot_id: cynoPilotId,
      }),
    });

    const shipment = await res.json();

    setSelected(shipment);

    setShipments((current) =>
      current.map((item) => (item.id === shipment.id ? shipment : item)),
    );

    setCynoPilotId("");
  }

  async function dispatchShipment(id) {
    setDispatching(true);

    try {
      const res = await fetch(`${API}/api/shipments/${id}/dispatch`, {
        method: "POST",
      });

      if (!res.ok) {
        const message = await res.text();
        alert(message);
        return;
      }

      const shipment = await res.json();

      setSelected(shipment);

      setShipments((current) =>
        current.map((item) => (item.id === shipment.id ? shipment : item)),
      );
    } finally {
      setDispatching(false);
    }
  }

  async function deliverShipment(id) {
    setDelivering(true);

    try {
      const res = await fetch(`${API}/api/shipments/${id}/deliver`, {
        method: "POST",
      });

      if (!res.ok) {
        const message = await res.text();
        alert(message);
        return;
      }

      const shipment = await res.json();

      setSelected(shipment);

      setShipments((current) =>
        current.map((item) => (item.id === shipment.id ? shipment : item)),
      );
    } finally {
      setDelivering(false);
    }
  }

  function updateNewShipment(field, value) {
    setNewShipment((current) => ({
      ...current,
      [field]: value,
    }));
  }

  function statusClass(status) {
    return status.toLowerCase();
  }

  function statusLabel(status) {
    return status.replace("_", " ");
  }

  const readyCount = shipments.filter(
    (shipment) => shipment.status === "READY",
  ).length;

  const blockedCount = shipments.filter(
    (shipment) => shipment.status === "BLOCKED",
  ).length;

  const inTransitCount = shipments.filter(
    (shipment) => shipment.status === "IN_TRANSIT",
  ).length;

  const deliveredCount = shipments.filter(
    (shipment) => shipment.status === "DELIVERED",
  ).length;

  return (
    <div className="app">
      <header className="header">
        <div>
          <h1>ShipRelay</h1>
          <span>EVE Online Logistics</span>
        </div>

        <div className="online">
          <span className="online-dot" />
          ONLINE
        </div>
      </header>

      <main>
        {selected ? (
          <>
            <button className="back-button" onClick={() => setSelected(null)}>
              ← Back to Operations
            </button>

            <div className="page-header">
              <div>
                <span className="shipment-id">SHIPMENT #{selected.id}</span>

                <h2>{selected.ship}</h2>

                <p>
                  {selected.origin} → {selected.destination}
                </p>
              </div>

              <span className={`status ${statusClass(selected.status)}`}>
                {statusLabel(selected.status)}
              </span>
            </div>

            <div className="detail-grid">
              <RouteMap
                apiURL={API}
                origin={selected.origin}
                destination={selected.destination}
              />
              <section className="panel">
                <h3>Contract</h3>

                <div className="row">
                  <span>Contract ID</span>
                  <strong>{selected.contract?.id || "Unknown"}</strong>
                </div>

                <div className="row">
                  <span>Payment</span>
                  <strong>
                    {selected.contract?.payment_verified
                      ? "✓ Verified"
                      : "✕ Not Verified"}
                  </strong>
                </div>

                <div className="row">
                  <span>Receipt Code</span>
                  <strong>{selected.contract?.receipt_code || "—"}</strong>
                </div>

                {!selected.contract?.payment_verified && (
                  <div className="form-actions">
                    <button onClick={() => verifyPayment(selected.id)}>
                      Verify Payment
                    </button>
                  </div>
                )}
              </section>
              <section className="panel">
                <h3>Abyssal Modules</h3>

                {(selected.abyssal_modules ?? []).length === 0 ? (
                  <p className="muted">No abyssal modules declared.</p>
                ) : (
                  selected.abyssal_modules.map((module) => (
                    <div className="module" key={module.name}>
                      <div>
                        <strong>{module.name}</strong>

                        <span>
                          Value: {module.value_isk.toLocaleString()} ISK
                        </span>

                        <span>
                          Deductible: {module.deductible_isk.toLocaleString()}{" "}
                          ISK
                        </span>
                      </div>

                      <span
                        className={
                          module.deductible_paid ? "success" : "failure"
                        }
                      >
                        {module.deductible_paid ? "PAID" : "UNPAID"}
                      </span>
                    </div>
                  ))
                )}
              </section>
              <section className="panel">
                <div className="panel-header">
                  <h3>Asset Scan</h3>

                  <button
                    onClick={() => runScan(selected.id)}
                    disabled={scanning}
                  >
                    {scanning ? "Scanning..." : "Run Asset Scan"}
                  </button>
                </div>

                {selected.asset_scan?.passed ? (
                  <div className="scan-result passed">
                    <div className="check">✓</div>

                    <div>
                      <strong>Scan Passed</strong>
                      <span>Ship contents match the expected assets.</span>
                    </div>
                  </div>
                ) : (
                  <div className="scan-result failed">
                    <div className="check">✕</div>

                    <div>
                      <strong>Scan Failed</strong>
                      <span>Expected assets were not detected.</span>
                    </div>
                  </div>
                )}

                <div className="scan-columns">
                  <div>
                    <h4>Expected</h4>

                    {(selected.asset_scan?.expected_modules ?? []).length ===
                    0 ? (
                      <span className="muted">No modules expected</span>
                    ) : (
                      selected.asset_scan.expected_modules.map((module) => (
                        <div key={module}>✓ {module}</div>
                      ))
                    )}
                  </div>

                  <div>
                    <h4>Detected</h4>

                    {(selected.asset_scan?.actual_modules ?? []).length ===
                    0 ? (
                      <span className="muted">No modules detected</span>
                    ) : (
                      selected.asset_scan.actual_modules.map((module) => (
                        <div key={module}>✓ {module}</div>
                      ))
                    )}
                  </div>
                </div>
              </section>
              <section className="panel">
                <h3>Logistics</h3>

                <div className="row">
                  <span>Route</span>
                  <strong>
                    {selected.origin} → {selected.destination}
                  </strong>
                </div>

                <div className="row">
                  <span>Carrier</span>
                  <strong>{selected.carrier_name || "Unassigned"}</strong>
                </div>

                {!selected.carrier_id && (
                  <div className="form-actions">
                    <select
                      value={carrierId}
                      onChange={(event) => setCarrierId(event.target.value)}
                    >
                      <option value="">Select carrier</option>

                      {carriers.map((carrier) => (
                        <option key={carrier.id} value={carrier.id}>
                          {carrier.name}
                        </option>
                      ))}
                    </select>

                    <button
                      disabled={!carrierId}
                      onClick={() => assignCarrier(selected.id)}
                    >
                      Assign Carrier
                    </button>
                  </div>
                )}

                <div className="row">
                  <span>Cyno Pilot</span>
                  <strong>{selected.cyno_pilot_name || "Unassigned"}</strong>
                </div>

                {!selected.cyno_pilot_id && (
                  <div className="form-actions">
                    <select
                      value={cynoPilotId}
                      onChange={(event) => setCynoPilotId(event.target.value)}
                    >
                      <option value="">Select cyno pilot</option>

                      {cynoPilots.map((pilot) => (
                        <option key={pilot.id} value={pilot.id}>
                          {pilot.name}
                        </option>
                      ))}
                    </select>

                    <button
                      disabled={!cynoPilotId}
                      onClick={() => assignCyno(selected.id)}
                    >
                      Assign Cyno
                    </button>
                  </div>
                )}

                {selected.status === "READY" && (
                  <div className="dispatch-section">
                    <button
                      className="dispatch-button"
                      disabled={
                        dispatching ||
                        !selected.asset_scan?.passed ||
                        !selected.carrier_id ||
                        !selected.cyno_pilot_id
                      }
                      onClick={() => dispatchShipment(selected.id)}
                    >
                      {dispatching ? "Dispatching..." : "Dispatch Shipment"}
                    </button>

                    {!selected.asset_scan?.passed && (
                      <span className="dispatch-hint">
                        Asset scan must pass
                      </span>
                    )}

                    {selected.asset_scan?.passed && !selected.carrier_id && (
                      <span className="dispatch-hint">
                        Carrier must be assigned
                      </span>
                    )}

                    {selected.asset_scan?.passed &&
                      selected.carrier_id &&
                      !selected.cyno_pilot_id && (
                        <span className="dispatch-hint">
                          Cyno pilot must be assigned
                        </span>
                      )}
                  </div>
                )}
                {selected.status === "IN_TRANSIT" && (
                  <div className="dispatch-section">
                    <button
                      className="dispatch-button"
                      disabled={delivering}
                      onClick={() => deliverShipment(selected.id)}
                    >
                      {delivering ? "Delivering..." : "Mark Delivered"}
                    </button>

                    <span className="dispatch-hint">
                      Confirm that the shipment reached its destination
                    </span>
                  </div>
                )}
              </section>
              <section className="panel activity-panel">
                <h3>Activity</h3>

                {!selected.activities?.length ? (
                  <p className="muted">No activity yet.</p>
                ) : (
                  <div className="activity-list">
                    {[...selected.activities]
                      .reverse()
                      .map((activity, index) => (
                        <div
                          className="activity-item"
                          key={`${activity.timestamp}-${index}`}
                        >
                          <div className="activity-dot" />

                          <div className="activity-content">
                            <strong>{activity.action}</strong>

                            {activity.details && <p>{activity.details}</p>}

                            <span>
                              {new Date(activity.timestamp).toLocaleTimeString(
                                [],
                                {
                                  hour: "2-digit",
                                  minute: "2-digit",
                                },
                              )}
                            </span>
                          </div>
                        </div>
                      ))}
                  </div>
                )}
              </section>
            </div>
          </>
        ) : (
          <>
            <div className="page-header">
              <div>
                <h2>Operations</h2>
                <p>Monitor current logistics operations.</p>
              </div>

              <button onClick={() => setShowNewShipment(true)}>
                + New Shipment
              </button>
            </div>

            {showNewShipment && (
              <section className="panel new-shipment">
                <div className="panel-header">
                  <h3>New Shipment</h3>

                  <button
                    className="secondary"
                    onClick={() => setShowNewShipment(false)}
                  >
                    Cancel
                  </button>
                </div>

                <form onSubmit={createShipment}>
                  <div className="form-grid">
                    <label>
                      Ship
                      <input
                        value={newShipment.ship}
                        onChange={(event) =>
                          updateNewShipment("ship", event.target.value)
                        }
                        placeholder="e.g. Ishtar"
                        required
                      />
                    </label>

                    <label>
                      Origin
                      <input
                        value={newShipment.origin}
                        onChange={(event) =>
                          updateNewShipment("origin", event.target.value)
                        }
                        placeholder="e.g. Jita"
                        required
                      />
                    </label>

                    <label>
                      Destination
                      <input
                        value={newShipment.destination}
                        onChange={(event) =>
                          updateNewShipment("destination", event.target.value)
                        }
                        placeholder="e.g. Amarr"
                        required
                      />
                    </label>

                    <label>
                      Contract ID
                      <input
                        value={newShipment.contract_id}
                        onChange={(event) =>
                          updateNewShipment("contract_id", event.target.value)
                        }
                        placeholder="CONTRACT-XXXX"
                        required
                      />
                    </label>

                    <label>
                      Receipt Code
                      <input
                        value={newShipment.receipt_code}
                        onChange={(event) =>
                          updateNewShipment("receipt_code", event.target.value)
                        }
                        placeholder="SR-XXXX"
                        required
                      />
                    </label>
                  </div>

                  <div className="form-actions">
                    <button type="submit">Create Shipment</button>
                  </div>
                </form>
              </section>
            )}

            <div className="stats-grid">
              <div className="stat-card">
                <span>READY</span>
                <strong>{readyCount}</strong>
              </div>

              <div className="stat-card">
                <span>BLOCKED</span>
                <strong>{blockedCount}</strong>
              </div>

              <div className="stat-card">
                <span>IN TRANSIT</span>
                <strong>{inTransitCount}</strong>
              </div>

              <div className="stat-card">
                <span>DELIVERED</span>
                <strong>{deliveredCount}</strong>
              </div>
            </div>

            <div className="section-header">
              <div>
                <h2>Active Shipments</h2>
                <p>Monitor current logistics operations.</p>
              </div>
            </div>

            <div className="operations-list">
              {shipments.map((shipment) => (
                <article
                  key={shipment.id}
                  className="operation-card"
                  onClick={() => openShipment(shipment.id)}
                >
                  <div className="operation-header">
                    <div>
                      <span className="shipment-id">#{shipment.id}</span>

                      <h3>{shipment.ship}</h3>
                    </div>

                    <span className={`status ${statusClass(shipment.status)}`}>
                      {statusLabel(shipment.status)}
                    </span>
                  </div>

                  <div className="route">
                    <strong>{shipment.origin}</strong>
                    <span>→</span>
                    <strong>{shipment.destination}</strong>
                  </div>

                  <div className="operation-details">
                    <div>
                      <span>Carrier</span>
                      <strong>{shipment.carrier_name || "Unassigned"}</strong>
                    </div>

                    <div>
                      <span>Cyno</span>
                      <strong>
                        {shipment.cyno_pilot_name || "Unassigned"}
                      </strong>
                    </div>

                    <div>
                      <span>Payment</span>
                      <strong>
                        {shipment.contract?.payment_verified
                          ? "✓ Verified"
                          : "✕ Pending"}
                      </strong>
                    </div>

                    <div>
                      <span>Asset Scan</span>
                      <strong
                        className={
                          shipment.asset_scan?.passed ? "success" : "failure"
                        }
                      >
                        {shipment.asset_scan?.passed ? "✓ Passed" : "✕ Failed"}
                      </strong>
                    </div>
                  </div>

                  {shipment.status === "BLOCKED" && (
                    <div className="blocked-message">
                      Asset verification failed — shipment blocked.
                    </div>
                  )}

                  <div className="operation-footer">
                    <span>Open shipment →</span>
                  </div>
                </article>
              ))}
            </div>
          </>
        )}
      </main>
    </div>
  );
}

export default App;
