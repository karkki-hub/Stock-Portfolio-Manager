package handlers

import (
	// "encoding/json"
	"net/http"

	// "karkki-hub/Stock-Portfolio-Manager/internal/models"
	"karkki-hub/Stock-Portfolio-Manager/internal/models"
	"karkki-hub/Stock-Portfolio-Manager/internal/services"
	"karkki-hub/Stock-Portfolio-Manager/internal/utilities"

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

func (h *ProfileHandler) Update(c echo.Context) error {
	type UpdateUser struct {
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		Address string `json:"address"`
	}

	var req UpdateUser

	// Bind request body to struct
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if req.Email == "" || req.Address == "" || req.Phone == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "all fields are required",
		})
	}

	if !utilities.IsValidEmail(req.Email) {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse("invalid email format"))
	}

	userID := getUserID(c)

	err := h.Service.Repo.Update(userID, req.Phone, req.Email, req.Address)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Profile updated successfully",
	})
}

func (h *ProfileHandler) Reset(c echo.Context) error {
	type ResetPassword struct {
		OldPassword     string `json:"oldpassword"`
		NewPassword     string `json:"newpassword"`
		ReEnterPassword string `json:"reenterpassword"`
	}

	var req ResetPassword

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if req.NewPassword != req.ReEnterPassword {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Re enered Password doesn't match",
		})
	}

	userID := getUserID(c)

	storedHash, err := h.Service.Repo.ExistingPassword(userID)
	if err != nil {
		return err
	}

	if err := utilities.CheckPasswordHash(req.OldPassword, storedHash); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Old password is incorrect",
		})
	}

	if err := utilities.CheckPasswordHash(req.NewPassword, storedHash); err == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Reusing the old password not allowed",
		})
	}

	er := h.Service.ChangePassword(userID, req.NewPassword)

	if er != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": er.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Password changed successfully",
	})
}
