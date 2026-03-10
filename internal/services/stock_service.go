package services

import (
	// "database/sql"
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

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	time.Sleep(2 * time.Second)

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

	if len(quote) == 0 {
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

	// fmt.Println(string(body)) // debug

	var result struct {
		Name string `json:"Name"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return symbol, err
	}

	// fmt.Printf("Decoded: %+v\n", result)

	if result.Name == "" {
		return symbol, nil
	}

	return result.Name, nil
}
