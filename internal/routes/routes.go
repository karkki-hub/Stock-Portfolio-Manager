package routes

import (
	"karkki-hub/Stock-Portfolio-Manager/internal/handlers"
	"karkki-hub/Stock-Portfolio-Manager/internal/middleware"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {

	e.POST("/register", handlers.Register)
	e.POST("/login", handlers.Login)

	protected := e.Group("/api")
	protected.Use(middleware.JWTMiddleware())
	protected.GET("/test", func(c echo.Context) error {
		return c.JSON(200, "JWT Protected Route Working")
	})
}
