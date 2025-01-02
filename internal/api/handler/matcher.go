package handler

import (
	"bitaksi-go-matcher/internal/client"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"bitaksi-go-matcher/internal/models"
)

// MatcherHandler is responsible for handling match requests.
type MatcherHandler struct {
	matcherService MatcherService
}

// MatcherService defines the methods required by MatcherHandler to find the nearest driver.
type MatcherService interface {
	FindNearestDriver(ctx context.Context, latitude, longitude float64, radius int) (*models.DriverWithDistance, error)
}

// NewMatcherHandler creates a new MatcherHandler instance.
func NewMatcherHandler(matcherService MatcherService) MatcherHandler {
	return MatcherHandler{matcherService: matcherService}
}

// MatchDriver handles searching for a nearby driver
// @Summary Search for a Driver
// @Description Find the nearest driver around a GeoJSON point
// @Tags Matcher
// @Param latitude query float64 true "Latitude"
// @Param longitude query float64 true "Longitude"
// @Param radius query int true "Search radius in meters"
// @Success 200 {object} models.DriverWithDistance
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Driver not found"
// @Security BearerAuth
// @Router /matcher/api/v1/search [get]
func (h *MatcherHandler) MatchDriver(w http.ResponseWriter, r *http.Request) {
	latitude, err := strconv.ParseFloat(r.URL.Query().Get("latitude"), 64)
	if err != nil || latitude < -90 || latitude > 90 {
		respondWithError(w, http.StatusBadRequest, "Invalid latitude: must be between -90 and 90")
		return
	}

	longitude, err := strconv.ParseFloat(r.URL.Query().Get("longitude"), 64)
	if err != nil || longitude < -180 || longitude > 180 {
		respondWithError(w, http.StatusBadRequest, "Invalid longitude: must be between -180 and 180")
		return
	}

	radius, err := strconv.Atoi(r.URL.Query().Get("radius"))
	if err != nil || radius <= 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid radius: must be a positive integer")
		return
	}

	driver, err := h.matcherService.FindNearestDriver(r.Context(), latitude, longitude, radius)
	if err != nil {
		log.Printf("FindNearestDriver error: %v\n", err)

		if errors.Is(err, client.ErrDriverNotFound) {
			respondWithError(w, http.StatusNotFound, "Driver not found in the search radius")

			return
		}

		respondWithError(w, http.StatusInternalServerError, "Internal server error")

		return
	}

	respondWithJSON(w, http.StatusOK, driver)
}

// respondWithError is a helper to encode an error message as JSON.
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// respondWithJSON is a helper to encode data as JSON with a specific status code.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("JSON encoding error: %v\n", err)
	}
}
