package risk

import "testing"

func TestRiskAssessment_Valid(t *testing.T) {
	tests := []struct {
		name string
		risk RiskAssessment
		want bool
	}{
		{
			name: "zero risk",
			risk: RiskAssessment{},
			want: true,
		},
		{
			name: "normal risk",
			risk: RiskAssessment{
				PickupRisk:      0.2,
				DestinationRisk: 0.4,
				RouteRisk:       0.6,
				GankRisk:        0.3,
				CollateralRisk:  0.5,
			},
			want: true,
		},
		{
			name: "maximum risk",
			risk: RiskAssessment{
				PickupRisk:      1,
				DestinationRisk: 1,
				RouteRisk:       1,
				GankRisk:        1,
				CollateralRisk:  1,
			},
			want: true,
		},
		{
			name: "negative pickup risk",
			risk: RiskAssessment{
				PickupRisk: -0.1,
			},
			want: false,
		},
		{
			name: "route risk above maximum",
			risk: RiskAssessment{
				RouteRisk: 1.1,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.risk.Valid(); got != tt.want {
				t.Fatalf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
