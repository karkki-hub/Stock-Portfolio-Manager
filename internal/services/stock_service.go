package services

import (
	// "database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"karkki-hub/Stock-Portfolio-Manager/internal/models"
	"karkki-hub/Stock-Portfolio-Manager/internal/repository"
)

type StockService struct {
	Repo     *repository.StockRepository
	AlphaKey string
}

func NewStockService(repo *repository.StockRepository, alphaKey string) *StockService {
	return &StockService{Repo: repo, AlphaKey: alphaKey}
}

func (s *StockService) SearchStock(symbol string) (*models.Stock, error) {
	stock, err := s.Repo.GetBySymbol(symbol)
	if err != nil {
		return nil, err
	}
	if stock != nil {
		return stock, nil
	}
	// if err == sql.ErrNoRows {
	// 	return nil, nil
	// }

	url := fmt.Sprintf("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s&apikey=%s", symbol, s.AlphaKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	quote := result["Global Quote"]

	price, _ := strconv.ParseFloat(quote["05. price"], 64)

	stock = &models.Stock{
		Symbol:    quote["01. symbol"],
		StockName: quote["01. symbol"], // Alpha Vantage doesn't provide name, using symbol as placeholder
		LastPrice: price,
	}

	err = s.Repo.Save(stock)
	if err != nil {
		return nil, err
	}

	return stock, nil
}
