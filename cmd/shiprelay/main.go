package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/istvzsig/eve-ship-relay/internal/api"
	"github.com/istvzsig/eve-ship-relay/internal/carrier"
	"github.com/istvzsig/eve-ship-relay/internal/route"
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

		CarrierID:     "CARRIER-02",
		CarrierName:   "Red Freighter",
		CynoPilotID:   "CYNO-02",
		CynoPilotName: "Nightwatch",

		Activities: []shipment.Activity{
			{
				Timestamp: time.Now().Add(-15 * time.Minute),
				Action:    "Shipment created",
			},
			{
				Timestamp: time.Now().Add(-13 * time.Minute),
				Action:    "Payment verified",
			},
			{
				Timestamp: time.Now().Add(-10 * time.Minute),
				Action:    "Asset scan passed",
			},
			{
				Timestamp: time.Now().Add(-8 * time.Minute),
				Action:    "Carrier assigned",
				Details:   "Red Freighter",
			},
			{
				Timestamp: time.Now().Add(-6 * time.Minute),
				Action:    "Cyno pilot assigned",
				Details:   "Nightwatch",
			},
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

		Activities: []shipment.Activity{
			{
				Timestamp: time.Now().Add(-45 * time.Minute),
				Action:    "Shipment created",
			},
			{
				Timestamp: time.Now().Add(-42 * time.Minute),
				Action:    "Payment verified",
			},
			{
				Timestamp: time.Now().Add(-38 * time.Minute),
				Action:    "Asset scan passed",
			},
			{
				Timestamp: time.Now().Add(-35 * time.Minute),
				Action:    "Shipment dispatched",
			},
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

		Activities: []shipment.Activity{
			{
				Timestamp: time.Now().Add(-20 * time.Minute),
				Action:    "Shipment created",
			},
			{
				Timestamp: time.Now().Add(-18 * time.Minute),
				Action:    "Payment verified",
			},
		},
	},
}

