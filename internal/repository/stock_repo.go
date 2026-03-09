package repository

import (
	"database/sql"
	"karkki-hub/Stock-Portfolio-Manager/internal/models"
)

type StockRepository struct {
	DB *sql.DB
}

func NewStockRepository(db *sql.DB) *StockRepository {
	return &StockRepository{DB: db}
}

func (r *StockRepository) Save(stock *models.Stock) error {
	query := `INSERT INTO stocks (symbol, stock_name, last_price) 
	VALUES (?, ?, ?) 
	ON DUPLICATE KEY UPDATE 
	last_price = VALUES(last_price),
	last_updated = CURRENT_TIMESTAMP`
	_, err := r.DB.Exec(query, stock.Symbol, stock.StockName, stock.LastPrice)
	return err
}

func (r *StockRepository) GetBySymbol(symbol string) (*models.Stock, error) {
	query := `SELECT stock_id, symbol, stock_name, last_price, last_updated, created_at FROM stocks WHERE symbol = ?`
	row := r.DB.QueryRow(query, symbol)

	var stock models.Stock
	err := row.Scan(&stock.ID, &stock.Symbol, &stock.StockName, &stock.LastPrice, &stock.LastUpdated, &stock.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return &stock, nil
}
