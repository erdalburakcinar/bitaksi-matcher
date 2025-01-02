package api

import (
	"net/http"

	"bitaksi-go-matcher/internal/config"
	"bitaksi-go-matcher/internal/middleware"

	"github.com/gorilla/mux"
)

type MatcherHandler interface {
	MatchDriver(w http.ResponseWriter, r *http.Request)
}

// SetupRouter initializes and configures all HTTP routes for the Matcher Service.
func SetupRouter(matcherHandler MatcherHandler, cfg *config.Config) *mux.Router {
	router := mux.NewRouter()

	// Public routes (e.g., health check)
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	}).Methods(http.MethodGet)

	// Protected routes (require JWT authentication)
	matcherRouter := router.PathPrefix("/matcher/api/v1").Subrouter()
	matcherRouter.Use(middleware.ValidateJWT(cfg.JWTSecretKey))

	// Register matcher endpoints
	matcherRouter.HandleFunc("/search", matcherHandler.MatchDriver).Methods(http.MethodGet)

	return router
}
