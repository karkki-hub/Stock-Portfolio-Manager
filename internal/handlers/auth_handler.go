package handlers

import (
	"net/http"
	"regexp"

	"karkki-hub/Stock-Portfolio-Manager/internal/models"
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
		return c.JSON(http.StatusBadRequest, models.ErrorResponse("invalid request"))
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse("email, password, and name are required"))
	}

	if !isValidEmail(req.Email) {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse("invalid email format"))
	}

	if len(req.Phone) != 10 {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse("phone number must be 10 digits"))
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
		return c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusCreated, models.SuccessResponse(
		"user registered successfully",
		map[string]string{
			"api_key": user.APIToken,
		},
	))

}

func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse("invalid request"))
	}

	token, err := h.Service.Login(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse("invalid credentials"))
	}

	return c.JSON(http.StatusOK, models.SuccessResponse(
		"login successful",
		map[string]string{
			"token": token,
		},
	))
}

func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
