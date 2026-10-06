package route

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestBuildRoute(t *testing.T) {
	tests := []struct {
		name         string
		systems      []System
		wantDistance int
		wantHighSec  int
		wantLowSec   int
		wantNullSec  int
	}{
		{
			name:         "empty route",
			systems:      nil,
			wantDistance: 0,
			wantHighSec:  0,
			wantLowSec:   0,
			wantNullSec:  0,
		},
		{
			name: "single system",
			systems: []System{
				{ID: 30000142, Name: "Jita", SecurityStatus: 0.95},
			},
			wantDistance: 0,
			wantHighSec:  0,
			wantLowSec:   0,
			wantNullSec:  0,
		},
		{
			name: "high sec route",
			systems: []System{
				{ID: 1, Name: "Jita", SecurityStatus: 0.95},
				{ID: 2, Name: "Perimeter", SecurityStatus: 0.90},
				{ID: 3, Name: "Urlen", SecurityStatus: 0.80},
			},
			wantDistance: 2,
			wantHighSec:  2,
			wantLowSec:   0,
			wantNullSec:  0,
		},
		{
			name: "mixed route",
			systems: []System{
				{ID: 1, Name: "Jita", SecurityStatus: 0.95},
				{ID: 2, Name: "LowSec", SecurityStatus: 0.40},
				{ID: 3, Name: "NullSec", SecurityStatus: 0.00},
				{ID: 4, Name: "DeepNull", SecurityStatus: -0.20},
			},
			wantDistance: 3,
			wantHighSec:  0,
			wantLowSec:   1,
			wantNullSec:  2,
		},
		{
			name: "jita to amarr style mixed route",
			systems: []System{
				{Name: "Jita", SecurityStatus: 0.9},
				{Name: "Ikuchi", SecurityStatus: 0.9},
				{Name: "Ahbazon", SecurityStatus: 0.4},
				{Name: "Shera", SecurityStatus: 0.6},
				{Name: "Amarr", SecurityStatus: 0.9},
			},
			wantDistance: 4,
			wantHighSec:  3,
			wantLowSec:   1,
			wantNullSec:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildRoute("Jita", "Amarr", tt.systems)

			if got.DistanceJumps != tt.wantDistance {
				t.Errorf(
					"DistanceJumps: got %d, want %d",
					got.DistanceJumps,
					tt.wantDistance,
				)
			}

			if got.HighSecJumps != tt.wantHighSec {
				t.Errorf(
					"HighSecJumps: got %d, want %d",
					got.HighSecJumps,
					tt.wantHighSec,
				)
			}

			if got.LowSecJumps != tt.wantLowSec {
				t.Errorf(
					"LowSecJumps: got %d, want %d",
					got.LowSecJumps,
					tt.wantLowSec,
				)
			}

			if got.NullSecJumps != tt.wantNullSec {
				t.Errorf(
					"NullSecJumps: got %d, want %d",
					got.NullSecJumps,
					tt.wantNullSec,
				)
			}
		})
	}
}

