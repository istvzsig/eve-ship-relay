package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/istvzsig/eve-ship-relay/internal/shipment"
)

var shipments = []shipment.Shipment{
	{
		ID:              1042,
		Ship:            "Ishtar",
		Origin:          "Jita",
		Destination:     "Amarr",
		Status:          shipment.StatusReady,
		PaymentVerified: true,
		AssetScanPassed: true,
	},
	{
		ID:              1041,
		Ship:            "Rifter",
		Origin:          "Jita",
		Destination:     "Hek",
		Status:          shipment.StatusInTransit,
		PaymentVerified: true,
		AssetScanPassed: true,
	},
	{
		ID:              1040,
		Ship:            "Vargur",
		Origin:          "Jita",
		Destination:     "Dodixie",
		Status:          shipment.StatusDelivered,
		PaymentVerified: true,
		AssetScanPassed: true,
	},
}

func main() {
	http.HandleFunc("/api/shipments", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(shipments)
	})

	log.Println("ShipRelay listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
