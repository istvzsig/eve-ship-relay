package risk

import (
	"math"
	"testing"

	"github.com/istvzsig/eve-ship-relay/internal/route"
	"github.com/istvzsig/eve-ship-relay/internal/shipment"
)

func TestAssess(t *testing.T) {
	tests := []struct {
		name          string
		lowSecJumps   int
		nullSecJumps  int
		gankExposure  float64
		collateralISK float64
		wantRouteRisk float64
		wantGankRisk  float64
	}{
		{
			name:          "high sec only",
			lowSecJumps:   0,
			nullSecJumps:  0,
			gankExposure:  0.15,
			collateralISK: 100_000_000,
			wantRouteRisk: 0,
			wantGankRisk:  0.15,
		},
		{
			name:          "low sec",
			lowSecJumps:   2,
			nullSecJumps:  0,
			gankExposure:  0.25,
			collateralISK: 500_000_000,
			wantRouteRisk: 0.08,
			wantGankRisk:  0.25,
		},
		{
			name:          "mixed low and null",
			lowSecJumps:   5,
			nullSecJumps:  2,
			gankExposure:  0.60,
			collateralISK: 1_000_000_000,
			wantRouteRisk: 0.36,
			wantGankRisk:  0.60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := shipment.Shipment{
				CollateralISK: tt.collateralISK,
				Route: route.Route{
					LowSecJumps:  tt.lowSecJumps,
					NullSecJumps: tt.nullSecJumps,
				},
			}

			got := Assess(
				input,
				0.04,
				0.08,
				tt.gankExposure,
			)

			if math.Abs(got.RouteRisk-tt.wantRouteRisk) > 1e-9 {
				t.Errorf(
					"RouteRisk: got %v, want %v",
					got.RouteRisk,
					tt.wantRouteRisk,
				)
			}

			if math.Abs(got.GankRisk-tt.wantGankRisk) > 1e-9 {
				t.Errorf(
					"GankRisk: got %v, want %v",
					got.GankRisk,
					tt.wantGankRisk,
				)
			}

			if got.CollateralRisk != tt.collateralISK {
				t.Errorf(
					"CollateralRisk: got %v, want %v",
					got.CollateralRisk,
					tt.collateralISK,
				)
			}
		})
	}
}

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
			name: "negative route risk",
			risk: RiskAssessment{
				RouteRisk: -0.1,
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

func TestAssessZeroRisk(t *testing.T) {
	input := shipment.Shipment{
		CollateralISK: 0,
		Route: route.Route{
			LowSecJumps:  0,
			NullSecJumps: 0,
		},
	}

	got := Assess(input, 0.04, 0.08, 0)

	if got.RouteRisk != 0 {
		t.Errorf("RouteRisk: got %v, want 0", got.RouteRisk)
	}

	if got.GankRisk != 0 {
		t.Errorf("GankRisk: got %v, want 0", got.GankRisk)
	}

	if got.CollateralRisk != 0 {
		t.Errorf("CollateralRisk: got %v, want 0", got.CollateralRisk)
	}
}
