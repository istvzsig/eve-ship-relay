package pricing

import (
	"fmt"

	"github.com/istvzsig/eve-ship-relay/internal/risk"
	"github.com/istvzsig/eve-ship-relay/internal/ship"
	"github.com/istvzsig/eve-ship-relay/internal/shipment"
)

type Calculator struct {
	coefficients Coefficients
}

func NewCalculator(coefficients Coefficients) *Calculator {
	return &Calculator{
		coefficients: coefficients,
	}
}

func (c *Calculator) Calculate(s shipment.Shipment) (Breakdown, error) {
	if !s.Valid() {
		return Breakdown{}, fmt.Errorf("invalid shipment")
	}

	volumeCost, err := c.VolumeCost(s)
	if err != nil {
		return Breakdown{}, err
	}

	routeRisk, err := c.RouteRisk(s)
	if err != nil {
		return Breakdown{}, err
	}

	gankRisk, err := c.GankRisk(s)
	if err != nil {
		return Breakdown{}, err
	}

	breakdown := Breakdown{
		VolumeCost:     volumeCost,
		DistanceCost:   c.DistanceCost(s),
		RouteRisk:      routeRisk,
		GankRisk:       gankRisk,
		CollateralCost: c.CollateralCost(s),
		JumpFuelCost:   c.JumpFuelCost(s),
		CynoCost:       c.CynoCost(s),
	}

	breakdown.Total =
		breakdown.VolumeCost +
			breakdown.DistanceCost +
			breakdown.CollateralCost +
			breakdown.JumpFuelCost +
			breakdown.CynoCost

	return breakdown, nil
}

func (c *Calculator) VolumeCost(s shipment.Shipment) (float64, error) {
	rate, ok := c.coefficients.VolumeISKPerM3[s.ShipClass]
	if !ok {
		return 0, fmt.Errorf(
			"no volume rate configured for ship class %s",
			s.ShipClass,
		)
	}

	return s.VolumeM3 * rate, nil
}

func (c *Calculator) DistanceCost(s shipment.Shipment) float64 {
	return float64(s.Route.DistanceJumps) *
		c.coefficients.DistanceISKPerJump
}

func (c *Calculator) RiskAssessment(s shipment.Shipment) (risk.RiskAssessment, error) {
	gankExposure, ok := c.coefficients.GankRiskByShip[s.ShipClass]
	if !ok {
		return risk.RiskAssessment{}, fmt.Errorf(
			"no gank risk configured for ship class %s",
			s.ShipClass,
		)
	}

	assessment := risk.Assess(
		s,
		c.coefficients.LowSecRiskPerJump,
		c.coefficients.NullSecRiskPerJump,
		gankExposure,
	)

	if !assessment.Valid() {
		return risk.RiskAssessment{}, fmt.Errorf(
			"invalid risk assessment: %+v",
			assessment,
		)
	}

	return assessment, nil
}

func (c *Calculator) RouteRisk(s shipment.Shipment) (float64, error) {
	routeRiskRate :=
		float64(s.Route.LowSecJumps)*c.coefficients.LowSecRiskPerJump +
			float64(s.Route.NullSecJumps)*c.coefficients.NullSecRiskPerJump

	return routeRiskRate, nil
}

func (c *Calculator) GankRisk(s shipment.Shipment) (float64, error) {
	assessment, err := c.RiskAssessment(s)
	if err != nil {
		return 0, err
	}

	return assessment.GankRisk, nil
}

func (c *Calculator) CollateralCost(s shipment.Shipment) float64 {
	return s.CollateralISK *
		c.coefficients.CollateralRiskRate
}

func (c *Calculator) JumpFuelCost(s shipment.Shipment) float64 {
	if s.ShipClass != ship.JumpFreighter {
		return 0
	}

	return float64(s.Route.DistanceJumps) *
		c.coefficients.JumpFuelISKPerLeg
}

func (c *Calculator) CynoCost(s shipment.Shipment) float64 {
	if s.ShipClass != ship.JumpFreighter {
		return 0
	}

	return float64(s.Route.DistanceJumps) *
		c.coefficients.CynoISKPerLeg
}
