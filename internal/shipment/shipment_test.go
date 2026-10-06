package shipment

import (
	"testing"

	"github.com/istvzsig/eve-ship-relay/internal/ship"
)

func TestShipmentValid(t *testing.T) {
	tests := []struct {
		name     string
		shipment Shipment
		want     bool
	}{
		{
			name: "valid blockade runner",
			shipment: Shipment{
				Origin:        "Jita",
				Destination:   "Amarr",
				VolumeM3:      100,
				CollateralISK: 1_000_000,
				ShipClass:     ship.BlockadeRunner,
			},
			want: true,
		},
		{
			name: "zero volume",
			shipment: Shipment{
				Origin:      "Jita",
				Destination: "Amarr",
				VolumeM3:    0,
				ShipClass:   ship.BlockadeRunner,
			},
			want: false,
		},
		{
			name: "negative volume",
			shipment: Shipment{
				Origin:      "Jita",
				Destination: "Amarr",
				VolumeM3:    -100,
				ShipClass:   ship.BlockadeRunner,
			},
			want: false,
		},
		{
			name: "negative collateral",
			shipment: Shipment{
				Origin:        "Jita",
				Destination:   "Amarr",
				VolumeM3:      100,
				CollateralISK: -1,
				ShipClass:     ship.BlockadeRunner,
			},
			want: false,
		},
		{
			name: "missing origin",
			shipment: Shipment{
				Destination: "Amarr",
				VolumeM3:    100,
				ShipClass:   ship.BlockadeRunner,
			},
			want: false,
		},
		{
			name: "missing destination",
			shipment: Shipment{
				Origin:    "Jita",
				VolumeM3:  100,
				ShipClass: ship.BlockadeRunner,
			},
			want: false,
		},
		{
			name: "unknown ship class",
			shipment: Shipment{
				Origin:      "Jita",
				Destination: "Amarr",
				VolumeM3:    100,
				ShipClass:   ship.ShipClass("unknown"),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.shipment.Valid(); got != tt.want {
				t.Fatalf("Valid(): got %v, want %v", got, tt.want)
			}
		})
	}
}
