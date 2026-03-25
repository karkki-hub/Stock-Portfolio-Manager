package handlers

import (
	"net/http"

	// "karkki-hub/Stock-Portfolio-Manager/internal/models"
	"karkki-hub/Stock-Portfolio-Manager/internal/services"
	// "karkki-hub/Stock-Portfolio-Manager/internal/utilities"

	"github.com/labstack/echo/v4"
)

type ProfileHandler struct {
	Service *services.ProfileService
}

func NewProfileHandler(s *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{Service: s}
}

func (h *ProfileHandler) Get(c echo.Context) error {

	userId := getUserID(c)

	transactions, err := h.Service.GetProfile(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, transactions)
}
