package risk

import "github.com/istvzsig/eve-ship-relay/internal/shipment"

type RiskAssessment struct {
	RouteRisk       float64
	GankRisk        float64
	CollateralRisk  float64
	PickupRisk      float64
	DestinationRisk float64
}

func Assess(
	s shipment.Shipment,
	lowSecRiskPerJump float64,
	nullSecRiskPerJump float64,
	gankExposure float64,
) RiskAssessment {
	routeRiskRate :=
		float64(s.Route.LowSecJumps)*lowSecRiskPerJump +
			float64(s.Route.NullSecJumps)*nullSecRiskPerJump

	return RiskAssessment{
		RouteRisk:      routeRiskRate,
		GankRisk:       gankExposure,
		CollateralRisk: s.CollateralISK,
	}
}

func (r RiskAssessment) Valid() bool {
	return r.PickupRisk >= 0 && r.PickupRisk <= 1 &&
		r.DestinationRisk >= 0 && r.DestinationRisk <= 1 &&
		r.RouteRisk >= 0 &&
		r.GankRisk >= 0 && r.GankRisk <= 1 &&
		r.CollateralRisk >= 0
}
