package pricing

import (
	"fmt"

	"github.com/istvzsig/eve-ship-relay/internal/route"
	"github.com/istvzsig/eve-ship-relay/internal/shipment"
)

type Calculator struct {
	coefficients Coefficients
}

func NewCalculator(coefficients Coefficients) Calculator {
	return Calculator{
		coefficients: coefficients,
	}
}

// VolumeCost = shipment.VolumeM3 x coefficient for shipment.ShipClass
func (c *Calculator) VolumeCost(shipment shipment.Shipment) (float64, error) {
	rate, ok := c.coefficients.VolumeISKPerM3[shipment.ShipClass]
	if !ok {
		return 0, fmt.Errorf("ship class not found: %q", shipment.ShipClass)
	}

	return shipment.VolumeM3 * rate, nil
}

// DistanceCost = DistanceJumps x DistanceISKPerJump
func (c *Calculator) DistanceCost(route route.Route) (float64, error) {
	if !route.Valid() {
		return 0, fmt.Errorf("invalid route")
	}

	return float64(route.DistanceJumps) * c.coefficients.DistanceISKPerJump, nil
}
