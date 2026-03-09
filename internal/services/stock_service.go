package services

import (
	"database/sql"
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

	if err == nil && stock != nil {
		return stock, nil
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	url := fmt.Sprintf(
		"https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s&apikey=%s",
		symbol,
		s.AlphaKey,
	)

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

	if note, ok := result["Note"]; ok {
		return nil, fmt.Errorf("alphavantage rate limit: %v", note)
	}

	quote := result["Global Quote"]

	if quote == nil || len(quote) == 0 {
		return nil, fmt.Errorf("no quote data for %s", symbol)
	}

	price, err := strconv.ParseFloat(quote["05. price"], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid price format")
	}

	name, err := s.GetStockName(symbol)
	if err != nil {
		name = symbol
	}

	stock = &models.Stock{
		Symbol:    quote["01. symbol"],
		StockName: name,
		LastPrice: price,
	}

	err = s.Repo.Save(stock)
	if err != nil {
		return nil, err
	}

	return stock, nil
}

func (s *StockService) GetStockName(symbol string) (string, error) {
	url := fmt.Sprintf(
		"https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=%s&apikey=%s",
		symbol,
		s.AlphaKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		BestMatches []map[string]string `json:"bestMatches"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	if len(result.BestMatches) == 0 {
		return symbol, nil // fallback
	}

	fmt.Printf("Best matches for %s: %+v\n", symbol, result.BestMatches)

	return result.BestMatches[0]["2. name"], nil
}
