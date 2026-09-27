package route

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const esiBaseURL = "https://esi.evetech.net"

type System struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Route struct {
	Origin      string   `json:"origin"`
	Destination string   `json:"destination"`
	Systems     []System `json:"systems"`
}

type ESIClient struct {
	httpClient *http.Client
}

func NewESIClient() *ESIClient {
	return &ESIClient{
		httpClient: &http.Client{},
	}
}

func (c *ESIClient) SearchSystem(ctx context.Context, name string) (int64, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		esiBaseURL+"/universe/ids/",
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
	endpoint := fmt.Sprintf("%s/latest/universe/systems/%d/", esiBaseURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
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
		return System{}, fmt.Errorf("ESI system lookup returned %s", resp.Status)
	}

	var result struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return System{}, err
	}

	return System{
		ID:   id,
		Name: result.Name,
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
		esiBaseURL,
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

	ids, err := c.CalculateRoute(ctx, origin, destination, preference)
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

	return Route{
		Origin:      originName,
		Destination: destinationName,
		Systems:     systems,
	}, nil
}
