package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	symbol = strings.ToUpper(symbol)

	// Check DB first
	stock, err := s.Repo.GetBySymbol(symbol)
	if err != nil {
		return nil, err
	}
	if stock != nil {
		return stock, nil
	}

	url := fmt.Sprintf(
		"https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s&apikey=%s",
		symbol,
		s.AlphaKey,
	)

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read full body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Rate-limit check
	var check map[string]interface{}
	json.Unmarshal(body, &check)
	if msg, ok := check["Note"]; ok {
		return nil, fmt.Errorf("alphavantage rate limit: %v", msg)
	}
	if msg, ok := check["Information"]; ok {
		return nil, fmt.Errorf("alphavantage rate limit: %v", msg)
	}

	// Decode normally
	var result map[string]map[string]string
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	quote, ok := result["Global Quote"]
	if !ok || len(quote) == 0 {
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
		"https://www.alphavantage.co/query?function=OVERVIEW&symbol=%s&apikey=%s",
		symbol,
		s.AlphaKey,
	)

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return symbol, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return symbol, err
	}

	var check map[string]interface{}
	json.Unmarshal(body, &check)

	// Detect rate-limit message
	if msg, ok := check["Information"]; ok {
		return symbol, fmt.Errorf("alphavantage rate limit: %v", msg)
	}
	if msg, ok := check["Note"]; ok {
		return symbol, fmt.Errorf("alphavantage rate limit: %v", msg)
	}

	// Decode normally if no rate-limit
	var result struct {
		Name string `json:"Name"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return symbol, err
	}

	if result.Name == "" {
		return symbol, nil
	}

	return result.Name, nil
}
