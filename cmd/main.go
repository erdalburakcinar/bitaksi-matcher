package main

import (
	"context" // Import context package
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bitaksi-go-matcher/internal/api"
	"bitaksi-go-matcher/internal/api/handler"
	"bitaksi-go-matcher/internal/client"
	"bitaksi-go-matcher/internal/config"
	"bitaksi-go-matcher/internal/service"
)

// @title Matcher Service API
// @version 1.0
// @description This is the API documentation for the Matcher Service.
// @termsOfService http://example.com/terms/
// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Driver API client
	driverAPI := client.NewDriverAPI(cfg.DriverServiceClient.URL, cfg.DriverServiceClient.ApiKey)

	// Initialize Matcher Service
	matcherService := service.NewMatcherService(driverAPI)

	// Initialize handlers
	matcherHandler := handler.NewMatcherHandler(matcherService)

	// Setup router
	router := api.SetupRouter(&matcherHandler, cfg)

	// Graceful shutdown setup
	server := &http.Server{
		Addr:         cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start the server in a goroutine
	go func() {
		log.Println("Starting Matcher Service on port:", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for an interrupt signal to gracefully shut down
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down gracefully...")

	// Graceful shutdown with a timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shut down: %v", err)
	}

	log.Println("Server exited cleanly")
}
