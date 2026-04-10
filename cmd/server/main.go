package main

import (
	"karkki-hub/Stock-Portfolio-Manager/internal/config"
	"karkki-hub/Stock-Portfolio-Manager/internal/database"
	"karkki-hub/Stock-Portfolio-Manager/internal/handlers"
	"karkki-hub/Stock-Portfolio-Manager/internal/repository"
	"karkki-hub/Stock-Portfolio-Manager/internal/routes"
	"karkki-hub/Stock-Portfolio-Manager/internal/services"
	"karkki-hub/Stock-Portfolio-Manager/pkg/utilities"

	"github.com/labstack/echo/v4"
)

func main() {

	cfg := config.LoadConfig()

	e := echo.New()

	e.Static("/", "UI")

	db := database.NewMySQL(cfg)

	cronrepo := repository.NewCronRepository(db)

	cronservice := services.NewCronService(cronrepo)

	userRepo := repository.NewUserRepository(db)

	authService := services.NewAuthService(userRepo, cfg.JWTSecret)

	authHandler := handlers.NewAuthHandler(authService)

	stockRepo := repository.NewStockRepository(db)

	stockService := services.NewStockService(stockRepo, cfg.AlphaKey)

	priceService := services.NewPriceService(stockRepo, cronservice)

	cronManager := utilities.NewCronManager()

	// Update prices every 30 minutes

	cronManager.AddJob("*/30 * * * *", func() {
		priceService.UpdatePrices()
	})

	stockHandler := handlers.NewStockHandler(stockService)

	watchRepo := repository.NewWatchlistRepository(db)

	watchService := services.NewWatchlistService(watchRepo, stockRepo, stockService)

	watchHandler := handlers.NewWatchlistHandler(watchService)

	portfolioRepo := repository.NewPortfolioRepository(db)

	portfolioService := services.NewPortfolioService(portfolioRepo)

	portfolioHandler := handlers.NewPortfolioHandler(portfolioService)

	txRepo := repository.NewTransactionRepository(db)

	txService := services.NewTransactionService(txRepo, stockRepo, portfolioService, stockService)

	txHandler := handlers.NewTransactionHandler(txService)

	profileRepo := repository.NewProfileRepository(db)

	profileService := services.NewProfileService(profileRepo)

	profileHandler := handlers.NewProfileHandler(profileService)

	reportRepo := repository.NewReportRepository(db)

	reportService := services.NewReportService(reportRepo)

	reportHandler := handlers.NewReportHandler(reportService, profileService, cronservice)

	// Schedule daily report generation at midnight

	cronManager.AddJob("00 00 * * *", func() {
		reportHandler.DailyReport()
	})

	cronManager.Start()

	routes.RegisterRoutes(e, authHandler, cfg.JWTSecret, stockHandler, watchHandler, txHandler, portfolioHandler, profileHandler, reportHandler)

	e.Logger.Fatal(e.Start(":" + cfg.AppPort))
}
