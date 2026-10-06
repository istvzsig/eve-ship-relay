package pricing

import (
	"testing"

	"github.com/istvzsig/eve-ship-relay/internal/ship"
)

func TestDefaultCoefficients(t *testing.T) {
	coefficients := DefaultCoefficients()

	tests := []struct {
		name string
		got  float64
		want float64
	}{
		{
			name: "distance per jump",
			got:  coefficients.DistanceISKPerJump,
			want: 1_500_000,
		},
		{
			name: "low sec risk per jump",
			got:  coefficients.LowSecRiskPerJump,
			want: 0.04,
		},
		{
			name: "null sec risk per jump",
			got:  coefficients.NullSecRiskPerJump,
			want: 0.08,
		},
		{
			name: "collateral risk rate",
			got:  coefficients.CollateralRiskRate,
			want: 0.005,
		},
		{
			name: "jump fuel per leg",
			got:  coefficients.JumpFuelISKPerLeg,
			want: 1_000_000,
		},
		{
			name: "cyno per leg",
			got:  coefficients.CynoISKPerLeg,
			want: 2_000_000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestDefaultCoefficients_ShipValues(t *testing.T) {
	coefficients := DefaultCoefficients()

	volumeTests := []struct {
		name string
		ship ship.ShipClass
		want float64
	}{
		{
			name: "blockade runner volume",
			ship: ship.BlockadeRunner,
			want: 4000,
		},
		{
			name: "DST volume",
			ship: ship.DST,
			want: 3000,
		},
		{
			name: "freighter volume",
			ship: ship.Freighter,
			want: 2000,
		},
		{
			name: "jump freighter volume",
			ship: ship.JumpFreighter,
			want: 5000,
		},
	}

	for _, tt := range volumeTests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := coefficients.VolumeISKPerM3[tt.ship]
			if !ok {
				t.Fatalf("missing volume coefficient for %q", tt.ship)
			}

			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}

	gankTests := []struct {
		name string
		ship ship.ShipClass
		want float64
	}{
		{
			name: "blockade runner gank risk",
			ship: ship.BlockadeRunner,
			want: 0.15,
		},
		{
			name: "DST gank risk",
			ship: ship.DST,
			want: 0.25,
		},
		{
			name: "freighter gank risk",
			ship: ship.Freighter,
			want: 0.60,
		},
		{
			name: "jump freighter gank risk",
			ship: ship.JumpFreighter,
			want: 0.35,
		},
	}

	for _, tt := range gankTests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := coefficients.GankRiskByShip[tt.ship]
			if !ok {
				t.Fatalf("missing gank coefficient for %q", tt.ship)
			}

			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
