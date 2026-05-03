// @title           Safa API
// @version         1.0
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	_ "safa-went/docs"
	"safa-went/http/middlewares"
	"safa-went/internal/config"
	"safa-went/internal/providers"
	"safa-went/routes"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Load the configurations
	config := &config.Config{}
	config.Initialize()

	// Initialize database connection
	dbp := &providers.DatabaseProvider{}
	dbp.Connect(config)

	// Initialize Chi router
	r := chi.NewRouter()

	// Middlewares
	r.Use(
		middleware.Logger,
		middleware.Recoverer,
		middlewares.Cors(),
		middlewares.Auth(config.JWTSecret, dbp.DB),
		// Add other middlewares as needed
	)

	// Setup routes
	routes.SetupRoutes(r, dbp.DB)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: r,
	}

	// Create context for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", config.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server gracefully
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
