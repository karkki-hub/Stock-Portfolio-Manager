package models

type Watchlist struct {
	ID      uint `json:"id"`
	UserID  uint `json:"user_id"`
	StockID uint `json:"stock_id"`
}

type WatchlistStock struct {
	Symbol    string  `json:"symbol"`
	StockName string  `json:"stock_name"`
	LastPrice float64 `json:"last_price"`
}
