package service

import (
	"context"
	"fmt"
	"log"

	"bitaksi-go-matcher/internal/models"
)

// MatcherService defines the interface for matching operations.
type MatcherService interface {
	FindNearestDriver(ctx context.Context, latitude, longitude float64, radius int) (*models.DriverWithDistance, error)
}

type matcherService struct {
	driverClient DriverClient
}

// DriverClient defines the interface for interacting with the Driver Service.
type DriverClient interface {
	SearchDriver(latitude, longitude float64, radius int) (*models.DriverWithDistance, error)
}

// NewMatcherService creates a new instance of the MatcherService.
// driverClient should be an implementation of DriverClient that communicates
func NewMatcherService(driverClient DriverClient) MatcherService {
	return &matcherService{
		driverClient: driverClient,
	}
}

// FindNearestDriver uses the DriverClient to retrieve a driver closest to the given coordinates.
func (s *matcherService) FindNearestDriver(ctx context.Context, latitude, longitude float64, radius int) (*models.DriverWithDistance, error) {
	log.Printf("[MatcherService] Finding driver near lat=%.6f, lon=%.6f, radius=%d", latitude, longitude, radius)

	driver, err := s.driverClient.SearchDriver(latitude, longitude, radius)
	if err != nil {
		return nil, fmt.Errorf("failed to find nearest driver: %w", err)
	}
	return driver, nil
}
