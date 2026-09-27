package assetscanner

import (
	"context"

	"github.com/istvzsig/eve-ship-relay/internal/shipment"
)

type AssetScanner interface {
	Scan(ctx context.Context, shipment shipment.Shipment) (
		shipment.AssetScan,
		error,
	)
}
