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

	routes.RegisterRoutes(e, authHandler, cfg.JWTSecret, stockHandler, watchHandler)

	e.Logger.Fatal(e.Start(":" + cfg.AppPort))
}
