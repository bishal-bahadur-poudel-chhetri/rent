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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load the configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(cfg.DBConnStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	// Get the underlying *sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error getting underlying *sql.DB: %v", err)
	}
	defer sqlDB.Close()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(sqlDB)
	vehicleRepo := repositories.NewVehicleRepository(sqlDB)
	saleRepo := repositories.NewSaleRepository(sqlDB)
	returnRepo := repositories.NewReturnRepository(sqlDB)
	videoRepo := repositories.NewVideoRepository(sqlDB, cfg)
	futurBookingRepo := repositories.NewFuturBookingRepository(sqlDB)
	paymentVerificationRepo := repositories.NewPaymentVerificationRepository(sqlDB)
	paymentRepo := repositories.NewPaymentRepository(sqlDB)
	saleDetailRepo := repositories.NewSaleDetailRepository(sqlDB)
	dataRepo := repositories.NewDataAggregateRepository(sqlDB)
	disableDateRepo := repositories.NewDisableDateRepository(sqlDB)
	statementRepo := repositories.NewStatementRepository(sqlDB)
	expenseRepo := repositories.NewExpenseRepository(sqlDB)
	revenueRepo := repositories.NewRevenueRepository(sqlDB)
	reminderRepo := repositories.NewReminderRepository(db)

	// Initialize services
	returnService := services.NewReturnService(returnRepo)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret, cfg.TokenExpiry)
	vehicleService := services.NewVehicleService(vehicleRepo)
	saleService := services.NewSaleService(saleRepo)
	videoService := services.NewVideoService(videoRepo)
	futurBookingService := services.NewFuturBookingService(futurBookingRepo)
	paymentVerificationService := services.NewPaymentVerificationService(paymentVerificationRepo)
	paymentService := services.NewPaymentService(paymentRepo)
	saleDetailService := services.NewSaleDetailService(saleDetailRepo)
	dataService := services.NewDataAggregateService(dataRepo, dataRepo)
	disableDateService := services.NewDisableDateService(disableDateRepo)
	statementService := services.NewStatementService(statementRepo)
	expenseService := services.NewExpenseService(expenseRepo)
	revenueService := services.NewRevenueService(revenueRepo)
	reminderService := services.NewReminderService(reminderRepo)

	// Initialize handlers
	returnHandler := handlers.NewReturnHandler(returnService, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)
	vehicleHandler := handlers.NewVehicleHandler(vehicleService)
	saleHandler := handlers.NewSaleHandler(saleService, cfg.JWTSecret)
	videoHandler := handlers.NewVideoHandler(videoService, "https://pub-8da91f66939f4cdc9e4206024a0e68e9.r2.dev")
	futurBookingHandler := handlers.NewFuturBookingHandler(futurBookingService)
	paymentVerificationHandler := handlers.NewPaymentVerification(paymentVerificationService, cfg.JWTSecret)
	paymentHandler := handlers.NewPaymentHandler(paymentService, cfg.JWTSecret)
	saleDetailHandler := handlers.NewSaleDetailHandler(saleDetailService)
	dataHandler := handlers.NewDataAggregateHandler(dataService)
	disableDateHandler := handlers.NewDisableDateHandler(disableDateService)
	statementHandler := handlers.NewStatementHandler(statementService)
	expenseHandler := handlers.NewExpenseHandler(expenseService)
	revenueHandler := handlers.NewRevenueHandler(revenueService)
	reminderHandler := handlers.NewReminderHandler(reminderService)

	// Initialize Gin router
	router := gin.Default()

	// Set trusted proxies if needed
	router.SetTrustedProxies([]string{"0.0.0.0"})

	// CORS middleware to allow cross-origin requests
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

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
			protected.GET("/sales", saleHandler.GetSales)
			protected.PUT("/sales/:saleID", saleHandler.UpdateSaleByUserID)

			// Video upload route
			protected.POST("/sales/upload/video", videoHandler.UploadVideo)
			protected.GET("/payment", paymentHandler.GetPaymentsWithSales)

			// Payment verification routes
			protected.PUT("/payment/:payment_id/sales/:sale_id/verify", paymentVerificationHandler.VerifyPayment) // Updated path
			protected.GET("/payment/:payment_id", paymentVerificationHandler.GetPaymentDetails)                   // New GET endpoint
			protected.POST("/payment/:payment_id/cancel", paymentVerificationHandler.CancelPayment)               // New Cancel endpoint

			// Payment routes
			protected.PUT("/payment/:payment_id", paymentHandler.UpdatePayment)
			protected.POST("/sales/:id/payment", paymentHandler.InsertPayment)

			// FuturBooking routes
			protected.GET("/futur-bookings", futurBookingHandler.GetFuturBookingsByMonth)
			protected.POST("/futur-bookings/:saleID/cancel", futurBookingHandler.CancelFuturBooking)

			// SaleDetail route for filtering sales
			protected.GET("/sales/filter", saleDetailHandler.GetSalesWithFilters)

			// DataAggregate route
			protected.GET("/aggregate", dataHandler.GetAggregatedData)
			protected.GET("/disabled-dates", disableDateHandler.GetDisabledDates)

			// Statement routes
			protected.GET("/statements", statementHandler.GetOutstandingStatements)

			// Expense routes
			protected.POST("/expenses", expenseHandler.CreateExpense)
			protected.PUT("/expenses/:id", expenseHandler.UpdateExpense)
			protected.DELETE("/expenses/:id", expenseHandler.DeleteExpense)
			protected.GET("/expenses/:id", expenseHandler.GetExpenseByID)
			protected.GET("/expenses", expenseHandler.GetAllExpenses)

			// Revenue route
			protected.GET("/revenue", revenueHandler.GetRevenue)
			protected.GET("/revenue/monthly", revenueHandler.GetMonthlyRevenue)
			protected.GET("/revenue/mobile-visualization", revenueHandler.GetMobileRevenueVisualization)

			// Reminder routes
			protected.POST("/reminders", reminderHandler.CreateReminder)
			protected.POST("/reminders/:reminder_id/acknowledge", reminderHandler.AcknowledgeReminder)
			protected.GET("/vehicles/:vehicle_id/reminders", reminderHandler.GetRemindersByVehicle)
			protected.GET("/reminders/due", reminderHandler.GetDueReminders)
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
