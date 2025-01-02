package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"bitaksi-go-matcher/internal/models"
	"github.com/eapache/go-resiliency/breaker"
)

type DriverAPIClient struct {
	BaseURL string
	APIKey  string
	Breaker *breaker.Breaker
	Timeout time.Duration
}

func NewDriverAPI(baseURL, apiKey string) *DriverAPIClient {
	return &DriverAPIClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Breaker: breaker.New(3, 1, 5*time.Second), // Circuit breaker settings
		Timeout: 5 * time.Second,                  // HTTP request timeout
	}
}

func (api *DriverAPIClient) SearchDriver(latitude, longitude float64, radius int) (*models.DriverWithDistance, error) {
	url := fmt.Sprintf("%s/driver/api/v1/search?latitude=%f&longitude=%f&radius=%d", api.BaseURL, latitude, longitude, radius)
	fmt.Println("apiBaseURL: ", api.BaseURL)
	fmt.Println("URL: ", url)
	var driver models.DriverWithDistance

	// Use the circuit breaker to protect the API call
	err := api.Breaker.Run(func() error {
		client := &http.Client{Timeout: api.Timeout}
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		// Add API key in the Authorization header
		req.Header.Set("Authorization", api.APIKey)
		fmt.Println("APIKey: ", api.APIKey)
		fmt.Println(fmt.Sprintf("Request: %+v", req))
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		fmt.Println("Response: ", resp)
		if resp.StatusCode == http.StatusNotFound {
			return errors.New("no driver found")
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("driver service error: %s", resp.Status)
		}

		if err := json.NewDecoder(resp.Body).Decode(&driver); err != nil {
			return fmt.Errorf("failed to decode driver response: %w", err)
		}

		return nil
	})

	// Handle errors from the circuit breaker
	if err != nil {
		if errors.Is(err, breaker.ErrBreakerOpen) {
			return nil, fmt.Errorf("circuit breaker open: failing fast")
		}
		return nil, err
	}

	// No error means the operation succeeded
	// Assuming the last successful call stored the driver response, you may want to return it from storage
	return &driver, nil // Update as necessary if you plan to store the response locally
}
