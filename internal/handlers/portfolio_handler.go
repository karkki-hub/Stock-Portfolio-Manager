package handlers

import (
	"net/http"

	"karkki-hub/Stock-Portfolio-Manager/internal/models"
	"karkki-hub/Stock-Portfolio-Manager/internal/services"

	"github.com/labstack/echo/v4"
)

type PortfolioHandler struct {
	Service *services.PortfolioService
}

func NewPortfolioHandler(s *services.PortfolioService) *PortfolioHandler {
	return &PortfolioHandler{Service: s}
}

func (h *PortfolioHandler) Get(c echo.Context) error {

	userId := getUserID(c)

	data, err := h.Service.Get(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, models.SuccessResponse("portfolio data retrieved", data))
}
