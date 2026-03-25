package main

import (
	"karkki-hub/Stock-Portfolio-Manager/internal/config"
	"karkki-hub/Stock-Portfolio-Manager/internal/database"
	"karkki-hub/Stock-Portfolio-Manager/internal/handlers"
	"karkki-hub/Stock-Portfolio-Manager/internal/repository"
	"karkki-hub/Stock-Portfolio-Manager/internal/routes"
	"karkki-hub/Stock-Portfolio-Manager/internal/services"

	"github.com/labstack/echo/v4"
)

func main() {

	cfg := config.LoadConfig()

	e := echo.New()

	e.Static("/", "UI")

	db := database.NewMySQL(cfg)

	userRepo := repository.NewUserRepository(db)

	authService := services.NewAuthService(userRepo, cfg.JWTSecret)

	authHandler := handlers.NewAuthHandler(authService)

	stockRepo := repository.NewStockRepository(db)

	stockService := services.NewStockService(stockRepo, cfg.AlphaKey)

	stockHandler := handlers.NewStockHandler(stockService)

	watchRepo := repository.NewWatchlistRepository(db)

	watchService := services.NewWatchlistService(watchRepo, stockRepo)

	watchHandler := handlers.NewWatchlistHandler(watchService)

	portfolioRepo := repository.NewPortfolioRepository(db)

	portfolioService := services.NewPortfolioService(portfolioRepo)

	portfolioHandler := handlers.NewPortfolioHandler(portfolioService)

	txRepo := repository.NewTransactionRepository(db)

	txService := services.NewTransactionService(txRepo, stockRepo, portfolioService)

	txHandler := handlers.NewTransactionHandler(txService)

	profileRepo := repository.NewProfileRepository(db)

	profileService := services.NewProfileService(profileRepo)

	profileHandler := handlers.NewProfileHandler(profileService)

	routes.RegisterRoutes(e, authHandler, cfg.JWTSecret, stockHandler, watchHandler, txHandler, portfolioHandler, profileHandler)

	e.Logger.Fatal(e.Start(":" + cfg.AppPort))
}
