import { useEffect, useMemo, useState } from "react";
import "./RouteMap.css";

function positionSystems(systems) {
  if (systems.length === 0) {
    return [];
  }

  const xs = systems.map((system) => system.position.x);
  const zs = systems.map((system) => system.position.z);

  const minX = Math.min(...xs);
  const maxX = Math.max(...xs);
  const minZ = Math.min(...zs);
  const maxZ = Math.max(...zs);

  const rangeX = maxX - minX || 1;
  const rangeZ = maxZ - minZ || 1;

  return systems.map((system, index) => ({
    ...system,
    x: 8 + ((system.position.x - minX) / rangeX) * 84,
    y: 8 + ((maxZ - system.position.z) / rangeZ) * 84,
    type:
      index === 0
        ? "origin"
        : index === systems.length - 1
          ? "destination"
          : "",
  }));
}

export default function RouteMap({ apiURL, origin, destination }) {
  const [route, setRoute] = useState([]);
  const [selectedSystem, setSelectedSystem] = useState(null);
  const [zoom, setZoom] = useState(1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    async function loadRoute() {
      if (!origin || !destination) {
        setRoute([]);
        return;
      }

      setLoading(true);
      setError(null);

      try {
        const response = await fetch(
          `${apiURL}/api/route?origin=${encodeURIComponent(
            origin,
          )}&destination=${encodeURIComponent(destination)}&preference=Shorter`,
        );

        if (!response.ok) {
          throw new Error(`Route request failed: ${response.status}`);
        }

        const data = await response.json();

        setRoute(positionSystems(data.systems));
      } catch (err) {
        setError(err.message);
        setRoute([]);
      } finally {
        setLoading(false);
      }
    }

    loadRoute();
  }, [apiURL, origin, destination]);

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

  if (loading) {
    return (
      <section className="route-panel">
        <div className="route-panel-header">
          <div>
            <span className="shipment-id">NAVIGATION</span>
            <h3>
              {origin} → {destination}
            </h3>
          </div>
        </div>

        <div className="route-map-shell">
          <div className="system-inspector">
            <span>ROUTE</span>
            <strong>Loading ESI route...</strong>
          </div>
        </div>
      </section>
    );
  }

  if (error) {
    return (
      <section className="route-panel">
        <div className="route-panel-header">
          <div>
            <span className="shipment-id">NAVIGATION</span>
            <h3>
              {origin} → {destination}
            </h3>
          </div>
        </div>

        <div className="route-map-shell">
          <div className="system-inspector">
            <span>ROUTE ERROR</span>
            <strong>{error}</strong>
          </div>
        </div>
      </section>
    );
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
                  key={`${system.id}-${next.id}`}
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
              key={system.id}
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
                {system.security_status.toFixed(1)}
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
          <div className="waypoint" key={system.id}>
            <span className="waypoint-number">{index + 1}</span>

            <div>
              <strong>{system.name}</strong>
              <span>{system.security_status.toFixed(1)} security</span>
            </div>

            {index === 0 && <em>ORIGIN</em>}
            {index === route.length - 1 && <em>DESTINATION</em>}
          </div>
        ))}
      </div>
    </section>
  );
}
