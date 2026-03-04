package handlers

import (
	"net/http"

	"karkki-hub/Stock-Portfolio-Manager/internal/services"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	Service *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{Service: s}
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Password string `json:"password"`
	APIKey   string `json:"api_key"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	user, err := h.Service.Register(
		req.Name,
		req.Email,
		req.Phone,
		req.Address,
		req.Password,
		req.APIKey,
	)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "user registered successfully",
		"api_key": user.APIToken,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	token, err := h.Service.Login(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid credentials",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
