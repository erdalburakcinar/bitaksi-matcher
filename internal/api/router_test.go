package api

import (
	"bitaksi-go-matcher/internal/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// mockMatcherHandler is a stubbed-out MatcherHandler for testing router logic.
// We won't test the actual handler logic here since that's covered in handler tests.
type mockMatcherHandler struct{}

func (mh mockMatcherHandler) MatchDriver(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message": "mock MatchDriver called"}`))
}

func TestSetupRouter(t *testing.T) {
	// Create a mock config; the JWT secret could be anything here.
	cfg := &config.Config{
		JWTSecretKey: "secret-token-api-key",
	}

	// Initialize our router with a mock handler.
	mh := mockMatcherHandler{}
	router := SetupRouter(mh, cfg)

	// We’ll run table-driven tests for various routes/methods.
	tests := []struct {
		name           string
		method         string
		path           string
		wantStatusCode int
		wantBodySubstr string
		token          string
		Authorization  string
	}{
		{
			name:           "Health endpoint (GET) - success",
			method:         http.MethodGet,
			path:           "/health",
			wantStatusCode: http.StatusOK,
			wantBodySubstr: `"status": "ok"`,
		},
		{
			name:           "Health endpoint (POST) - method not allowed",
			method:         http.MethodPost,
			path:           "/health",
			wantStatusCode: http.StatusMethodNotAllowed,
			wantBodySubstr: "", // Gorilla by default returns empty body for 405
		},
		{
			name:           "Matcher search endpoint (GET) without token - middleware likely fails",
			method:         http.MethodGet,
			path:           "/matcher/api/v1/search",
			wantStatusCode: http.StatusUnauthorized,
			wantBodySubstr: "Missing Authorization header",
		},
		{
			name:           "Matcher search endpoint (GET) with mock token - expected success",
			method:         http.MethodGet,
			path:           "/matcher/api/v1/search?latitude=40.94289771&longitude=28.0390297&radius=500000",
			token:          "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYXV0aGVudGljYXRlZCI6dHJ1ZSwiaWF0IjoxNTE2MjM5MDIyfQ.KVIuNoXVPw-0qh5d8Mtwx9jmHv97nk0XQaNoBOsatlI",
			wantStatusCode: http.StatusOK,
			wantBodySubstr: `mock MatchDriver called`,
		},
		{
			name:           "Invalid route - returns 404",
			method:         http.MethodGet,
			path:           "/matcher/api/v1/unknown",
			wantStatusCode: http.StatusNotFound,
			wantBodySubstr: "404 page not found", // Gorilla default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)

			// Set Authorization header if token is provided
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.wantStatusCode, rr.Code)
			}

			if tt.wantBodySubstr != "" && !strings.Contains(rr.Body.String(), tt.wantBodySubstr) {
				t.Errorf("Expected response body to contain %q, got: %q",
					tt.wantBodySubstr, rr.Body.String())
			}
		})
	}
}