type CynoPilot struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var cynoPilots = []CynoPilot{
	{ID: "CYNO-01", Name: "Dark Angel"},
	{ID: "CYNO-02", Name: "Nightwatch"},
	{ID: "CYNO-03", Name: "Black Lantern"},
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	esiClient := route.NewESIClient()

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

		addActivity(&newShipment, "Shipment created", "")

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

			addActivity(
				&shipments[i],
				"Payment verified",
				"",
			)

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

		var selectedCarrier *carrier.Carrier

		for i := range carriers {
			if carriers[i].ID == input.CarrierID {
				selectedCarrier = &carriers[i]
				break
			}
		}

		if selectedCarrier == nil {
			http.Error(w, "carrier not found", http.StatusNotFound)
			return
		}

		for i := range shipments {
			if shipments[i].ID != id {
				continue
			}

			shipments[i].CarrierID = selectedCarrier.ID
			shipments[i].CarrierName = selectedCarrier.Name

			addActivity(
				&shipments[i],
				"Carrier assigned",
				selectedCarrier.Name,
			)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(shipments[i])
			return
		}

		http.Error(w, "shipment not found", http.StatusNotFound)
	})

	http.HandleFunc("POST /api/shipments/{id}/assign-cyno", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid shipment id", http.StatusBadRequest)
			return
		}

		var input struct {
			CynoPilotID string `json:"cyno_pilot_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		var selectedPilot *CynoPilot

		for i := range cynoPilots {
			if cynoPilots[i].ID == input.CynoPilotID {
				selectedPilot = &cynoPilots[i]
				break
			}
		}

		if selectedPilot == nil {
			http.Error(w, "cyno pilot not found", http.StatusNotFound)
			return
		}

		for i := range shipments {
			if shipments[i].ID != id {
				continue
			}

			shipments[i].CynoPilotID = selectedPilot.ID
			shipments[i].CynoPilotName = selectedPilot.Name

			addActivity(
				&shipments[i],
				"Cyno pilot assigned",
				selectedPilot.Name,
			)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(shipments[i])
			return
		}

		http.Error(w, "shipment not found", http.StatusNotFound)
	})

	http.HandleFunc("POST /api/shipments/{id}/dispatch", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid shipment id", http.StatusBadRequest)
			return
		}

		for i := range shipments {
			if shipments[i].ID != id {
				continue
			}

			current := &shipments[i]

			if current.Status != shipment.StatusReady {
				http.Error(
					w,
					"shipment is not ready for dispatch",
					http.StatusConflict,
				)
				return
			}

			if !current.AssetScan.Passed {
				http.Error(
					w,
					"asset scan failed",
					http.StatusConflict,
				)
				return
			}

			if current.CarrierID == "" {
				http.Error(
					w,
					"carrier not assigned",
					http.StatusConflict,
				)
				return
			}

			if current.CynoPilotID == "" {
				http.Error(
					w,
					"cyno pilot not assigned",
					http.StatusConflict,
				)
				return
			}

			current.Status = shipment.StatusInTransit

			addActivity(
				current,
				"Shipment dispatched",
				"",
			)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(*current)
			return
		}

		http.Error(w, "shipment not found", http.StatusNotFound)
	})

	http.HandleFunc("POST /api/shipments/{id}/deliver", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))

		if err != nil {
			http.Error(w, "invalid shipment id", http.StatusBadRequest)
			return
		}

		for i := range shipments {
			if shipments[i].ID != id {
				continue
			}

			current := &shipments[i]

			if current.Status != shipment.StatusInTransit {
				http.Error(
					w,
					"shipment is not in transit",
					http.StatusConflict,
				)
				return
			}

			current.Status = shipment.StatusDelivered

			addActivity(
				current,
				"Shipment delivered",
				"",
			)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(*current)
			return
		}

		http.Error(w, "shipment not found", http.StatusNotFound)
	})

	http.HandleFunc("GET /api/ships", route.ShipSearchHandler)

	http.HandleFunc("GET /api/systems", route.SystemSearchHandler)

	http.HandleFunc("GET /api/route", func(w http.ResponseWriter, r *http.Request) {
		origin := r.URL.Query().Get("origin")
		destination := r.URL.Query().Get("destination")
		preference := r.URL.Query().Get("preference")

		if origin == "" || destination == "" {
			http.Error(
				w,
				"origin and destination are required",
				http.StatusBadRequest,
			)
			return
		}

		if preference == "" {
			preference = "Shorter"
		}

		switch preference {
		case "Shorter", "Safer", "LessSecure":
		default:
			http.Error(
				w,
				"invalid route preference",
				http.StatusBadRequest,
			)
			return
		}

		ctx := r.Context()

		result, err := esiClient.CalculateRouteByName(
			ctx,
			origin,
			destination,
			preference,
		)
		if err != nil {
			log.Printf("route calculation failed: %v", err)

			http.Error(
				w,
				"route calculation failed",
				http.StatusBadGateway,
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
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

			current := &shipments[i]

			current.AssetScan.Passed = modulesMatch(
				current.AssetScan.ExpectedModules,
				current.AssetScan.ActualModules,
			)

			if current.AssetScan.Passed {
				addActivity(
					current,
					"Asset scan passed",
					"Expected assets are present",
				)
			} else {
				current.Status = shipment.StatusBlocked

				addActivity(
					current,
					"Asset scan failed",
					"Expected assets do not match the ship",
				)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(*current)
			return
		}

		http.Error(w, "shipment not found", http.StatusNotFound)
	})

	log.Println("ShipRelay listening on :8080")
	log.Fatal(http.ListenAndServe(":"+port, handler))
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

func addActivity(s *shipment.Shipment, action, details string) {
	s.Activities = append(s.Activities, shipment.Activity{
		Timestamp: time.Now(),
		Action:    action,
		Details:   details,
	})
}
