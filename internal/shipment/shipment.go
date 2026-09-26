package shipment

type Status string

const (
	StatusReady     Status = "READY"
	StatusInTransit Status = "IN_TRANSIT"
	StatusDelivered Status = "DELIVERED"
	StatusBlocked   Status = "BLOCKED"
)

type Shipment struct {
	ID          int    `json:"id"`
	Ship        string `json:"ship"`
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Status      Status `json:"status"`

	PaymentVerified bool `json:"payment_verified"`
	AssetScanPassed bool `json:"asset_scan_passed"`
}
