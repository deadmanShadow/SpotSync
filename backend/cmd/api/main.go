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

	"spotsync-backend/internal/config"
	"spotsync-backend/internal/database"
	"spotsync-backend/internal/handler"
	"spotsync-backend/internal/middleware"
	"spotsync-backend/internal/repository"
	"spotsync-backend/internal/service"
	"spotsync-backend/pkg/utils"
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

	// 3. Health check route (public)
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "UP"})
	})

	// 4. Connect to the database. We do this BEFORE the server starts so the
	// route registration block below can be entered when the connection is
	// healthy. Auto-migration is run in a background goroutine because it can
	// be slow on remote databases (e.g. Neon) and we don't want it to delay
	// the HTTP server start.
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Printf("WARNING: Database connection failed: %v", err)
		log.Println("Server will continue running but DB-dependent endpoints will not work.")
	} else {
		// Repositories
		userRepo := repository.NewUserRepository(db)
		zoneRepo := repository.NewZoneRepository(db)
		reservationRepo := repository.NewReservationRepository(db)

		// Services
		authService := service.NewAuthService(userRepo, cfg.JWTSecret)
		zoneService := service.NewZoneService(zoneRepo)
		reservationService := service.NewReservationService(reservationRepo, zoneRepo)

		// Handlers
		authHandler := handler.NewAuthHandler(authService)
		zoneHandler := handler.NewZoneHandler(zoneService)
		reservationHandler := handler.NewReservationHandler(reservationService)

		// Run GORM auto-migrations in the background so the server can
		// accept traffic immediately.
		go func() {
			if err := database.AutoMigrate(db); err != nil {
				log.Printf("WARNING: Auto-migration failed: %v", err)
			} else {
				log.Println("Database migrations applied successfully")
			}
		}()

		// --- Public auth routes ---
		authGroup := e.Group("/api/v1/auth")
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)

		// --- Authenticated routes ---
		auth := middleware.AuthMiddleware(cfg.JWTSecret)

		// Zones: any authenticated user can read; admins mutate.
		zoneGroup := e.Group("/api/v1/zones", auth)
		zoneGroup.GET("", zoneHandler.List)
		zoneGroup.GET("/:id", zoneHandler.GetByID)
		zoneGroup.POST("", zoneHandler.Create, middleware.AdminOnly())
		zoneGroup.PUT("/:id", zoneHandler.Update, middleware.AdminOnly())
		zoneGroup.DELETE("/:id", zoneHandler.Delete, middleware.AdminOnly())

		// Reservations: any authenticated user can manage their own; admins
		// can see/modify all.
		reservationGroup := e.Group("/api/v1/reservations", auth)
		reservationGroup.GET("/mine", reservationHandler.ListMine)
		reservationGroup.GET("", reservationHandler.ListAll, middleware.AdminOnly())
		reservationGroup.GET("/:id", reservationHandler.GetByID)
		reservationGroup.POST("", reservationHandler.Create)
		reservationGroup.PATCH("/:id/status", reservationHandler.UpdateStatus)
		reservationGroup.DELETE("/:id", reservationHandler.Delete)
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