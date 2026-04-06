package utilities

import (
	"encoding/csv"
	"fmt"
	"io"

	"karkki-hub/Stock-Portfolio-Manager/internal/models"
)

func WriteReportCSV(w io.Writer, report *models.Report) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Title
	if err := writer.Write([]string{"", "NAME:", report.Name, ""}); err != nil {
		return err
	}

	// Header
	if err := writer.Write([]string{
		"Symbol",
		"Stock Name",
		"Quantity",
		"Avg Buy Price",
		"Current Price",
		"Total Investment",
		"Current Value",
		"Profit/Loss",
	}); err != nil {
		return err
	}

	// Data
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

		if err := writer.Write(row); err != nil {
			return err
		}
	}

	// Empty row
	_ = writer.Write([]string{""})

	// Total
	if err := writer.Write([]string{
		"TOTAL",
		"",
		"",
		"",
		"",
		fmt.Sprintf("%.2f", report.TotalInvestment),
		fmt.Sprintf("%.2f", report.TotCurrentValue),
		fmt.Sprintf("%.2f", report.TotalProfitLoss),
	}); err != nil {
		return err
	}

	return writer.Error()
}
