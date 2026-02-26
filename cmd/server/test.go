package main

import (
	"karkki-hub/Stock-Portfolio-Manager/internal/routes"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	godotenv.Load()

	e := echo.New()
	routes.SetupRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
