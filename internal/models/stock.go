package models

type Stock struct {
	Symbol    string  `json:"symbol"`
	StockName string  `json:"stock_name"`
	LastPrice float64 `json:"last_price"`
}
