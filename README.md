cat > README.md <<'EOF'

# EVE Ship Relay

A prototype logistics operations platform for EVE Online ship-shipping workflows.

ShipRelay models the operational lifecycle of a shipped EVE ship, from contract intake through delivery.

## Current Demo

The prototype currently supports:

- shipment intake
- contract and receipt-code tracking
- payment verification
- abyssal module tracking
- asset verification / scanning
- shipment blocking when assets do not match
- carrier assignment
- cyno pilot assignment
- dispatch validation
- in-transit tracking
- delivery confirmation
- activity history / audit timeline
- dark operator-style React dashboard

## Shipment Lifecycle

READY → IN_TRANSIT → DELIVERED

A shipment cannot be dispatched unless:

- payment is verified
- asset scan passes
- carrier is assigned
- cyno pilot is assigned

Failed asset verification moves the shipment to `BLOCKED`.

## Architecture

- Go HTTP API
- React + Vite frontend
- in-memory demo data
- no database required
- no external EVE Online integrations yet

The current implementation intentionally uses fake data so the operational workflow can be demonstrated without requiring EVE credentials or external services.

## Development

### Backend

```bash
go run ./cmd/server
```

### API

```bash
http://localhost:8080
```

Example:

```bash
curl http://localhost:8080/api/shipments
```

### Frontend

```bash
http://localhost:5173
```

### Production build

```bash
cd frontend
npm run build
```

## Project

Built with Go, React, and Vite.

Module:

```bash
github.com/istvzsig/eve-ship-relay
```

## Status

Early prototype / demo.

The next stage would be replacing the in-memory demo layer with real EVE Online integrations and persistent shipment data.
