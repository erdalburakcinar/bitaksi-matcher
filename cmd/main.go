package main

import (
	"log"
	"net/http"

	"bitaksi-go-matcher/internal/api"
	"bitaksi-go-matcher/internal/api/handler"
	"bitaksi-go-matcher/internal/client"
	"bitaksi-go-matcher/internal/config"
	"bitaksi-go-matcher/internal/middleware"
	"bitaksi-go-matcher/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize JWT Middleware
	jwtMiddleware := middleware.NewJWTMiddleware(cfg.JWTSecretKey)

	// Initialize Driver API client
	driverAPI := client.NewDriverAPI(cfg.DriverServiceClient.URL, cfg.DriverServiceClient.ApiKey)

	// Initialize Matcher Service
	matcherService := service.NewMatcherService(driverAPI)

	// Initialize handlers
	matcherHandler := handler.NewMatcherHandler(matcherService)

	// Setup router
	router := api.SetupRouter(matcherHandler, jwtMiddleware)

	// Start the server
	log.Println("Starting Matcher Service on port:", cfg.Server.Port)
	if err := http.ListenAndServe(cfg.Server.Port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
