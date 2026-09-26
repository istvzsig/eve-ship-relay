# EVE Ship Relay

A private logistics platform for EVE Online ship-shipping operations.

The goal is to automate and manage:

- shipping contract intake
- payment verification
- ship and asset verification
- abyssal module valuation and deductibles
- shipment tracking
- carrier and cyno logistics

## Status

Early prototype.

Currently includes:

- in-memory shipment data
- shipment listing API
- basic shipment statuses

## Development

Run the server:

```bash
go run ./cmd/server
```

## API

The API is available at:

```text
http://localhost:8080
```
