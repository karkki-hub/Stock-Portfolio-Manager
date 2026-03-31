package handlers

import (
	"net/http"

	// "karkki-hub/Stock-Portfolio-Manager/internal/models"
	"karkki-hub/Stock-Portfolio-Manager/internal/models"
	"karkki-hub/Stock-Portfolio-Manager/internal/services"

	"github.com/labstack/echo/v4"
)

type TransactionHandler struct {
	Service *services.TransactionService
}

func NewTransactionHandler(s *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{Service: s}
}

type TransactionRequest struct {
	Symbol   string  `json:"symbol"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
}

func (h *TransactionHandler) Buy(c echo.Context) error {

	var req TransactionRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse("invalid request"))
	}

	userId := getUserID(c)

	err := h.Service.Buy(userId, req.Symbol, req.Quantity, req.Price)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, models.SuccessResponse("transaction successful", nil))
}

func (h *TransactionHandler) Sell(c echo.Context) error {

	var req TransactionRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse("invalid request"))
	}

	userId := getUserID(c)

	err := h.Service.Sell(userId, req.Symbol, req.Quantity, req.Price)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, models.SuccessResponse("transaction successful", nil))
}

func (h *TransactionHandler) History(c echo.Context) error {

	userId := getUserID(c)

	transactions, err := h.Service.History(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, models.SuccessResponse("transaction history retrieved", transactions))
}
