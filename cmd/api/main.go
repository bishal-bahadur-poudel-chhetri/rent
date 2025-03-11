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
	vehicleRepo := repositories.NewVehicleRepository(db)
	saleRepo := repositories.NewSaleRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret, cfg.TokenExpiry)
	vehicleService := services.NewVehicleService(vehicleRepo)
	saleService := services.NewSaleService(saleRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	vehicleHandler := handlers.NewVehicleHandler(vehicleService)
	saleHandler := handlers.NewSaleHandler(saleService, cfg.JWTSecret)

	// Set up router
	router := gin.Default()

	// Trust all proxies (if your app is behind a reverse proxy)
	router.SetTrustedProxies([]string{"0.0.0.0"}) // Replace with your proxy IP(s)
	router.Use(cors.Default())

	v1 := router.Group("/api/v1")
	{
		// Public routes
		v1.POST("/register", authHandler.Register)
		v1.POST("/login", authHandler.Login)

		// Protected routes (require JWT authentication)
		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			// Vehicle routes
			// protected.POST("/vehicle", vehicleHandler.RegisterVehicleHandler)
			protected.GET("/vehicle", vehicleHandler.ListVehicles) // Updated to handle future booking details

			// Sale routes
			protected.POST("/sales", saleHandler.CreateSale)
			protected.GET("/sales/:id", saleHandler.GetSaleByID)
		}
	}

	// Start server
	log.Printf("Server running on %s", cfg.ServerAddress)
	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
