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

type SolarSystem struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type SystemCatalog struct {
	Systems []SolarSystem `json:"systems"`
}

const systemCatalogFile = "data/systems.json"

func LoadSystems() ([]SolarSystem, error) {
	data, err := os.ReadFile(systemCatalogFile)
	if err != nil {
		return nil, fmt.Errorf("read system catalog: %w", err)
	}

	var catalog SystemCatalog
	if err := json.Unmarshal(data, &catalog); err != nil {
		return nil, fmt.Errorf("decode system catalog: %w", err)
	}

	return catalog.Systems, nil
}

func SearchSystems(query string, limit int) ([]SolarSystem, error) {
	query = strings.TrimSpace(strings.ToLower(query))

	if query == "" {
		return []SolarSystem{}, nil
	}

	systems, err := LoadSystems()
	if err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 10
	}

	var exact []SolarSystem
	var prefix []SolarSystem
	var contains []SolarSystem

	for _, system := range systems {
		name := strings.ToLower(system.Name)

		switch {
		case name == query:
			exact = append(exact, system)
		case strings.HasPrefix(name, query):
			prefix = append(prefix, system)
		case strings.Contains(name, query):
			contains = append(contains, system)
		}
	}

	sort.Slice(prefix, func(i, j int) bool {
		return prefix[i].Name < prefix[j].Name
	})

	sort.Slice(contains, func(i, j int) bool {
		return contains[i].Name < contains[j].Name
	})

	result := make([]SolarSystem, 0, limit)

	appendResults := func(items []SolarSystem) {
		for _, system := range items {
			if len(result) >= limit {
				return
			}

			result = append(result, system)
		}
	}

	appendResults(exact)
	appendResults(prefix)
	appendResults(contains)

	return result, nil
}

func SystemSearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	limit := 10

	if raw := r.URL.Query().Get("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 && parsed <= 50 {
			limit = parsed
		}
	}

	systems, err := SearchSystems(query, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(systems); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
