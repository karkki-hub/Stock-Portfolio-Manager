package services

import (
	"encoding/json"
	"fmt"
	"karkki-hub/Stock-Portfolio-Manager/internal/repository"
	"net/http"
	"strconv"
	"time"
)

type PriceService struct {
	StockRepo   *repository.StockRepository
	CronService *CronService
}

func NewPriceService(stockRepo *repository.StockRepository, cronService *CronService) *PriceService {
	return &PriceService{StockRepo: stockRepo, CronService: cronService}
}

type TimeSeriesResponse struct {
	TimeSeries map[string]map[string]string `json:"Time Series (Daily)"`
}

func (s *PriceService) UpdatePrices() {

	if !MarketclosingWindow() { // Check if the market is currently open, and if so, skip the price update
		fmt.Println("Market open, skipping price update")
		return
	}

	stocks, err := s.StockRepo.GetAllStocks()
	if err != nil {
		fmt.Println("Error fetching stocks:", err)
		return
	}
	today := time.Now().Format("2006-01-02") // Get the current date in YYYY-MM-DD format for price update checks

	for _, stock := range stocks {

		check, err := s.StockRepo.PriceUpdateCheck(stock.ID, today)
		if err != nil {
			fmt.Println("Price check error:", err)
			continue
		}

		if check {
			fmt.Printf("Already updated: %s\n", stock.Symbol)
			continue
		}

		url := fmt.Sprintf( // Fetch the latest stock price data from the Alpha Vantage API for the given stock symbol
			"https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&apikey=YOUR_API_KEY",
			stock.Symbol,
		)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("HTTP error:", err)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Bad response for:", stock.Symbol)
			continue
		}

		var data TimeSeriesResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			fmt.Println("JSON decode error:", err)
			continue
		}

		if data.TimeSeries == nil {
			fmt.Println("No data for:", stock.Symbol)
			continue
		}

		var latestDate time.Time
		var latestClose float64

		for dateStr, values := range data.TimeSeries { // Iterate through the time series data and update history and current price

			date, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				continue
			}

			closePrice, err := strconv.ParseFloat(values["4. close"], 64)
			if err != nil {
				continue
			}

			err = s.StockRepo.UpdateHistory(stock.ID, closePrice, date)
			if err != nil {
				fmt.Println("History insert error")
			}

			if date.After(latestDate) {
				latestDate = date
				latestClose = closePrice
			}
		}

		if latestClose > 0 {
			err = s.StockRepo.UpdateStockPrice(stock.ID, latestClose)
			if err != nil {
				fmt.Println("Update price error:", err)
			}
		}

		time.Sleep(12 * time.Second)
	}

	fmt.Println("Prices & History Updated")
	s.CronService.CreateLog("Price Update", "SUCCESS", "Prices and history updated successfully")

}

func MarketclosingWindow() bool { // Determine if the current time is within the market closing window (weekends and after-hours)
	now := time.Now()

	loc, _ := time.LoadLocation("Asia/Kolkata")
	now = now.In(loc)

	weekday := now.Weekday()
	hour := now.Hour()

	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}
	if hour >= 9 && hour <= 18 {
		return false
	}
	return true
}
