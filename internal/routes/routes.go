package routes

import (
	"karkki-hub/Stock-Portfolio-Manager/internal/handlers"
	"karkki-hub/Stock-Portfolio-Manager/internal/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo,
	authHandler *handlers.AuthHandler,
	jwtSecret string,
	stockHandler *handlers.StockHandler,
	watchHandler *handlers.WatchlistHandler,
	txHandler *handlers.TransactionHandler,
	portfolioHandler *handlers.PortfolioHandler,
	profileHandler *handlers.ProfileHandler,
	reportHandler *handlers.ReportHandler,
) {

	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)

	// Protected routes
	api := e.Group("/api")
	api.Use(middleware.JWTMiddleware(jwtSecret))

	api.GET("/test", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"message": "JWT auth successful"})
	})

	api.GET("/stocks/:keyword", stockHandler.SearchStock)

	api.GET("/watchlist", watchHandler.Get)
	api.POST("/watchlist", watchHandler.Add)
	api.DELETE("/watchlist/:symbol", watchHandler.Remove)

	api.POST("/transactions/buy", txHandler.Buy)
	api.POST("/transactions/sell", txHandler.Sell)
	api.GET("/transactions", txHandler.History)
	api.GET("/portfolio", portfolioHandler.Get)

	api.GET("/profile", profileHandler.Get)
	api.PUT("/profile/update", profileHandler.Update)
	api.POST("/profile/reset_pswd", profileHandler.Reset)

	api.GET("/report", reportHandler.ExportReportCSV)
}
