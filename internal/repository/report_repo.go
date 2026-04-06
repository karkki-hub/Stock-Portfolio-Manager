package repository

import (
	"database/sql"

	"karkki-hub/Stock-Portfolio-Manager/internal/models"
)

type ReportRepository struct {
	DB *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{DB: db}
}

func (r *ReportRepository) GetReportById(userID uint) (*models.Report, error) {
	query := `
SELECT 
	s.symbol,
	s.stock_name,
	p.qty,
	p.avg_buy_price,
	s.last_price,
	p.tot_investment
FROM portfolios p
JOIN stocks s ON p.stock_id = s.stock_id
WHERE p.user_id = ?`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []models.UserReport
	var totalInvestment, totalCurrentValue float64

	for rows.Next() {
		var stock models.UserReport
		err := rows.Scan(
			&stock.Symbol,
			&stock.StockName,
			&stock.Qty,
			&stock.AvgBuyPrice,
			&stock.CurrentPrice,
			&stock.TotalInvestment,
		)
		if err != nil {
			return nil, err
		}
		stock.CurrentValue = stock.Qty * stock.CurrentPrice
		stock.ProfitLoss = stock.CurrentValue - stock.TotalInvestment
		stocks = append(stocks, stock)
		totalInvestment += stock.TotalInvestment
		totalCurrentValue += stock.CurrentValue
	}

	totalProfitLoss := totalCurrentValue - totalInvestment

	var name string
	err = r.DB.QueryRow(`SELECT name FROM users WHERE user_id = ?`, userID).Scan(&name)
	if err != nil {
		return nil, err
	}

	report := &models.Report{
		Name:            name,
		TotalInvestment: totalInvestment,
		TotCurrentValue: totalCurrentValue,
		TotalProfitLoss: totalProfitLoss,
		StocksOwned:     stocks,
	}
	return report, nil
}

func (r *ReportRepository) LogReport(filename string, action string, status string) error {
	query := `INSERT INTO reports (report_type, file_name, status) VALUES (?, ?, ?)`
	_, err := r.DB.Exec(query, filename, action, status)
	return err
}
