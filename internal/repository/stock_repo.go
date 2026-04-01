package repository

import (
	"database/sql"
	"fmt"
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
	query := `SELECT stock_id, symbol, stock_name, last_price FROM stocks WHERE symbol = ?`
	row := r.DB.QueryRow(query, symbol)

	var stock models.Stock
	err := row.Scan(&stock.ID, &stock.Symbol, &stock.StockName, &stock.LastPrice)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		fmt.Printf("Error fetching stock by symbol: %v\n", err)
		return nil, err
	}
	return &stock, nil
}

func (r *StockRepository) SearchByKeyword(keyword string) ([]models.StockDetails, error) {
	query := `SELECT symbol, stock_name FROM stock_list 
	WHERE symbol LIKE ? OR stock_name LIKE ? LIMIT 10`
	search := "%" + keyword + "%"
	rows, err := r.DB.Query(query, search, search)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []models.StockDetails

	for rows.Next() {
		var stock models.StockDetails
		err := rows.Scan(&stock.Symbol, &stock.StockName)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}

func (r *StockRepository) GetStockName(symbol string) string {
	query := `SELECT stock_name FROM stock_list WHERE symbol = ?`
	row := r.DB.QueryRow(query, symbol)
	var name string
	err := row.Scan(&name)
	if err != nil {
		return ""
	}
	return name
}
