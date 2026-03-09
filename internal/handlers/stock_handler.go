package handlers

import (
	"net/http"

	"karkki-hub/Stock-Portfolio-Manager/internal/models"
	"karkki-hub/Stock-Portfolio-Manager/internal/services"

	"github.com/labstack/echo/v4"
)

type StockHandler struct {
	Service *services.StockService
}

func NewStockHandler(s *services.StockService) *StockHandler {
	return &StockHandler{Service: s}
}

func (h *StockHandler) SearchStock(c echo.Context) error {
	symbol := c.Param("symbol")

	stock, err := h.Service.SearchStock(symbol)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse("failed to search stock"))
	}
	return c.JSON(http.StatusOK, models.SuccessResponse("stock found", stock))
}
