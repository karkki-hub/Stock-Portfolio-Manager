package models

type Stock struct {
	ID        uint    `json:"-"`
	Symbol    string  `json:"symbol"`
	StockName string  `json:"stock_name"`
	LastPrice float64 `json:"last_price"`
}
type StockDetails struct {
	Symbol    string `json:"symbol"`
	StockName string `json:"stock_name"`
}
