package handlers

import (
	"encoding/csv"
	"net/http"
	"os"

	"fmt"

	// "karkki-hub/Stock-Portfolio-Manager/internal/models"
	"karkki-hub/Stock-Portfolio-Manager/internal/services"

	"github.com/labstack/echo/v4"
)

type ReportHandler struct {
	Service *services.ReportService
}

func NewReportHandler(s *services.ReportService) *ReportHandler {
	return &ReportHandler{Service: s}
}

func (h *ReportHandler) ExportReportCSV(c echo.Context) error {

	userID := getUserID(c)

	report, err := h.Service.GetReport(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	filepath := "C:\\Users\\karkki\\Desktop\\reports\\report.csv"

	file, err := os.Create(filepath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	defer file.Close()

	// Set headers for download
	c.Response().Header().Set("Content-Type", "text/csv")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=report.csv")

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{
		"",
		report.Name,
		"",
	})

	// ✅ Header row
	writer.Write([]string{
		"Symbol",
		"Stock Name",
		"Quantity",
		"Avg Buy Price",
		"Current Price",
		"Total Investment",
		"Current Value",
		"Profit/Loss",
	})

	// ✅ Data rows
	for _, s := range report.StocksOwned {
		row := []string{
			s.Symbol,
			s.StockName,
			fmt.Sprintf("%.2f", s.Qty),
			fmt.Sprintf("%.2f", s.AvgBuyPrice),
			fmt.Sprintf("%.2f", s.CurrentPrice),
			fmt.Sprintf("%.2f", s.TotalInvestment),
			fmt.Sprintf("%.2f", s.CurrentValue),
			fmt.Sprintf("%.2f", s.ProfitLoss),
		}
		writer.Write(row)
	}

	writer.Write([]string{
		"",
	})

	// ✅ Total row (very important)
	writer.Write([]string{
		"TOTAL",
		"",
		"",
		"",
		"",
		fmt.Sprintf("%.2f", report.TotalInvestment),
		fmt.Sprintf("%.2f", report.TotCurrentValue),
		fmt.Sprintf("%.2f", report.TotalProfitLoss),
	})

	return nil
}
