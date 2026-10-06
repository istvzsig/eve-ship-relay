package pricing

type Breakdown struct {
	VolumeCost     float64 `json:"volume_cost"`
	DistanceCost   float64 `json:"distance_cost"`
	CollateralCost float64 `json:"collateral_cost"`
	JumpFuelCost   float64 `json:"jump_fuel_cost"`
	CynoCost       float64 `json:"cyno_cost"`

	RouteRisk float64 `json:"route_risk"`
	GankRisk  float64 `json:"gank_risk"`

	Total float64 `json:"total"`
}
