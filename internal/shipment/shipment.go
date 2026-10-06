package shipment

import (
	"time"

	"github.com/istvzsig/eve-ship-relay/internal/route"
	"github.com/istvzsig/eve-ship-relay/internal/ship"
)

type Status string

const (
	StatusReady     Status = "READY"
	StatusInTransit Status = "IN_TRANSIT"
	StatusDelivered Status = "DELIVERED"
	StatusBlocked   Status = "BLOCKED"
)

type Shipment struct {
	ID          int         `json:"id"`
	Ship        string      `json:"ship"`
	Origin      string      `json:"origin"`
	Destination string      `json:"destination"`
	Status      Status      `json:"status"`
	Route       route.Route `json:"route"`

	VolumeM3      float64 `json:"volume_m3"`
	CollateralISK float64 `json:"collateral_isk"`

	ShipClass ship.ShipClass `json:"ship_class"`

	CarrierID   string `json:"carrier_id"`
	CarrierName string `json:"carrier_name"`

	Contract       Contract        `json:"contract"`
	AbyssalModules []AbyssalModule `json:"abyssal_modules"`

	AssetScan AssetScan `json:"asset_scan"`

	CynoPilotID   string `json:"cyno_pilot_id"`
	CynoPilotName string `json:"cyno_pilot_name"`

	Activities []Activity `json:"activities"`
}

type Contract struct {
	ID              string `json:"id"`
	PaymentVerified bool   `json:"payment_verified"`
	ReceiptCode     string `json:"receipt_code"`
}

type AbyssalModule struct {
	Name           string `json:"name"`
	ValueISK       int64  `json:"value_isk"`
	Deductible     int64  `json:"deductible_isk"`
	DeductiblePaid bool   `json:"deductible_paid"`
}

type AssetScan struct {
	ExpectedModules []string `json:"expected_modules"`
	ActualModules   []string `json:"actual_modules"`
	Passed          bool     `json:"passed"`
}

type Activity struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Details   string    `json:"details,omitempty"`
}

func (s Shipment) Valid() bool {
	if s.VolumeM3 <= 0 {
		return false
	}

	if s.CollateralISK < 0 {
		return false
	}

	if s.Origin == "" || s.Destination == "" {
		return false
	}

	if !s.Route.Valid() {
		return false
	}

	switch s.ShipClass {
	case ship.BlockadeRunner,
		ship.DST,
		ship.Freighter,
		ship.JumpFreighter:
		return true

	default:
		return false
	}
}
