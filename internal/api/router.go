package api

import (
	"bitaksi-go-matcher/internal/api/handler"
	"bitaksi-go-matcher/internal/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func SetupRouter(matcherHandler handler.MatcherHandler, jwtMiddleware *middleware.JWTMiddleware) *mux.Router {
	router := mux.NewRouter()

	// Public routes (e.g., health check)
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	// Protected routes
	matcherRouter := router.PathPrefix("/matcher/api/v1").Subrouter()
	matcherRouter.Use(jwtMiddleware.ValidateJWT)

	// Register endpoints
	matcherRouter.HandleFunc("/search", matcherHandler.MatchDriver).Methods("GET")

	return router
}
