package models

import "time"

type Transaction struct {
	ID        uint      `json:"-"` // Exclude ID from JSON responses
	UserID    uint      `json:"-"` // Exclude UserID from JSON responses
	StockID   uint      `json:"-"` // Exclude StockID from JSON responses
	Symbol    string    `json:"symbol"`
	Type      string    `json:"type"`
	Quantity  float64   `json:"quantity"`
	Price     float64   `json:"price"`
	TotalAmt  float64   `json:"total_amount"`
	CreatedAt time.Time `json:"created_at"`
}
