package models

import "time"

type Transaction struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	StockID   uint      `json:"stock_id"`
	Type      string    `json:"type"` // "buy" or "sell"
	Quantity  float64   `json:"quantity"`
	Price     float64   `json:"price"`
	TotalAmt  float64   `json:"total_amount"`
	CreatedAt time.Time `json:"created_at"`
}
