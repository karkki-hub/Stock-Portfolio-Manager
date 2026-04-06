package handlers

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"karkki-hub/Stock-Portfolio-Manager/internal/models"
	"karkki-hub/Stock-Portfolio-Manager/internal/services"
)

type WatchlistHandler struct {
	Service *services.WatchlistService
}

func NewWatchlistHandler(s *services.WatchlistService) *WatchlistHandler {
	return &WatchlistHandler{Service: s}
}

func getUserID(c echo.Context) uint {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return uint(claims["user_id"].(float64))
}

func (h *WatchlistHandler) Add(c echo.Context) error {

	symbol := c.QueryParam("symbol")

	userID := getUserID(c)
	err := h.Service.Add(userID, symbol)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, models.SuccessResponse("stock added to watchlist", nil))
}

func (h *WatchlistHandler) Remove(c echo.Context) error {

	symbol := c.Param("symbol")
	userID := getUserID(c)
	err := h.Service.Remove(userID, symbol)

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, models.SuccessResponse("stock removed from watchlist", nil))
}

func (h *WatchlistHandler) Get(c echo.Context) error {

	userID := getUserID(c)
	fmt.Print(userID)
	stocks, err := h.Service.Get(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, models.SuccessResponse("watchlist retrieved", stocks))
}

func (h *WatchlistHandler) GetStockHistory(c echo.Context) error {

	symbol := c.Param("symbol")
	stock, err := h.Service.StockRepo.GetBySymbol(symbol)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error()))
	}
	history, err := h.Service.GetStockHistory(stock.Symbol)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, models.SuccessResponse("stock history retrieved", history))
}
