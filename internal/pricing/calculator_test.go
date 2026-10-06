package pricing

import (
	"testing"

	"github.com/istvzsig/eve-ship-relay/internal/ship"
	"github.com/istvzsig/eve-ship-relay/internal/shipment"
)

func TestVolumeCost(t *testing.T) {
	coefficients := DefaultCoefficients()
	calculator := NewCalculator(coefficients)

	tests := []struct {
		name string
		ship ship.ShipClass
		vol  int
		want float64
	}{
		{
			name: "cost for Blockade Runner",
			ship: ship.BlockadeRunner,
			vol:  100,
			want: 400_000,
		},
		{
			name: "cost for DST",
			ship: ship.DST,
			vol:  100,
			want: 300_000,
		},
		{
			name: "cost for Freighter",
			ship: ship.Freighter,
			vol:  100,
			want: 200_000,
		},
		{
			name: "cost for Jump Freighter",
			ship: ship.JumpFreighter,
			vol:  100,
			want: 500_000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := shipment.Shipment{
				VolumeM3:  float64(tt.vol),
				ShipClass: tt.ship,
			}

			got, err := calculator.VolumeCost(input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
