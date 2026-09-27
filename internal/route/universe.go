package route

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type Universe struct {
	httpClient *http.Client

	once sync.Once
	ids  map[string]int64
	err  error
}

func NewUniverse(client *http.Client) *Universe {
	if client == nil {
		client = &http.Client{}
	}

	return &Universe{
		httpClient: client,
		ids:        make(map[string]int64),
	}
}

func (u *Universe) SystemID(ctx context.Context, name string) (int64, error) {
	u.once.Do(func() {
		u.err = u.load(ctx)
	})

	if u.err != nil {
		return 0, u.err
	}

	id, ok := u.ids[strings.ToLower(strings.TrimSpace(name))]
	if !ok {
		return 0, fmt.Errorf("solar system %q not found", name)
	}

	return id, nil
}

func (u *Universe) load(ctx context.Context) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://esi.evetech.net/latest/universe/systems/",
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("X-Compatibility-Date", "2025-09-30")

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ESI universe systems returned %s", resp.Status)
	}

	var ids []int64

	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		return err
	}

	for _, id := range ids {
		system, err := getSystem(ctx, u.httpClient, id)
		if err != nil {
			return err
		}

		u.ids[strings.ToLower(system.Name)] = id
	}

	return nil
}

func getSystem(
	ctx context.Context,
	client *http.Client,
	id int64,
) (System, error) {
	endpoint := fmt.Sprintf(
		"https://esi.evetech.net/latest/universe/systems/%d/",
		id,
	)

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

	resp, err := client.Do(req)
	if err != nil {
		return System{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return System{}, fmt.Errorf(
			"ESI system %d returned %s",
			id,
			resp.Status,
		)
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
