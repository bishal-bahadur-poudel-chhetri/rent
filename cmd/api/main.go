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

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := config.ConnectDB(cfg.DBConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := repositories.NewUserRepository(db)
	vehicleRepo := repositories.NewVehicleRepository(db)
	saleRepo := repositories.NewSaleRepository(db)
	returnRepo := repositories.NewReturnRepository(db)
	returnService := services.NewReturnService(returnRepo)
	returnHandler := handlers.NewReturnHandler(returnService, cfg.JWTSecret)

	authService := services.NewAuthService(userRepo, cfg.JWTSecret, cfg.TokenExpiry)
	vehicleService := services.NewVehicleService(vehicleRepo)
	saleService := services.NewSaleService(saleRepo)

	authHandler := handlers.NewAuthHandler(authService)
	vehicleHandler := handlers.NewVehicleHandler(vehicleService)
	saleHandler := handlers.NewSaleHandler(saleService, cfg.JWTSecret)

	videoRepo := repositories.NewVideoRepository(db, cfg)
	videoService := services.NewVideoService(videoRepo)
	videoHandler := handlers.NewVideoHandler(videoService)

	router := gin.Default()

	router.SetTrustedProxies([]string{"0.0.0.0"})
	router.Use(cors.Default())

	v1 := router.Group("/api/v1")
	{

		v1.POST("/register", authHandler.Register)
		v1.POST("/login", authHandler.Login)

		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(cfg.JWTSecret))
		{

			protected.GET("/vehicle", vehicleHandler.ListVehicles)
			protected.POST("/sales/:id/return", returnHandler.CreateReturn)

			protected.POST("/sales", saleHandler.CreateSale)
			protected.GET("/sales/:id", saleHandler.GetSaleByID)

			protected.POST("/sales/upload/video", videoHandler.UploadVideo)
		}
	}

	log.Printf("Server running on %s", cfg.ServerAddress)
	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
