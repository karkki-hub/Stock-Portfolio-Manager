package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"karkki-hub/Stock-Portfolio-Manager/internal/models"
	"karkki-hub/Stock-Portfolio-Manager/internal/services"
	"karkki-hub/Stock-Portfolio-Manager/pkg/utilities"
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
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error()))
	}

	c.Response().Header().Set("Content-Type", "text/csv")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%d-report.csv", userID))

	err = utilities.WriteReportCSV(c.Response(), report)
	if err != nil {
		err = h.Service.LogReport(fmt.Sprintf("%d-report.csv", userID), "download", "FAILED")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error()))
	}

	err = h.Service.LogReport(fmt.Sprintf("%d-report.csv", userID), "download", "SUCCESS")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error()))
	}

	return nil
}

func (h *ReportHandler) DailyReport() error {
	users, err := h.Profile.GetAllUserIDs()
	if err != nil {
		return err
	}

	basepath := "reports"
	if err := os.MkdirAll(basepath, os.ModePerm); err != nil {
		return err
	}

	for _, user := range users {
		var report *models.Report
		var err error

		// Retry fetching report 3 times
		for attempt := 1; attempt <= 3; attempt++ {
			report, err = h.Service.GetReport(user.ID)
			if err == nil {
				break
			}
			time.Sleep(2 * time.Second)
		}

		// Prepare safe filename
		filenameSafe := strings.ReplaceAll(user.Name, " ", "_")
		filename := fmt.Sprintf("%d-%s-%s.csv", user.ID, filenameSafe, time.Now().Format("2006-01-02"))
		filepath := fmt.Sprintf("%s/%s", basepath, filename)

		if err != nil {
			h.Service.LogReport(filename, "daily", "FAILED")
			h.Cron.CreateLog("Daily Report", "FAILED",
				fmt.Sprintf("Failed to generate report for user %d after retries: %s", user.ID, err.Error()))
			continue
		}

		file, err := os.Create(filepath)
		if err != nil {
			h.Service.LogReport(filename, "daily", "FAILED")
			h.Cron.CreateLog("Daily Report", "FAILED",
				fmt.Sprintf("Failed to create file for user %d: %s", user.ID, err.Error()))
			continue
		}
		defer file.Close()

		// Retry writing CSV 3 times
		for attempt := 1; attempt <= 3; attempt++ {
			err = utilities.WriteReportCSV(file, report)
			if err == nil {
				break
			} else if attempt == 3 {
				h.Service.LogReport(filename, "daily", "FAILED")
				h.Cron.CreateLog("Daily Report", "FAILED",
					fmt.Sprintf("Failed to write report for user %d after retries: %s", user.ID, err.Error()))
			} else {
				time.Sleep(2 * time.Second)
			}
		}

		if err == nil {
			if logErr := h.Service.LogReport(filename, "daily", "SUCCESS"); logErr != nil {
				h.Cron.CreateLog("Daily Report", "FAILED",
					fmt.Sprintf("Failed to log report success for user %d: %s", user.ID, logErr.Error()))
			}
		}
	}

	h.Cron.CreateLog("Daily Report", "SUCCESS", "Daily reports generated successfully")
	return nil
}

// func (h *ReportHandler) ListReports(c echo.Context) error {
// 	userID := getUserID(c)
// 	basepath := "reports"

// 	files, err := os.ReadDir(basepath)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error()))
// 	}

// 	var result []string

// 	for _, file := range files {
// 		if strings.HasPrefix(file.Name(), fmt.Sprintf("%d-", userID)) {
// 			result = append(result, file.Name())
// 		}
// 	}

// 	return c.JSON(http.StatusOK, models.SuccessResponse("reports fetched", result))
// }

func (h *ReportHandler) DownloadReport(c echo.Context) error {

	filename := c.Param("filename")

	basepath := "reports"
	filepath := fmt.Sprintf("%s/%s", basepath, filename)

	// 🔒 security: prevent ../ attacks
	if strings.Contains(filename, "..") {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse("invalid filename"))
	}

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return c.JSON(http.StatusNotFound, models.ErrorResponse("file not found"))
	}

	return c.Attachment(filepath, filename)
}

func (h *ReportHandler) ListReports(c echo.Context) error {
	userID := getUserID(c)

	var result [][]string
	result, err := h.Service.ListReports(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, models.SuccessResponse("reports fetched", result))
}
