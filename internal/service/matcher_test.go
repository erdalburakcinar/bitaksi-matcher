package service

import (
	"bitaksi-go-matcher/internal/models"
	"context"
	"errors"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// mockDriverClient is a mock implementation of DriverClient.
type mockDriverClient struct {
	SearchDriverFunc func(latitude, longitude float64, radius int) (*models.DriverWithDistance, error)
}

func (m *mockDriverClient) SearchDriver(latitude, longitude float64, radius int) (*models.DriverWithDistance, error) {
	if m.SearchDriverFunc == nil {
		return nil, errors.New("not implemented")
	}
	return m.SearchDriverFunc(latitude, longitude, radius)
}

func TestMatcherService_FindNearestDriver(t *testing.T) {
	tests := []struct {
		name          string
		mockFunc      func(lat, lon float64, r int) (*models.DriverWithDistance, error)
		latitude      float64
		longitude     float64
		radius        int
		wantErr       bool
		wantErrSubstr string
		wantDriver    *models.DriverWithDistance
	}{
		{
			name: "Success - found driver",
			mockFunc: func(lat, lon float64, r int) (*models.DriverWithDistance, error) {
				return &models.DriverWithDistance{
					ID:       primitive.NewObjectID(),
					Distance: 120.5,
				}, nil
			},
			latitude:  40.0,
			longitude: 29.0,
			radius:    1000,
			wantErr:   false,
			wantDriver: &models.DriverWithDistance{
				Distance: 120.5,
			},
		},
		{
			name: "Driver not found - service error",
			mockFunc: func(lat, lon float64, r int) (*models.DriverWithDistance, error) {
				return nil, errors.New("no driver available")
			},
			latitude:      40.0,
			longitude:     29.0,
			radius:        1000,
			wantErr:       true,
			wantErrSubstr: "no driver available",
			wantDriver:    nil,
		},
		{
			name: "Generic error from driver client",
			mockFunc: func(lat, lon float64, r int) (*models.DriverWithDistance, error) {
				return nil, errors.New("service unreachable")
			},
			latitude:      41.5,
			longitude:     28.0,
			radius:        200,
			wantErr:       true,
			wantErrSubstr: "service unreachable",
			wantDriver:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockDriverClient{
				SearchDriverFunc: tt.mockFunc,
			}
			service := NewMatcherService(mockClient)

			driver, err := service.FindNearestDriver(context.Background(), tt.latitude, tt.longitude, tt.radius)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got none")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("did not expect error, got: %v", err)
			}
			if tt.wantErr && err != nil && tt.wantErrSubstr != "" {
				if !containsSubstring(err.Error(), tt.wantErrSubstr) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.wantErrSubstr)
				}
			}

			// If we expect a driver, compare the important fields
			if !tt.wantErr && tt.wantDriver != nil {
				if driver.Distance != tt.wantDriver.Distance {
					t.Errorf("expected Distance=%.1f, got %.1f",
						tt.wantDriver.Distance, driver.Distance)
				}
			}
		})
	}
}

func containsSubstring(s, substr string) bool {
	return len(substr) == 0 || (len(substr) > 0 && len(s) >= len(substr) &&
		contains(s, substr))
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr ||
		(len(s) > len(substr) && (s[0:len(substr)] == substr || contains(s[1:], substr))))
}
