import { useMemo, useState } from "react";
import "./RouteMap.css";

const demoSystems = [
  { name: "Jita", security: 0.9, x: 8, y: 48, type: "origin" },
  { name: "Perimeter", security: 0.9, x: 25, y: 36 },
  { name: "Sivala", security: 0.9, x: 43, y: 45 },
  { name: "Niyabainen", security: 0.9, x: 62, y: 32 },
  { name: "Amarr", security: 1.0, x: 88, y: 48, type: "destination" },
];

export default function RouteMap({ origin, destination, route = demoSystems }) {
  const [selectedSystem, setSelectedSystem] = useState(null);
  const [zoom, setZoom] = useState(1);

  const routeText = useMemo(
    () => route.map((system) => system.name).join(" → "),
    [route],
  );

  const systemList = useMemo(
    () => route.map((system) => system.name).join("\n"),
    [route],
  );

  async function copyRoute() {
    await navigator.clipboard.writeText(systemList);
  }

  async function copyDiscord() {
    const text = [
      "🚀 **ShipRelay Route**",
      "",
      `${origin} → ${destination}`,
      `${route.length - 1} jumps`,
      "",
      routeText,
    ].join("\n");

    await navigator.clipboard.writeText(text);
  }

  return (
    <section className="route-panel">
      <div className="route-panel-header">
        <div>
          <span className="shipment-id">NAVIGATION</span>
          <h3>
            {origin} → {destination}
          </h3>
        </div>

        <div className="route-actions">
          <button onClick={copyRoute}>Copy Route</button>
          <button onClick={copyDiscord}>Copy Discord</button>
        </div>
      </div>

      <div className="route-map-shell">
        <div className="route-map" style={{ transform: `scale(${zoom})` }}>
          <div className="map-grid" />

          <svg
            className="route-lines"
            viewBox="0 0 100 100"
            preserveAspectRatio="none"
          >
            {route.slice(0, -1).map((system, index) => {
              const next = route[index + 1];

              return (
                <line
                  key={`${system.name}-${next.name}`}
                  x1={system.x}
                  y1={system.y}
                  x2={next.x}
                  y2={next.y}
                />
              );
            })}
          </svg>

          {route.map((system, index) => (
            <button
              key={system.name}
              className={`system-node ${
                system.type || ""
              } ${selectedSystem === system.name ? "selected" : ""}`}
              style={{
                left: `${system.x}%`,
                top: `${system.y}%`,
              }}
              onClick={() => setSelectedSystem(system.name)}
              title={system.name}
            >
              <span className="system-dot" />
              <span className="system-label">{system.name}</span>
              <span className="system-security">
                {system.security.toFixed(1)}
              </span>
              <span className="system-index">{index + 1}</span>
            </button>
          ))}
        </div>

        <div className="map-controls">
          <button
            onClick={() => setZoom((value) => Math.min(1.5, value + 0.1))}
          >
            +
          </button>

          <button
            onClick={() => setZoom((value) => Math.max(0.7, value - 0.1))}
          >
            −
          </button>

          <button onClick={() => setZoom(1)}>Reset</button>
        </div>

        {selectedSystem && (
          <div className="system-inspector">
            <span>SYSTEM</span>
            <strong>{selectedSystem}</strong>
          </div>
        )}
      </div>

      <div className="route-summary">
        <div>
          <span>ROUTE</span>
          <strong>
            {origin} → {destination}
          </strong>
        </div>

        <div>
          <span>JUMPS</span>
          <strong>{route.length - 1}</strong>
        </div>

        <div>
          <span>PREFERENCE</span>
          <strong>Shorter</strong>
        </div>
      </div>

      <div className="route-waypoints">
        <div className="route-waypoints-header">
          <h4>Waypoints</h4>
          <span>{route.length} systems</span>
        </div>

        {route.map((system, index) => (
          <div className="waypoint" key={system.name}>
            <span className="waypoint-number">{index + 1}</span>

            <div>
              <strong>{system.name}</strong>
              <span>{system.security.toFixed(1)} security</span>
            </div>

            {index === 0 && <em>ORIGIN</em>}
            {index === route.length - 1 && <em>DESTINATION</em>}
          </div>
        ))}
      </div>
    </section>
  );
}
