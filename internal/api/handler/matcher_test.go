package handler

import (
	"bitaksi-go-matcher/internal/models"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// mockMatcherService is a mock implementation of the MatcherService interface
// for testing purposes.
type mockMatcherService struct {
	FindNearestDriverFunc func(ctx context.Context, lat, long float64, radius int) (*models.DriverWithDistance, error)
}

func (m *mockMatcherService) FindNearestDriver(ctx context.Context, lat, long float64, radius int) (*models.DriverWithDistance, error) {
	if m.FindNearestDriverFunc == nil {
		return nil, errors.New("not implemented")
	}
	return m.FindNearestDriverFunc(ctx, lat, long, radius)
}

func TestMatcherHandler_MatchDriver(t *testing.T) {
	tests := []struct {
		name            string
		queryParams     map[string]string
		mockServiceFunc func(ctx context.Context, lat, long float64, radius int) (*models.DriverWithDistance, error)
		wantStatus      int
		wantContains    string
	}{
		{
			name: "Valid request returns driver",
			queryParams: map[string]string{
				"latitude":  "40.0",
				"longitude": "-70.0",
				"radius":    "100",
			},
			mockServiceFunc: func(ctx context.Context, lat, long float64, radius int) (*models.DriverWithDistance, error) {
				return &models.DriverWithDistance{
					ID:       primitive.NewObjectID(),
					Distance: 10.5,
					Location: models.Location{
						Type:        "Point",
						Coordinates: []float64{-70.0, 40.0},
					},
				}, nil
			},
			wantStatus:   http.StatusOK,
			wantContains: `"distance":10.5`,
		},
		{
			name: "Invalid latitude - out of range",
			queryParams: map[string]string{
				"latitude":  "95.0",
				"longitude": "-70.0",
				"radius":    "100",
			},
			mockServiceFunc: nil, // won't be called
			wantStatus:      http.StatusBadRequest,
			wantContains:    "Invalid latitude",
		},
		{
			name: "Invalid latitude - not a number",
			queryParams: map[string]string{
				"latitude":  "abc",
				"longitude": "-70.0",
				"radius":    "100",
			},
			mockServiceFunc: nil, // won't be called
			wantStatus:      http.StatusBadRequest,
			wantContains:    "Invalid latitude",
		},
		{
			name: "Invalid longitude - out of range",
			queryParams: map[string]string{
				"latitude":  "40.0",
				"longitude": "200",
				"radius":    "100",
			},
			mockServiceFunc: nil, // won't be called
			wantStatus:      http.StatusBadRequest,
			wantContains:    "Invalid longitude",
		},
		{
			name: "Invalid longitude - not a number",
			queryParams: map[string]string{
				"latitude":  "40.0",
				"longitude": "abc",
				"radius":    "100",
			},
			mockServiceFunc: nil,
			wantStatus:      http.StatusBadRequest,
			wantContains:    "Invalid longitude",
		},
		{
			name: "Invalid radius - not a number",
			queryParams: map[string]string{
				"latitude":  "40.0",
				"longitude": "-70.0",
				"radius":    "abc",
			},
			mockServiceFunc: nil,
			wantStatus:      http.StatusBadRequest,
			wantContains:    "Invalid radius",
		},
		{
			name: "Invalid radius - zero or negative",
			queryParams: map[string]string{
				"latitude":  "40.0",
				"longitude": "-70.0",
				"radius":    "0",
			},
			mockServiceFunc: nil,
			wantStatus:      http.StatusBadRequest,
			wantContains:    "Invalid radius",
		},
		{
			name: "Driver not found",
			queryParams: map[string]string{
				"latitude":  "40.0",
				"longitude": "-70.0",
				"radius":    "100",
			},
			mockServiceFunc: func(ctx context.Context, lat, long float64, radius int) (*models.DriverWithDistance, error) {
				return nil, errors.New("driver not found")
			},
			wantStatus:   http.StatusNotFound,
			wantContains: "Driver not found",
		},
		{
			name: "Internal error",
			queryParams: map[string]string{
				"latitude":  "40.0",
				"longitude": "-70.0",
				"radius":    "100",
			},
			mockServiceFunc: func(ctx context.Context, lat, long float64, radius int) (*models.DriverWithDistance, error) {
				return nil, errors.New("some internal error")
			},
			wantStatus:   http.StatusNotFound,
			wantContains: "Driver not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &mockMatcherService{
				FindNearestDriverFunc: tc.mockServiceFunc,
			}
			handler := NewMatcherHandler(mockService)

			// Construct the query string
			values := url.Values{}
			for k, v := range tc.queryParams {
				values.Set(k, v)
			}
			req, err := http.NewRequest(http.MethodGet, "/matcher/api/v1/search?"+values.Encode(), nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler.MatchDriver(rr, req)

			if status := rr.Code; status != tc.wantStatus {
				t.Errorf("Expected status %d, got %d", tc.wantStatus, status)
			}

			// Check if the response body contains the expected substring
			if !strings.Contains(rr.Body.String(), tc.wantContains) {
				t.Errorf("Response body does not contain expected substring.\nExpected: %s\nGot: %s",
					tc.wantContains, rr.Body.String())
			}
		})
	}
}
