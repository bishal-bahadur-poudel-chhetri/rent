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
	// Load the configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to the database
	db, err := config.ConnectDB(cfg.DBConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	vehicleRepo := repositories.NewVehicleRepository(db)
	saleRepo := repositories.NewSaleRepository(db)
	returnRepo := repositories.NewReturnRepository(db)
	videoRepo := repositories.NewVideoRepository(db, cfg)
	futurBookingRepo := repositories.NewFuturBookingRepository(db)
	paymentVerificationRepo := repositories.NewPaymentVerificationRepository(db)
	paymentRepo := repositories.NewPaymentRepository(db)
	saleDetailRepo := repositories.NewSaleDetailRepository(db) // Add this line

	// Initialize services
	returnService := services.NewReturnService(returnRepo)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret, cfg.TokenExpiry)
	vehicleService := services.NewVehicleService(vehicleRepo)
	saleService := services.NewSaleService(saleRepo)
	videoService := services.NewVideoService(videoRepo)
	futurBookingService := services.NewFuturBookingService(futurBookingRepo)
	paymentVerificationService := services.NewPaymentVerificationService(paymentVerificationRepo)
	paymentService := services.NewPaymentService(paymentRepo)
	saleDetailService := services.NewSaleDetailService(saleDetailRepo) // Add this line

	// Initialize handlers
	returnHandler := handlers.NewReturnHandler(returnService, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)
	vehicleHandler := handlers.NewVehicleHandler(vehicleService)
	saleHandler := handlers.NewSaleHandler(saleService, cfg.JWTSecret)
	videoHandler := handlers.NewVideoHandler(videoService)
	futurBookingHandler := handlers.NewFuturBookingHandler(futurBookingService)
	paymentVerificationHandler := handlers.NewPaymentVerification(paymentVerificationService, cfg.JWTSecret)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	saleDetailHandler := handlers.NewSaleDetailHandler(saleDetailService) // Add this line

	// Initialize Gin router
	router := gin.Default()

	// Set trusted proxies if needed (this can be helpful if you're behind a reverse proxy)
	router.SetTrustedProxies([]string{"0.0.0.0"})

	// CORS middleware to allow cross-origin requests
	router.Use(cors.Default())

	// Define API routes under /api/v1
	v1 := router.Group("/api/v1")
	{
		// Public routes
		v1.POST("/register", authHandler.Register)
		v1.POST("/login", authHandler.Login)

		// Protected routes, JWT authentication required
		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			// Vehicle routes
			protected.GET("/vehicle", vehicleHandler.ListVehicles)

			// Return routes
			protected.POST("/sales/:id/return", returnHandler.CreateReturn)

			// Sale routes
			protected.POST("/sales", saleHandler.CreateSale)
			protected.GET("/sales/:id", saleHandler.GetSaleByID)

			// Video upload route
			protected.POST("/sales/upload/video", videoHandler.UploadVideo)
			protected.GET("/payment", paymentHandler.GetPaymentsWithSales)

			// Payment verification route
			protected.PUT("/sales/:payment_id/verify", paymentVerificationHandler.VerifyPayment)

			// FuturBooking route
			protected.GET("/futur-bookings", futurBookingHandler.GetFuturBookingsByMonth)

			// SaleDetail route for filtering sales
			protected.GET("/sales/filter", saleDetailHandler.GetSalesWithFilters) // Add this line
		}
	}

	// Health Check route for monitoring
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
		})
	})

	// Start the server and listen on the configured address
	log.Printf("Server running on %s", cfg.ServerAddress)
	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
