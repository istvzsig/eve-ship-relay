package pricing

import (
	"math"
	"testing"

	"github.com/istvzsig/eve-ship-relay/internal/route"
	"github.com/istvzsig/eve-ship-relay/internal/ship"
	"github.com/istvzsig/eve-ship-relay/internal/shipment"
)

func TestCalculator(t *testing.T) {
	coefficients := DefaultCoefficients()
	calculator := NewCalculator(coefficients)

	input := shipment.Shipment{
		ShipClass:     ship.JumpFreighter,
		VolumeM3:      10_000,
		CollateralISK: 1_000_000_000,
		Origin:        "Jita",
		Destination:   "Amarr",
		Route: route.Route{
			Origin:        "Jita",
			Destination:   "Amarr",
			DistanceJumps: 10,
			HighSecJumps:  7,
			LowSecJumps:   2,
			NullSecJumps:  1,
		},
	}

	got, err := calculator.Calculate(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := Breakdown{
		VolumeCost:     50_000_000,
		DistanceCost:   15_000_000,
		RouteRisk:      0.16,
		GankRisk:       0.35,
		CollateralCost: 5_000_000,
		JumpFuelCost:   10_000_000,
		CynoCost:       20_000_000,
		Total:          100_000_000,
	}

	if math.Abs(got.VolumeCost-want.VolumeCost) > 1e-9 {
		t.Errorf("VolumeCost: got %v, want %v", got.VolumeCost, want.VolumeCost)
	}

	if math.Abs(got.DistanceCost-want.DistanceCost) > 1e-9 {
		t.Errorf("DistanceCost: got %v, want %v", got.DistanceCost, want.DistanceCost)
	}

	if math.Abs(got.RouteRisk-want.RouteRisk) > 1e-9 {
		t.Errorf("RouteRisk: got %v, want %v", got.RouteRisk, want.RouteRisk)
	}

	if math.Abs(got.GankRisk-want.GankRisk) > 1e-9 {
		t.Errorf("GankRisk: got %v, want %v", got.GankRisk, want.GankRisk)
	}

	if math.Abs(got.CollateralCost-want.CollateralCost) > 1e-9 {
		t.Errorf("CollateralCost: got %v, want %v", got.CollateralCost, want.CollateralCost)
	}

	if math.Abs(got.JumpFuelCost-want.JumpFuelCost) > 1e-9 {
		t.Errorf("JumpFuelCost: got %v, want %v", got.JumpFuelCost, want.JumpFuelCost)
	}

	if math.Abs(got.CynoCost-want.CynoCost) > 1e-9 {
		t.Errorf("CynoCost: got %v, want %v", got.CynoCost, want.CynoCost)
	}

	if math.Abs(got.Total-want.Total) > 1e-9 {
		t.Errorf("Total: got %v, want %v", got.Total, want.Total)
	}
}

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

func TestDistanceCost(t *testing.T) {
	coefficients := DefaultCoefficients()
	calculator := NewCalculator(coefficients)

	tests := []struct {
		name  string
		jumps int
		want  float64
	}{
		{
			name:  "zero jumps",
			jumps: 0,
			want:  0,
		},
		{
			name:  "one jump",
			jumps: 1,
			want:  1_500_000,
		},
		{
			name:  "ten jumps",
			jumps: 10,
			want:  15_000_000,
		},
		{
			name:  "twenty jumps",
			jumps: 20,
			want:  30_000_000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := shipment.Shipment{
				Route: route.Route{
					DistanceJumps: tt.jumps,
				},
			}

			got := calculator.DistanceCost(input)
			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteRisk(t *testing.T) {
	coefficients := DefaultCoefficients()
	calculator := NewCalculator(coefficients)

	tests := []struct {
		name         string
		lowSecJumps  int
		nullSecJumps int
		want         float64
	}{
		{
			name:         "high sec only",
			lowSecJumps:  0,
			nullSecJumps: 0,
			want:         0,
		},
		{
			name:         "one low sec jump",
			lowSecJumps:  1,
			nullSecJumps: 0,
			want:         0.04,
		},
		{
			name:         "one null sec jump",
			lowSecJumps:  0,
			nullSecJumps: 1,
			want:         0.08,
		},
		{
			name:         "five low sec and two null sec jumps",
			lowSecJumps:  5,
			nullSecJumps: 2,
			want:         0.36,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := shipment.Shipment{
				Route: route.Route{
					LowSecJumps:  tt.lowSecJumps,
					NullSecJumps: tt.nullSecJumps,
				},
			}

			got, err := calculator.RouteRisk(input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGankRisk(t *testing.T) {
	coefficients := DefaultCoefficients()
	calculator := NewCalculator(coefficients)

	tests := []struct {
		name      string
		shipClass ship.ShipClass
		want      float64
	}{
		{
			name:      "blockade runner",
			shipClass: ship.BlockadeRunner,
			want:      0.15,
		},
		{
			name:      "DST",
			shipClass: ship.DST,
			want:      0.25,
		},
		{
			name:      "freighter",
			shipClass: ship.Freighter,
			want:      0.60,
		},
		{
			name:      "jump freighter",
			shipClass: ship.JumpFreighter,
			want:      0.35,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := shipment.Shipment{
				ShipClass: tt.shipClass,
			}

			got, err := calculator.GankRisk(input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if math.Abs(got-tt.want) > 1e-9 {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollateralCost(t *testing.T) {
	coefficients := DefaultCoefficients()
	calculator := NewCalculator(coefficients)

	tests := []struct {
		name          string
		collateralISK float64
		want          float64
	}{
		{
			name:          "zero collateral",
			collateralISK: 0,
			want:          0,
		},
		{
			name:          "100 million collateral",
			collateralISK: 100_000_000,
			want:          500_000,
		},
		{
			name:          "1 billion collateral",
			collateralISK: 1_000_000_000,
			want:          5_000_000,
		},
		{
			name:          "5 billion collateral",
			collateralISK: 5_000_000_000,
			want:          25_000_000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := shipment.Shipment{
				CollateralISK: tt.collateralISK,
			}

			got := calculator.CollateralCost(input)
			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJumpFuelCost(t *testing.T) {
	coefficients := DefaultCoefficients()
	calculator := NewCalculator(coefficients)

	tests := []struct {
		name          string
		shipClass     ship.ShipClass
		distanceJumps int
		want          float64
	}{
		{
			name:          "blockade runner",
			shipClass:     ship.BlockadeRunner,
			distanceJumps: 10,
			want:          0,
		},
		{
			name:          "DST",
			shipClass:     ship.DST,
			distanceJumps: 10,
			want:          0,
		},
		{
			name:          "freighter",
			shipClass:     ship.Freighter,
			distanceJumps: 10,
			want:          0,
		},
		{
			name:          "jump freighter one jump",
			shipClass:     ship.JumpFreighter,
			distanceJumps: 1,
			want:          1_000_000,
		},
		{
			name:          "jump freighter ten jumps",
			shipClass:     ship.JumpFreighter,
			distanceJumps: 10,
			want:          10_000_000,
		},
		{
			name:          "jump freighter zero jumps",
			shipClass:     ship.JumpFreighter,
			distanceJumps: 0,
			want:          0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := shipment.Shipment{
				ShipClass: tt.shipClass,
				Route: route.Route{
					DistanceJumps: tt.distanceJumps,
				},
			}

			got := calculator.JumpFuelCost(input)
			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCynoCost(t *testing.T) {
	coefficients := DefaultCoefficients()
	calculator := NewCalculator(coefficients)

	tests := []struct {
		name          string
		shipClass     ship.ShipClass
		distanceJumps int
		want          float64
	}{
		{
			name:          "blockade runner",
			shipClass:     ship.BlockadeRunner,
			distanceJumps: 10,
			want:          0,
		},
		{
			name:          "DST",
			shipClass:     ship.DST,
			distanceJumps: 10,
			want:          0,
		},
		{
			name:          "freighter",
			shipClass:     ship.Freighter,
			distanceJumps: 10,
			want:          0,
		},
		{
			name:          "jump freighter one jump",
			shipClass:     ship.JumpFreighter,
			distanceJumps: 1,
			want:          2_000_000,
		},
		{
			name:          "jump freighter ten jumps",
			shipClass:     ship.JumpFreighter,
			distanceJumps: 10,
			want:          20_000_000,
		},
		{
			name:          "jump freighter zero jumps",
			shipClass:     ship.JumpFreighter,
			distanceJumps: 0,
			want:          0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := shipment.Shipment{
				ShipClass: tt.shipClass,
				Route: route.Route{
					DistanceJumps: tt.distanceJumps,
				},
			}

			got := calculator.CynoCost(input)
			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateInvalidShipment(t *testing.T) {
	coefficients := DefaultCoefficients()
	calculator := NewCalculator(coefficients)

	tests := []struct {
		name  string
		input shipment.Shipment
	}{
		{
			name: "zero volume",
			input: shipment.Shipment{
				ShipClass:   ship.DST,
				VolumeM3:    0,
				Origin:      "Jita",
				Destination: "Amarr",
			},
		},
		{
			name: "negative collateral",
			input: shipment.Shipment{
				ShipClass:     ship.DST,
				VolumeM3:      100,
				CollateralISK: -1,
				Origin:        "Jita",
				Destination:   "Amarr",
			},
		},
		{
			name: "missing origin",
			input: shipment.Shipment{
				ShipClass:   ship.DST,
				VolumeM3:    100,
				Origin:      "",
				Destination: "Amarr",
			},
		},
		{
			name: "missing destination",
			input: shipment.Shipment{
				ShipClass:   ship.DST,
				VolumeM3:    100,
				Origin:      "Jita",
				Destination: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := calculator.Calculate(tt.input)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestVolumeCostUnknownShip(t *testing.T) {
	coefficients := DefaultCoefficients()
	calculator := NewCalculator(coefficients)

	input := shipment.Shipment{
		ShipClass: ship.ShipClass("unknown"),
		VolumeM3:  100,
	}

	_, err := calculator.VolumeCost(input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestJumpFreighterCosts_NonJumpFreighter(t *testing.T) {
	calc := NewCalculator(DefaultCoefficients())

	s := shipment.Shipment{
		ShipClass:     ship.DST,
		VolumeM3:      1000,
		CollateralISK: 1_000_000,
		Origin:        "Jita",
		Destination:   "Amarr",
		Route: route.Route{
			DistanceJumps: 10,
			HighSecJumps:  10,
		},
	}

	if got := calc.JumpFuelCost(s); got != 0 {
		t.Fatalf("JumpFuelCost(): got %v, want 0", got)
	}

	if got := calc.CynoCost(s); got != 0 {
		t.Fatalf("CynoCost(): got %v, want 0", got)
	}
}

func TestVolumeCost_UnknownShipClass(t *testing.T) {
	calc := NewCalculator(DefaultCoefficients())

	s := shipment.Shipment{
		ShipClass:   ship.ShipClass("unknown"),
		VolumeM3:    1000,
		Origin:      "Jita",
		Destination: "Amarr",
	}

	_, err := calc.VolumeCost(s)
	if err == nil {
		t.Fatal("VolumeCost(): expected error for unknown ship class")
	}
}
