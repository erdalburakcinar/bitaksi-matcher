package handler

import (
	"bitaksi-go-matcher/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type MatcherHandler struct {
	matcherService MatcherService
}

type MatcherService interface {
	FindNearestDriver(ctx context.Context, latitude, longitude float64, radius int) (*models.DriverWithDistance, error)
}

// NewMatcherHandler creates a new MatcherHandler instance
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
// @Router /matcher/api/v1/search [get]
func (h *MatcherHandler) MatchDriver(w http.ResponseWriter, r *http.Request) {
	latitude, err := strconv.ParseFloat(r.URL.Query().Get("latitude"), 64)
	if err != nil || latitude < -90 || latitude > 90 {
		http.Error(w, `{"error": "Invalid latitude: must be between -90 and 90"}`, http.StatusBadRequest)
		return
	}

	longitude, err := strconv.ParseFloat(r.URL.Query().Get("longitude"), 64)
	if err != nil || longitude < -180 || longitude > 180 {
		http.Error(w, `{"error": "Invalid longitude: must be between -180 and 180"}`, http.StatusBadRequest)
		return
	}

	radius, err := strconv.Atoi(r.URL.Query().Get("radius"))
	if err != nil || radius <= 0 {
		http.Error(w, `{"error": "Invalid radius: must be a positive integer"}`, http.StatusBadRequest)
		return
	}

	driver, err := h.matcherService.FindNearestDriver(r.Context(), latitude, longitude, radius)
	if err != nil {
		fmt.Println("Error: ", err)
		http.Error(w, `{"error": "Driver not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(driver)
}
