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

	routes.RegisterRoutes(e, authHandler, cfg.JWTSecret)

	e.Logger.Fatal(e.Start(":" + cfg.AppPort))
}
