package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

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

func (s *StockService) AddStock(symbol string, name string) error {

	url := fmt.Sprintf(
		"https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&apikey=%s",
		symbol,
		s.AlphaKey,
	)

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var check map[string]interface{}
	_ = json.Unmarshal(body, &check)

	if msg, ok := check["Note"]; ok {
		return fmt.Errorf("alphavantage rate limit: %v", msg)
	}
	if msg, ok := check["Information"]; ok {
		return fmt.Errorf("alphavantage rate limit: %v", msg)
	}

	type TimeSeriesResponse struct {
		TimeSeries map[string]map[string]string `json:"Time Series (Daily)"`
	}

	var data TimeSeriesResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	if data.TimeSeries == nil {
		return fmt.Errorf("no time series data for %s", symbol)
	}

	stock := &models.Stock{
		Symbol:    symbol,
		StockName: name,
		LastPrice: 0,
	}

	err = s.Repo.Save(stock)
	if err != nil {
		return err
	}

	savedStock, err := s.Repo.GetBySymbol(symbol)
	if err != nil {
		return err
	}
	stockID := savedStock.ID

	var latestDate time.Time
	var latestPrice float64

	for dateStr, values := range data.TimeSeries {

		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}

		price, err := strconv.ParseFloat(values["4. close"], 64)
		if err != nil {
			continue
		}

		err = s.Repo.UpdateHistory(stockID, price, date)
		if err != nil {
			fmt.Println("History insert error:", err)
			continue
		}

		if date.After(latestDate) {
			latestDate = date
			latestPrice = price
		}
	}

	if latestPrice == 0 {
		return fmt.Errorf("failed to extract latest price")
	}

	err = s.Repo.UpdateStockPrice(stockID, latestPrice)
	if err != nil {
		return err
	}

	return nil
}

func (s *StockService) GetStockName(keyword string) ([]models.StockDetails, error) {
	time.Sleep(1 * time.Second)

	url := fmt.Sprintf(
		"https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=%s&apikey=%s",
		keyword,
		s.AlphaKey,
	)

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var check map[string]interface{}
	_ = json.Unmarshal(body, &check)

	if msg, ok := check["Information"]; ok {
		return nil, fmt.Errorf("alphavantage rate limit: %v", msg)
	}
	if msg, ok := check["Note"]; ok {
		return nil, fmt.Errorf("alphavantage rate limit: %v", msg)
	}

	type AlphaResponse struct {
		BestMatches []struct {
			Symbol string `json:"1. symbol"`
			Name   string `json:"2. name"`
		} `json:"bestMatches"`
	}

	var apiResp AlphaResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}

	var result []models.StockDetails

	for _, match := range apiResp.BestMatches {
		result = append(result, models.StockDetails{
			Symbol:    match.Symbol,
			StockName: match.Name,
		})
	}

	return result, nil
}

func (s *StockService) SearchStocksByKeyword(keyword string) ([]models.StockDetails, error) {
	if len(keyword) < 3 {
		return nil, fmt.Errorf("keyword must be at least 3 characters")
	}
	stock, err := s.Repo.SearchByKeyword(keyword)
	if err != nil {
		return s.GetStockName(keyword)
	}
	if stock == nil {
		return s.GetStockName(keyword)
	}
	return stock, nil
}

func (s *StockService) GetOrCreateStock(symbol string) (*models.Stock, error) {
	stock, err := s.Repo.GetBySymbol(symbol)
	if err == nil && stock != nil {
		return stock, nil
	}

	name := s.Repo.GetStockName(symbol)

	if err := s.AddStock(symbol, name); err != nil {
		return nil, err
	}

	stock, err = s.Repo.GetBySymbol(symbol)
	if err != nil {
		return nil, err
	}

	return stock, nil
}
