package repository

import (
	"database/sql"
	"karkki-hub/Stock-Portfolio-Manager/internal/models"
)

type WatchlistRepository struct {
	DB *sql.DB
}

func NewWatchlistRepository(db *sql.DB) *WatchlistRepository {
	return &WatchlistRepository{DB: db}
}

func (r *WatchlistRepository) Add(userID, stockID uint) error {

	query := "INSERT INTO watchlist (user_id, stock_id) VALUES (?, ?)"
	_, err := r.DB.Exec(query, userID, stockID)
	return err
}

func (r *WatchlistRepository) Remove(userID uint, symbol string) error {
	query := `DELETE w FROM watchlist w
JOIN stocks s ON w.stock_id = s.stock_id
WHERE w.user_id = ? AND s.symbol = ?`
	_, err := r.DB.Exec(query, userID, symbol)
	return err
}

func (r *WatchlistRepository) GetByUser(userID uint) ([]models.WatchlistStock, error) {
	query := `SELECT s.symbol, s.stock_name, s.last_price
FROM watchlist w
JOIN stocks s ON w.stock_id = s.stock_id
WHERE w.user_id = ?`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []models.WatchlistStock
	for rows.Next() {
		var stock models.WatchlistStock
		err := rows.Scan(&stock.Symbol, &stock.StockName, &stock.LastPrice)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}
	return stocks, nil
}

func (r *WatchlistRepository) GetStockHistory(stockid uint) ([]models.StockPriceHistory, error) {
	query := `SELECT price_date, closing_price FROM stock_price_history WHERE stock_id = ? ORDER BY price_date DESC`
	rows, err := r.DB.Query(query, stockid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.StockPriceHistory
	for rows.Next() {
		var entry models.StockPriceHistory
		err := rows.Scan(&entry.Date, &entry.Price)
		if err != nil {
			return nil, err
		}
		history = append(history, entry)
	}
	return history, nil
}
