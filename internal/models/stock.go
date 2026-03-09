package models

import "time"

type Stock struct {
	ID          int       `json:"id"`
	Symbol      string    `json:"symbol"`
	StockName   string    `json:"stock_name"`
	LastPrice   float64   `json:"last_price"`
	LastUpdated time.Time `json:"last_update"`
	CreatedAt   time.Time `json:"created_at"`
}
