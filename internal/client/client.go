package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"bitaksi-go-matcher/internal/models"

	"github.com/eapache/go-resiliency/breaker"
)

// DriverAPIClient handles communication with the Driver Service.
type DriverAPIClient struct {
	BaseURL string
	APIKey  string
	Breaker *breaker.Breaker
	Timeout time.Duration
}

// NewDriverAPI constructs a new DriverAPIClient.
func NewDriverAPI(baseURL, apiKey string) *DriverAPIClient {
	return &DriverAPIClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Breaker: breaker.New(3, 1, 5*time.Second), // e.g., 3 failures allowed before opening
		Timeout: 5 * time.Second,                  // http.Client timeout
	}
}

// SearchDriver queries the Driver Service to find a nearby driver.
func (api *DriverAPIClient) SearchDriver(latitude, longitude float64, radius int) (*models.DriverWithDistance, error) {
	url := fmt.Sprintf("%s/driver/api/v1/search?latitude=%f&longitude=%f&radius=%d",
		api.BaseURL, latitude, longitude, radius)
	log.Printf("[DriverAPIClient] Requesting: %s", url)

	var driver models.DriverWithDistance

	// Wrap our HTTP call in a circuit-breaker to protect from repeated failures
	err := api.Breaker.Run(func() error {
		client := &http.Client{Timeout: api.Timeout}
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		// Add API key in the Authorization header
		req.Header.Set("Authorization", api.APIKey)
		log.Printf("[DriverAPIClient] Authorization: %s", api.APIKey)

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to execute request: %w", err)
		}
		defer resp.Body.Close()

		// Handle status codes
		if resp.StatusCode == http.StatusNotFound {
			return errors.New("no driver found")
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("driver service error: received status %s", resp.Status)
		}

		if err := json.NewDecoder(resp.Body).Decode(&driver); err != nil {
			return fmt.Errorf("failed to decode JSON response: %w", err)
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, breaker.ErrBreakerOpen) {
			return nil, errors.New("circuit breaker open: failing fast")
		}
		return nil, err
	}

	return &driver, nil
}
