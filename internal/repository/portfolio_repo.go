package repository

import (
	"database/sql"
	"fmt"
	"karkki-hub/Stock-Portfolio-Manager/internal/models"
	"karkki-hub/Stock-Portfolio-Manager/internal/utilities"
)

type PortfolioRepository struct {
	DB *sql.DB
}

func NewPortfolioRepository(db *sql.DB) *PortfolioRepository {
	return &PortfolioRepository{DB: db}
}

func (r *PortfolioRepository) Update(userID uint, stockID uint, qty float64, price float64) error {
	query := `
INSERT INTO portfolios (user_id, stock_id, qty, avg_buy_price, tot_investment)
VALUES (?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE 
avg_buy_price = (tot_investment + VALUES(tot_investment)) / (qty + VALUES(qty)),
tot_investment = tot_investment + VALUES(tot_investment),
qty = qty + VALUES(qty)`

	fmt.Print(price)

	_, err := r.DB.Exec(
		query,
		userID,
		stockID,
		qty,
		price,
		qty*price,
	)
	return err
}

func (r *PortfolioRepository) Sell(userID uint, stockID uint, qty float64, price float64) error {
	query := `
UPDATE portfolios 
SET qty = qty - ?, tot_investment = tot_investment - ( ? * ? ),
UPDATED_AT = CURRENT_TIMESTAMP
WHERE user_id = ? AND stock_id = ?`

	_, err := r.DB.Exec(
		query,
		qty,
		qty,
		price,
		userID,
		stockID,
	)
	return err
}

func (r *PortfolioRepository) CheckStock() error {
	query := `DELETE FROM portfolios WHERE qty = 0`

	_, nil := r.DB.Exec(
		query,
	)
	return nil
}

func (r *PortfolioRepository) GetByUser(userID uint) ([]models.Portfolio, error) {
	query := `
SELECT 
	p.user_id,
	p.stock_id,
	s.symbol,
	p.qty,
	p.avg_buy_price,
	p.tot_investment,
	s.last_price,
	p.updated_at
FROM portfolios p
JOIN stocks s ON p.stock_id = s.stock_id
WHERE p.user_id = ?`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var portfolios []models.Portfolio

	for rows.Next() {
		var p models.Portfolio
		err := rows.Scan(
			&p.UserID,
			&p.StockID,
			&p.Symbol,
			&p.Quantity,
			&p.AvgBuyPrice,
			&p.TotalInvest,
			&p.CurrentPrice,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		p.CurrentValue = utilities.RoundUp(p.CurrentPrice * p.Quantity)
		p.ProfitLoss = utilities.RoundUp(p.CurrentValue - p.TotalInvest)

		portfolios = append(portfolios, p)
	}

	return portfolios, nil
}
