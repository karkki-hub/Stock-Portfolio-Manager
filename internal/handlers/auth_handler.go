package handlers

import (
	"net/http"

	"karkki-hub/Stock-Portfolio-Manager/internal/services"

	"github.com/labstack/echo/v4"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c echo.Context) error {
	req := new(AuthRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	err := services.Register(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "User registered")
}

func Login(c echo.Context) error {
	req := new(AuthRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	token, err := services.Login(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
