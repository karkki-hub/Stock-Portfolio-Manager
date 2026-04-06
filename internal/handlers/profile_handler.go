package handlers

import (
	"net/http"

	"karkki-hub/Stock-Portfolio-Manager/internal/models"
	"karkki-hub/Stock-Portfolio-Manager/internal/services"
	"karkki-hub/Stock-Portfolio-Manager/pkg/utilities"

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
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, models.SuccessResponse("profile data retrieved", transactions))
}

func (h *ProfileHandler) Update(c echo.Context) error {
	type UpdateUser struct {
		Name    string `json:"name"`
		Phone   string `json:"phone"`
		Address string `json:"address"`
	}

	var req UpdateUser

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ErrorResponse("Invalid request body"))
	}

	if req.Name == "" || req.Address == "" || req.Phone == "" {
		return c.JSON(http.StatusBadRequest,
			models.ErrorResponse("all fields are required"))
	}

	userID := getUserID(c)

	err := h.Service.ChangeProfile(userID, req.Phone, req.Name, req.Address)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			models.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, models.SuccessResponse("profile updated successfully", nil))
}

func (h *ProfileHandler) Reset(c echo.Context) error {
	type ResetPassword struct {
		OldPassword     string `json:"oldpassword"`
		NewPassword     string `json:"newpassword"`
		ReEnterPassword string `json:"reenterpassword"`
	}

	var req ResetPassword

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ErrorResponse("Invalid request body"))
	}

	if !utilities.IsValidPassword(req.NewPassword) {
		return c.JSON(http.StatusBadRequest,
			models.ErrorResponse("Password must be at least 8 characters long and contain at least one uppercase and one lowercase letter and a special character"))
	}

	if req.NewPassword != req.ReEnterPassword {
		return c.JSON(http.StatusBadRequest,
			models.ErrorResponse("Re-entered password doesn't match"))
	}

	userID := getUserID(c)

	storedHash, err := h.Service.Repo.ExistingPassword(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			models.ErrorResponse(err.Error()))
	}

	if err := utilities.CheckPasswordHash(req.OldPassword, storedHash); err != nil {
		return c.JSON(http.StatusUnauthorized,
			models.ErrorResponse("Old password is incorrect"))
	}

	if err := utilities.CheckPasswordHash(req.NewPassword, storedHash); err == nil {
		return c.JSON(http.StatusUnauthorized,
			models.ErrorResponse("Reusing the old password not allowed"))
	}

	er := h.Service.ChangePassword(userID, req.NewPassword)

	if er != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse(er.Error()))
	}

	return c.JSON(http.StatusOK, models.SuccessResponse("Password changed successfully", nil))
}
