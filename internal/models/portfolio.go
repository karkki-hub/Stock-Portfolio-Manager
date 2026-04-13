package models

import "time"

type Portfolio struct {
	UserID       uint      `json:"-"` // Exclude UserID from JSON responses
	StockID      uint      `json:"-"` // Exclude StockID from JSON responses
	Symbol       string    `json:"symbol"`
	Quantity     float64   `json:"quantity"`
	AvgBuyPrice  float64   `json:"avg_buy_price"`
	TotalInvest  float64   `json:"total_investment"`
	CurrentPrice float64   `json:"current_price"`
	CurrentValue float64   `json:"current_value"`
	ProfitLoss   float64   `json:"profit_loss"`
	UpdatedAt    time.Time `json:"-"`
}

type PortfolioSummary struct {
	TotalInvestment float64     `json:"tot_investment"`
	TotCurrentValue float64     `json:"tot_cur_investment"`
	TotalProfitLoss float64     `json:"total_profit_loss"`
	Stocks          []Portfolio `json:"stocks"`
}
