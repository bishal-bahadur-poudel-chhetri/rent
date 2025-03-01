package main

import (
	"log"
	"renting/internal/config"
	"renting/internal/handlers"
	"renting/internal/middleware"
	"renting/internal/repositories"
	"renting/internal/services"

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

	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret, cfg.TokenExpiry)
	authHandler := handlers.NewAuthHandler(authService)

	vehicleRepo := repositories.NewVehicleRepository(db)
	vehicleService := services.NewVehicleService(vehicleRepo)
	vehicleHandler := handlers.NewVehicleHandler(vehicleService)

	router := gin.Default()
	router.Use(middleware.CORS())
	v1 := router.Group("/api/v1")
	{
		v1.POST("/register", authHandler.Register)
		v1.POST("/login", authHandler.Login)
		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			protected.GET("/profile", authHandler.GetProfile)
			protected.PUT("/profile", authHandler.UpdateProfile)

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
