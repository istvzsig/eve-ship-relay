package risk

type RiskAssessment struct {
	PickupRisk      float64
	DestinationRisk float64
	RouteRisk       float64
	GankRisk        float64
	CollateralRisk  float64
}

func (r RiskAssessment) Valid() bool {
	return r.PickupRisk >= 0 && r.PickupRisk <= 1 &&
		r.DestinationRisk >= 0 && r.DestinationRisk <= 1 &&
		r.RouteRisk >= 0 && r.RouteRisk <= 1 &&
		r.GankRisk >= 0 && r.GankRisk <= 1 &&
		r.CollateralRisk >= 0 && r.CollateralRisk <= 1
}
