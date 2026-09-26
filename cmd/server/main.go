package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/istvzsig/eve-ship-relay/internal/api"
	"github.com/istvzsig/eve-ship-relay/internal/carrier"
	"github.com/istvzsig/eve-ship-relay/internal/shipment"
)

var carriers = []carrier.Carrier{
	{ID: "CARRIER-01", Name: "Night Hauler"},
	{ID: "CARRIER-02", Name: "Red Freighter"},
	{ID: "CARRIER-03", Name: "Void Runner"},
}

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

		AssetScan: shipment.AssetScan{
			ExpectedModules: []string{
				"Abyssal Energized Adaptive Nano Membrane",
			},
			ActualModules: []string{
				"Abyssal Energized Adaptive Nano Membrane",
			},
			Passed: true,
		},
	},

	{
		ID:          1041,
		Ship:        "Rifter",
		Origin:      "Jita",
		Destination: "Hek",
		Status:      shipment.StatusInTransit,

		Contract: shipment.Contract{
			ID:              "CONTRACT-1041",
			PaymentVerified: true,
			ReceiptCode:     "SR-4A21",
		},

		AbyssalModules: []shipment.AbyssalModule{},

		AssetScan: shipment.AssetScan{
			ExpectedModules: []string{},
			ActualModules:   []string{},
			Passed:          true,
		},
	},

	{
		ID:          1040,
		Ship:        "Vargur",
		Origin:      "Jita",
		Destination: "Dodixie",
		Status:      shipment.StatusReady,

		Contract: shipment.Contract{
			ID:              "CONTRACT-1040",
			PaymentVerified: true,
			ReceiptCode:     "SR-91BC",
		},

		AbyssalModules: []shipment.AbyssalModule{
			{
				Name:           "Abyssal Large Armor Repairer",
				ValueISK:       1_200_000_000,
				Deductible:     240_000_000,
				DeductiblePaid: false,
			},
		},

		AssetScan: shipment.AssetScan{
			ExpectedModules: []string{
				"Abyssal Large Armor Repairer",
			},
			ActualModules: []string{},
			Passed:        false,
		},
	},
}

func main() {
	handler := api.CORS(http.DefaultServeMux)

	http.HandleFunc("POST /api/shipments", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Ship        string `json:"ship"`
			Origin      string `json:"origin"`
			Destination string `json:"destination"`
			ContractID  string `json:"contract_id"`
			ReceiptCode string `json:"receipt_code"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		id := 1000
		for _, s := range shipments {
			if s.ID >= id {
				id = s.ID + 1
			}
		}

		newShipment := shipment.Shipment{
			ID:          id,
			Ship:        input.Ship,
			Origin:      input.Origin,
			Destination: input.Destination,
			Status:      shipment.StatusReady,

			Contract: shipment.Contract{
				ID:              input.ContractID,
				PaymentVerified: false,
				ReceiptCode:     input.ReceiptCode,
			},

			AbyssalModules: []shipment.AbyssalModule{},

			AssetScan: shipment.AssetScan{
				ExpectedModules: []string{},
				ActualModules:   []string{},
				Passed:          true,
			},
		}

		shipments = append(shipments, newShipment)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newShipment)
	})

	http.HandleFunc("POST /api/shipments/{id}/verify-payment", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid shipment id", http.StatusBadRequest)
			return
		}

		for i := range shipments {
			if shipments[i].ID != id {
				continue
			}

			shipments[i].Contract.PaymentVerified = true

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(shipments[i])
			return
		}

		http.Error(w, "shipment not found", http.StatusNotFound)
	})

	http.HandleFunc("POST /api/shipments/{id}/assign-carrier", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid shipment id", http.StatusBadRequest)
			return
		}

		var input struct {
			CarrierID string `json:"carrier_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		var carrier *carrier.Carrier

		for i := range carriers {
			if carriers[i].ID == input.CarrierID {
				carrier = &carriers[i]
				break
			}
		}

		if carrier == nil {
			http.Error(w, "carrier not found", http.StatusNotFound)
			return
		}

		for i := range shipments {
			if shipments[i].ID != id {
				continue
			}

			shipments[i].CarrierID = carrier.ID
			shipments[i].CarrierName = carrier.Name

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(shipments[i])
			return
		}

		http.Error(w, "shipment not found", http.StatusNotFound)
	})

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

	http.HandleFunc("/api/shipments/{id}/scan", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1500 * time.Millisecond)

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid shipment id", http.StatusBadRequest)
			return
		}

		for i := range shipments {
			if shipments[i].ID != id {
				continue
			}

			shipments[i].AssetScan.Passed = modulesMatch(
				shipments[i].AssetScan.ExpectedModules,
				shipments[i].AssetScan.ActualModules,
			)

			if !shipments[i].AssetScan.Passed {
				shipments[i].Status = shipment.StatusBlocked
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(shipments[i])
			return
		}

		http.Error(w, "shipment not found", http.StatusNotFound)
	})

	log.Println("ShipRelay listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func modulesMatch(expected, actual []string) bool {
	if len(expected) != len(actual) {
		return false
	}

	expectedSet := make(map[string]struct{}, len(expected))
	for _, module := range expected {
		expectedSet[module] = struct{}{}
	}

	for _, module := range actual {
		if _, ok := expectedSet[module]; !ok {
			return false
		}
	}

	return true
}
