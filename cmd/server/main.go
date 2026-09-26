package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/istvzsig/eve-ship-relay/internal/shipment"
)

var shipments = []shipment.Shipment{
	{
		ID:          1042,
		Ship:        "Ishtar",
		Origin:      "Jita",
		Destination: "Amarr",
		Status:      shipment.StatusReady,

		Contract: shipment.Contract{
			ID:              "CONTRACT-1042",
			PaymentVerified: true,
			ReceiptCode:     "SR-7F42",
		},

		AbyssalModules: []shipment.AbyssalModule{
			{
				Name:           "Abyssal Energized Adaptive Nano Membrane",
				ValueISK:       850_000_000,
				Deductible:     170_000_000,
				DeductiblePaid: true,
			},
		},

		AssetScanPassed: true,
	},
	{
		ID:          1041,
		Ship:        "Rifter",
		Origin:      "Jita",
		Destination: "Hek",
		Status:      shipment.StatusInTransit,

		Contract: shipment.Contract{
			ID:              "CONTRACT-1042",
			PaymentVerified: true,
			ReceiptCode:     "SR-7F42",
		},

		AssetScanPassed: true,
	},
	{
		ID:          1040,
		Ship:        "Vargur",
		Origin:      "Jita",
		Destination: "Dodixie",
		Status:      shipment.StatusDelivered,

		Contract: shipment.Contract{
			ID:              "CONTRACT-1040",
			PaymentVerified: true,
			ReceiptCode:     "SR-91BC",
		},

		AssetScanPassed: true,
	},
}

func main() {
	http.HandleFunc("/api/shipments", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(shipments)
	})

	http.HandleFunc("/api/shipments/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		for _, s := range shipments {
			if strconv.Itoa(s.ID) == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(s)
				return
			}
		}

		http.Error(w, "shipment not found", http.StatusNotFound)
	})

	log.Println("ShipRelay listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
