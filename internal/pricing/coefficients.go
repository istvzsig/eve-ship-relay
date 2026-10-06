package pricing

import "github.com/istvzsig/eve-ship-relay/internal/ship"

type Coefficients struct {
	// How much the prototype pays per m³ depending on ship.
	VolumeISKPerM3 map[ship.ShipClass]float64

	// Base payment associated with route length.
	DistanceISKPerJump float64

	// Prototype risk contribution from each security jump.
	LowSecRiskPerJump  float64
	NullSecRiskPerJump float64

	// Converts risk exposure into ISK.
	RouteRiskRate      float64
	GankRiskRate       float64
	CollateralRiskRate float64

	// Ship-specific prototype gank exposure.
	GankRiskByShip map[ship.ShipClass]float64

	// Prototype fuel allowance for JF operations.
	JumpFuelISKPerLeg float64
	// Prototype cyno allowance for JF operations.
	CynoISKPerLeg float64
}

func DefaultCoefficients() Coefficients {
	return Coefficients{
		VolumeISKPerM3: map[ship.ShipClass]float64{
			ship.BlockadeRunner: 4000,
			ship.DST:            3000,
			ship.Freighter:      2000,
			ship.JumpFreighter:  5000,
		},

		DistanceISKPerJump: 1_500_000,

		LowSecRiskPerJump:  0.04,
		NullSecRiskPerJump: 0.08,

		RouteRiskRate:      0.20,
		GankRiskRate:       0.20,
		CollateralRiskRate: 0.005,

		GankRiskByShip: map[ship.ShipClass]float64{
			ship.BlockadeRunner: 0.15,
			ship.DST:            0.25,
			ship.Freighter:      0.60,
			ship.JumpFreighter:  0.35,
		},

		JumpFuelISKPerLeg: 1_000_000,
		CynoISKPerLeg:     2_000_000,
	}
}
