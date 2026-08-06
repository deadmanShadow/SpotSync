package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"

	"spotsync-backend/config"
	"spotsync-backend/database"
	"spotsync-backend/utils"
)

func main() {
	// 1. Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Initialize Echo
	e := echo.New()
	e.HideBanner = true
	e.Validator = utils.NewValidator()

	// 3. Health check route (public, used for Step 1 verification)
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "UP"})
	})

	// 4. Connect to the database
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Printf("WARNING: Database connection failed: %v", err)
		log.Println("Server will continue running but DB-dependent endpoints will not work.")
	} else {
		// Auto-migration will be wired in Step 2
		_ = db
		log.Println("Database is ready")
	}

	// 5. Start server with graceful shutdown
	go func() {
		if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	log.Printf("SpotSync API server listening on port %s", cfg.Port)

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
	log.Println("Server stopped")
}
