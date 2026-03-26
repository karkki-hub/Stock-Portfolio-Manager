package routes

import (
	"karkki-hub/Stock-Portfolio-Manager/internal/handlers"
	"karkki-hub/Stock-Portfolio-Manager/internal/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, authHandler *handlers.AuthHandler, jwtSecret string,
) {

	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)

	// Protected routes
	api := e.Group("/api")
	api.Use(middleware.JWTMiddleware(jwtSecret))

	api.GET("/test", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"message": "JWT auth successful"})
	})

}
