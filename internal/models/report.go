package models

type UserReport struct {
	Symbol          string  `json:"symbol"`
	StockName       string  `json:"stock_name"`
	Qty             float64 `json:"quantity"`
	AvgBuyPrice     float64 `json:"avg_buy_price"`
	CurrentPrice    float64 `json:"current_price"`
	TotalInvestment float64 `json:"total_investment"`
	CurrentValue    float64 `json:"current_value"`
	ProfitLoss      float64 `json:"profit_loss"`
}

type Report struct {
	Name            string       `json:"name"`
	TotalInvestment float64      `json:"tot_investment"`
	TotCurrentValue float64      `json:"tot_cur_investment"`
	TotalProfitLoss float64      `json:"total_profit_loss"`
	StocksOwned     []UserReport `json:"stocks_owned"`
}
