package models

import "time"

type Transaction struct {
	ID        uint      `json:"-"`
	UserID    uint      `json:"-"`
	StockID   uint      `json:"-"`
	Symbol    string    `json:"symbol"`
	Type      string    `json:"type"`
	Quantity  float64   `json:"quantity"`
	Price     float64   `json:"price"`
	TotalAmt  float64   `json:"total_amount"`
	CreatedAt time.Time `json:"created_at"`
}