func TestRouteValid(t *testing.T) {
	tests := []struct {
		name  string
		route Route
		want  bool
	}{
		{
			name: "valid high sec route",
			route: Route{
				DistanceJumps: 10,
				HighSecJumps:  10,
			},
			want: true,
		},
		{
			name: "valid mixed route",
			route: Route{
				DistanceJumps: 10,
				HighSecJumps:  7,
				LowSecJumps:   2,
				NullSecJumps:  1,
			},
			want: true,
		},
		{
			name: "zero distance",
			route: Route{
				DistanceJumps: 0,
			},
			want: true,
		},
		{
			name: "negative distance",
			route: Route{
				DistanceJumps: -1,
			},
			want: false,
		},
		{
			name: "negative high sec jumps",
			route: Route{
				DistanceJumps: 1,
				HighSecJumps:  -1,
				LowSecJumps:   2,
			},
			want: false,
		},
		{
			name: "jump counts do not match distance",
			route: Route{
				DistanceJumps: 10,
				HighSecJumps:  5,
				LowSecJumps:   2,
				NullSecJumps:  2,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.route.Valid()

			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func newTestESIClient(serverURL string) *ESIClient {
	return &ESIClient{
		baseURL:    serverURL,
		httpClient: &http.Client{},
	}
}

func TestSearchSystem(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method: got %s, want POST", r.Method)
		}

		if r.URL.Path != "/universe/ids/" {
			t.Fatalf("path: got %s, want /universe/ids/", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, `{
			"systems": [
				{
					"id": 30000142,
					"name": "Jita"
				}
			]
		}`)
	}))
	defer server.Close()

	client := newTestESIClient(server.URL)

	got, err := client.SearchSystem(context.Background(), "Jita")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 30000142 {
		t.Fatalf("got %d, want %d", got, 30000142)
	}
}

func TestGetSystem(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method: got %s, want GET", r.Method)
		}

		if r.URL.Path != "/latest/universe/systems/30000142/" {
			t.Fatalf(
				"path: got %s, want /latest/universe/systems/30000142/",
				r.URL.Path,
			)
		}

		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, `{
			"name": "Jita",
			"security_status": 0.945,
			"position": {
				"x": 123.45,
				"y": 678.90,
				"z": -111.22
			}
		}`)
	}))
	defer server.Close()

	client := newTestESIClient(server.URL)

	got, err := client.GetSystem(
		context.Background(),
		30000142,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.ID != 30000142 {
		t.Errorf("ID: got %d, want %d", got.ID, 30000142)
	}

	if got.Name != "Jita" {
		t.Errorf("Name: got %q, want %q", got.Name, "Jita")
	}

	if got.SecurityStatus != 0.945 {
		t.Errorf(
			"SecurityStatus: got %v, want %v",
			got.SecurityStatus,
			0.945,
		)
	}

	if got.Position.X != 123.45 {
		t.Errorf("Position.X: got %v, want %v", got.Position.X, 123.45)
	}

	if got.Position.Y != 678.90 {
		t.Errorf("Position.Y: got %v, want %v", got.Position.Y, 678.90)
	}

	if got.Position.Z != -111.22 {
		t.Errorf("Position.Z: got %v, want %v", got.Position.Z, -111.22)
	}
}

func TestCalculateRoute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method: got %s, want POST", r.Method)
		}

		if r.URL.Path != "/route/30000142/30002187/" {
			t.Fatalf(
				"path: got %s, want /route/30000142/30002187/",
				r.URL.Path,
			)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"route": [30000142, 30000144, 30002187]
		}`)
	}))
	defer server.Close()

	client := newTestESIClient(server.URL)

	got, err := client.CalculateRoute(
		context.Background(),
		30000142,
		30002187,
		"shortest",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []int64{
		30000142,
		30000144,
		30002187,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("route: got %v, want %v", got, want)
	}
}

func TestCalculateRouteByName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/universe/ids/":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{
				"systems": [
					{"id": 30000142, "name": "Jita"},
					{"id": 30002187, "name": "Amarr"}
				]
			}`)

		case "/route/30000142/30002187/":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{
				"route": [30000142, 30000144, 30002187]
			}`)

		case "/latest/universe/systems/30000142/":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{
				"name": "Jita",
				"security_status": 0.945
			}`)

		case "/latest/universe/systems/30000144/":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{
				"name": "Perimeter",
				"security_status": 0.832
			}`)

		case "/latest/universe/systems/30002187/":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{
				"name": "Amarr",
				"security_status": 1.0
			}`)

		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := newTestESIClient(server.URL)

	got, err := client.CalculateRouteByName(
		context.Background(),
		"Jita",
		"Amarr",
		"shortest",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Origin != "Jita" {
		t.Fatalf("origin: got %q, want %q", got.Origin, "Jita")
	}

	if got.Destination != "Amarr" {
		t.Fatalf("destination: got %q, want %q", got.Destination, "Amarr")
	}

	if got.DistanceJumps != 2 {
		t.Fatalf("distance: got %d, want 2", got.DistanceJumps)
	}

	if got.HighSecJumps != 2 {
		t.Fatalf("high-sec: got %d, want 2", got.HighSecJumps)
	}

	if got.LowSecJumps != 0 {
		t.Fatalf("low-sec: got %d, want 0", got.LowSecJumps)
	}

	if got.NullSecJumps != 0 {
		t.Fatalf("null-sec: got %d, want 0", got.NullSecJumps)
	}

	if len(got.Systems) != 3 {
		t.Fatalf("systems: got %d, want 3", len(got.Systems))
	}
}
