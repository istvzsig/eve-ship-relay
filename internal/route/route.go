package route

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const esiBaseURL = "https://esi.evetech.net"

type Route struct {
	Origin      string   `json:"origin"`
	Destination string   `json:"destination"`
	Systems     []System `json:"systems"`

	DistanceJumps int `json:"distance_jumps"`

	HighSecJumps int `json:"high_sec_jumps"`
	LowSecJumps  int `json:"low_sec_jumps"`
	NullSecJumps int `json:"null_sec_jumps"`
}

type System struct {
	ID             int64    `json:"id"`
	Name           string   `json:"name"`
	SecurityStatus float64  `json:"security_status"`
	Position       Position `json:"position"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

func (r Route) Valid() bool {
	if r.DistanceJumps < 0 {
		return false
	}

	if r.HighSecJumps < 0 ||
		r.LowSecJumps < 0 ||
		r.NullSecJumps < 0 {
		return false
	}

	return r.DistanceJumps ==
		r.HighSecJumps+
			r.LowSecJumps+
			r.NullSecJumps
}

func BuildRoute(
	originName string,
	destinationName string,
	systems []System,
) Route {
	r := Route{
		Origin:      originName,
		Destination: destinationName,
		Systems:     systems,
	}

	if len(systems) <= 1 {
		return r
	}

	r.DistanceJumps = len(systems) - 1

	for _, system := range systems[1:] {
		switch {
		case system.SecurityStatus >= 0.5:
			r.HighSecJumps++
		case system.SecurityStatus > 0:
			r.LowSecJumps++
		default:
			r.NullSecJumps++
		}
	}

	return r
}

type ESIClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewESIClient() *ESIClient {
	return &ESIClient{
		baseURL:    esiBaseURL,
		httpClient: &http.Client{},
	}
}

func (c *ESIClient) SearchSystem(ctx context.Context, name string) (int64, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/universe/ids/",
		strings.NewReader(fmt.Sprintf(`["%s"]`, name)),
	)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Compatibility-Date", "2025-09-30")
	req.Header.Set("User-Agent", "ShipRelay/0.1")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("ESI universe ids: HTTP %d", resp.StatusCode)
	}

	var result struct {
		Systems []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"systems"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	for _, system := range result.Systems {
		if strings.EqualFold(system.Name, name) {
			return system.ID, nil
		}
	}

	return 0, fmt.Errorf("solar system %q not found", name)
}

func (c *ESIClient) GetSystem(ctx context.Context, id int64) (System, error) {
	endpoint := fmt.Sprintf("%s/latest/universe/systems/%d/", c.baseURL, id)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint,
		nil,
	)
	if err != nil {
		return System{}, err
	}

	req.Header.Set("X-Compatibility-Date", "2025-09-30")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return System{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return System{}, fmt.Errorf(
			"ESI system lookup returned %s",
			resp.Status,
		)
	}

	var result struct {
		Name           string   `json:"name"`
		SecurityStatus float64  `json:"security_status"`
		Position       Position `json:"position"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return System{}, err
	}

	return System{
		ID:             id,
		Name:           result.Name,
		SecurityStatus: result.SecurityStatus,
		Position:       result.Position,
	}, nil
}

func (c *ESIClient) CalculateRoute(
	ctx context.Context,
	originID int64,
	destinationID int64,
	preference string,
) ([]int64, error) {
	endpoint := fmt.Sprintf(
		"%s/route/%d/%d/",
		c.baseURL,
		originID,
		destinationID,
	)

	body := fmt.Sprintf(
		`{"preference":%q,"security_penalty":50}`,
		preference,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		strings.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Compatibility-Date", "2025-09-30")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ESI route returned %s", resp.Status)
	}

	var result struct {
		Route []int64 `json:"route"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Route, nil
}

func (c *ESIClient) CalculateRouteByName(
	ctx context.Context,
	originName string,
	destinationName string,
	preference string,
) (Route, error) {
	origin, err := c.SearchSystem(ctx, originName)
	if err != nil {
		return Route{}, err
	}

	destination, err := c.SearchSystem(ctx, destinationName)
	if err != nil {
		return Route{}, err
	}

	ids, err := c.CalculateRoute(
		ctx,
		origin,
		destination,
		preference,
	)
	if err != nil {
		return Route{}, err
	}

	systems := make([]System, 0, len(ids))

	for _, id := range ids {
		system, err := c.GetSystem(ctx, id)
		if err != nil {
			return Route{}, err
		}

		systems = append(systems, system)
	}

	result := BuildRoute(
		originName,
		destinationName,
		systems,
	)

	if !result.Valid() {
		return Route{}, fmt.Errorf("invalid calculated route")
	}

	return result, nil
}
