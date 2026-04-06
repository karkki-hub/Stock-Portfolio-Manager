package handlers

import (
	"net/http"
	"os"

	"fmt"
	"time"

	"karkki-hub/Stock-Portfolio-Manager/internal/services"
	"karkki-hub/Stock-Portfolio-Manager/internal/utilities"

	"github.com/labstack/echo/v4"
)

type ReportHandler struct {
	Service *services.ReportService
	Profile *services.ProfileService
	Cron    *services.CronService
}

func NewReportHandler(s *services.ReportService, p *services.ProfileService, c *services.CronService) *ReportHandler {
	return &ReportHandler{Service: s, Profile: p, Cron: c}
}

func (h *ReportHandler) ExportReportCSV(c echo.Context) error {

	userID := getUserID(c)

	report, err := h.Service.GetReport(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// Set headers
	c.Response().Header().Set("Content-Type", "text/csv")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%d-report.csv", userID))

	// Write directly to response
	err = utilities.WriteReportCSV(c.Response(), report)
	if err != nil {
		err = h.Service.LogReport(fmt.Sprintf("%d-report.csv", userID), "download", "FAILED")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	err = h.Service.LogReport(fmt.Sprintf("%d-report.csv", userID), "download", "SUCCESS")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return nil
}

func (h *ReportHandler) DailyReport() error {

	users, err := h.Profile.GetAllUserIDs()
	if err != nil {
		return err
	}

	for _, user := range users {

		report, err := h.Service.GetReport(user.ID)
		if err != nil {
			return err
		}

		filepath := fmt.Sprintf(
			"C:\\Users\\karkki\\Desktop\\reports\\%d-%s-%s.csv",
			user.ID,
			user.Name,
			time.Now().Format("2006-01-02"),
		)

		filename := fmt.Sprintf("%d-%s-%s.csv", user.ID, user.Name, time.Now().Format("2006-01-02"))
		file, err := os.Create(filepath)
		if err != nil {
			return err
		}

		if err := utilities.WriteReportCSV(file, report); err != nil {
			h.Service.LogReport(filename, "daily", "FAILED")
			h.Cron.CreateLog("Daily Report", "FAILED", fmt.Sprintf("Failed to generate report for user %d: %s", user.ID, err.Error()))
			file.Close()
			return err
		}

		err = h.Service.LogReport(filename, "daily", "SUCCESS")
		if err != nil {
			return err
		}
		file.Close()
	}

	h.Cron.CreateLog("Daily Report", "SUCCESS", "Daily reports generated successfully")

	return nil
}
