package main

import (
	"log"
	"renting/internal/config"
	"renting/internal/handlers"
	"renting/internal/middleware"
	"renting/internal/repositories"
	"renting/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := config.ConnectDB(cfg.DBConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret, cfg.TokenExpiry)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	vehicleRepo := repositories.NewVehicleRepository(db)
	vehicleService := services.NewVehicleService(vehicleRepo)
	vehicleHandler := handlers.NewVehicleHandler(vehicleService)

	// Set up router
	router := gin.Default()
	router.Use(cors.Default())

	v1 := router.Group("/api/v1")
	{
		v1.POST("/register", authHandler.Register) // Add the register route
		v1.POST("/login", authHandler.Login)
		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(cfg.JWTSecret))
		{

			protected.POST("/vehical", vehicleHandler.RegisterVehicleHandler)
			protected.GET("/vehical", vehicleHandler.ListVehicles)

		}
	}

	// Start server
	log.Printf("Server running on %s", cfg.ServerAddress)
	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
