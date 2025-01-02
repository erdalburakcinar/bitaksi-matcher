package service

import (
	"context"
	"fmt"

	"bitaksi-go-matcher/internal/models"
)

// MatcherService defines the interface for matching operations
type MatcherService interface {
	FindNearestDriver(ctx context.Context, latitude, longitude float64, radius int) (*models.DriverWithDistance, error)
}

type matcherService struct {
	driverClient DriverClient
}

// DriverClient defines the interface for interacting with the Driver Service
type DriverClient interface {
	SearchDriver(latitude, longitude float64, radius int) (*models.DriverWithDistance, error)
}

// NewMatcherService creates a new instance of MatcherService
func NewMatcherService(driverClient DriverClient) MatcherService {
	return &matcherService{driverClient: driverClient}
}

// FindNearestDriver calls the Driver Service to find the nearest driver
func (s *matcherService) FindNearestDriver(ctx context.Context, latitude, longitude float64, radius int) (*models.DriverWithDistance, error) {
	driver, err := s.driverClient.SearchDriver(latitude, longitude, radius)
	if err != nil {
		return nil, fmt.Errorf("failed to find nearest driver: %w", err)
	}

	return driver, nil
}
