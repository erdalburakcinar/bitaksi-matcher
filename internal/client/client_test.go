package client

import (
	"bitaksi-go-matcher/internal/models"
	"encoding/json"
	"github.com/eapache/go-resiliency/breaker"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// mockDriverSearchResponse is used to create JSON response bodies in tests.
func mockDriverSearchResponse(t *testing.T, driver *models.DriverWithDistance) string {
	t.Helper()
	data, err := json.Marshal(driver)
	if err != nil {
		t.Fatalf("failed to marshal mock driver response: %v", err)
	}
	return string(data)
}

func TestSearchDriver(t *testing.T) {
	tests := []struct {
		name            string
		serverHandler   http.HandlerFunc
		wantErr         bool
		wantErrContains string
		wantDriver      *models.DriverWithDistance
		statusCode      int
		failTimes       int // how many times the server should fail before success
	}{
		{
			name: "Success - 200 OK",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				// A normal successful response
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(mockDriverSearchResponse(t, &models.DriverWithDistance{
					Distance: 100.0,
				})))
			},
			wantErr: false,
			wantDriver: &models.DriverWithDistance{
				Distance: 100.0,
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Not Found - 404",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			wantErr:         true,
			wantErrContains: "no driver found",
			statusCode:      http.StatusNotFound,
		},
		{
			name: "Internal Server Error - 500",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr:         true,
			wantErrContains: "driver service error: received status 500 Internal Server Error",
			statusCode:      http.StatusInternalServerError,
		},
		{
			name: "Invalid JSON response",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`invalid-json`)) // Not valid JSON
			},
			wantErr:         true,
			wantErrContains: "failed to decode JSON",
			statusCode:      http.StatusOK,
		},
		{
			name: "Circuit breaker opens after repeated failure",
			// failTimes set higher than breaker threshold triggers open breaker
			failTimes:       4,
			wantErr:         true,
			wantErrContains: "driver service error: received status 500 Internal Server Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var failureCount int

			// Create a test server that can simulate failures or normal responses
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.failTimes > 0 && failureCount < tt.failTimes {
					w.WriteHeader(http.StatusInternalServerError)
					failureCount++
					return
				}
				// If a custom handler is specified, use it; otherwise, default to 200 OK
				if tt.serverHandler != nil {
					tt.serverHandler(w, r)
				} else {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(mockDriverSearchResponse(t, &models.DriverWithDistance{
						Distance: 42.0,
					})))
				}
			}))
			defer testServer.Close()

			// Initialize a new client with the test server's URL
			client := NewDriverAPI(testServer.URL, "test-api-key")

			// We lower the breaker threshold for quick demonstration
			// (3 failures allowed, then open).
			client.Breaker = NewTestBreaker(3, 1, time.Second)

			// Attempt the search
			respDriver, err := client.SearchDriver(40.0, 29.0, 1000)

			if tt.wantErr && err == nil {
				t.Fatalf("expected an error but got none")
			} else if !tt.wantErr && err != nil {
				t.Fatalf("did not expect an error but got: %v", err)
			}

			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.wantErrContains) {
				t.Errorf("expected error to contain %q, got %q", tt.wantErrContains, err.Error())
			}

			if !tt.wantErr && tt.wantDriver != nil {
				if respDriver.Distance != tt.wantDriver.Distance {
					t.Errorf("expected driver distance = %v, got %v",
						tt.wantDriver.Distance, respDriver.Distance)
				}
			}
		})
	}
}

// NewTestBreaker is a helper to create a circuit breaker with test-friendly settings.
func NewTestBreaker(failuresAllowed, successesNeeded int, timeout time.Duration) *breaker.Breaker {
	return breaker.New(failuresAllowed, successesNeeded, timeout)
}
