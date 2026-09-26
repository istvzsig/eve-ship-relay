package shipment

import "time"

type Status string

const (
	StatusReady     Status = "READY"
	StatusInTransit Status = "IN_TRANSIT"
	StatusDelivered Status = "DELIVERED"
	StatusBlocked   Status = "BLOCKED"
)

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

type Shipment struct {
	ID          int    `json:"id"`
	Ship        string `json:"ship"`
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Status      Status `json:"status"`

	CarrierID   string `json:"carrier_id"`
	CarrierName string `json:"carrier_name"`

	Contract       Contract        `json:"contract"`
	AbyssalModules []AbyssalModule `json:"abyssal_modules"`

	AssetScan AssetScan `json:"asset_scan"`

	CynoPilotID   string `json:"cyno_pilot_id"`
	CynoPilotName string `json:"cyno_pilot_name"`

	Activities []Activity `json:"activities"`
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
