package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

const shipCatalogFile = "data/ships.json"

type Ship struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ShipCatalog struct {
	Ships []Ship `json:"ships"`
}

func LoadShips() ([]Ship, error) {
	data, err := os.ReadFile(shipCatalogFile)
	if err != nil {
		return nil, fmt.Errorf("read ship catalog: %w", err)
	}

	var catalog ShipCatalog

	if err := json.Unmarshal(data, &catalog); err != nil {
		return nil, fmt.Errorf("decode ship catalog: %w", err)
	}

	return catalog.Ships, nil
}

func SearchShips(query string, limit int) ([]Ship, error) {
	query = strings.TrimSpace(strings.ToLower(query))

	if query == "" {
		return []Ship{}, nil
	}

	ships, err := LoadShips()
	if err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 10
	}

	var exact []Ship
	var prefix []Ship
	var contains []Ship

	for _, ship := range ships {
		name := strings.ToLower(ship.Name)

		switch {
		case name == query:
			exact = append(exact, ship)

		case strings.HasPrefix(name, query):
			prefix = append(prefix, ship)

		case strings.Contains(name, query):
			contains = append(contains, ship)
		}
	}

	sort.Slice(prefix, func(i, j int) bool {
		return prefix[i].Name < prefix[j].Name
	})

	sort.Slice(contains, func(i, j int) bool {
		return contains[i].Name < contains[j].Name
	})

	result := make([]Ship, 0, limit)

	appendResults := func(items []Ship) {
		for _, ship := range items {
			if len(result) >= limit {
				return
			}

			result = append(result, ship)
		}
	}

	appendResults(exact)
	appendResults(prefix)
	appendResults(contains)

	return result, nil
}

func ShipSearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	limit := 10

	if raw := r.URL.Query().Get("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 && parsed <= 50 {
			limit = parsed
		}
	}

	ships, err := SearchShips(query, limit)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(ships); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
