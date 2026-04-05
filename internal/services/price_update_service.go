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
	StockRepo *repository.StockRepository
}

func NewPriceService(stockRepo *repository.StockRepository) *PriceService {
	return &PriceService{StockRepo: stockRepo}
}

// ✅ Move struct outside
type TimeSeriesResponse struct {
	TimeSeries map[string]map[string]string `json:"Time Series (Daily)"`
}

func (s *PriceService) UpdatePrices() {

	stocks, err := s.StockRepo.GetAllStocks()
	if err != nil {
		fmt.Println("Error fetching stocks:", err)
		return
	}

	for _, stock := range stocks {

		url := fmt.Sprintf(
			"https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&apikey=YOUR_API_KEY",
			stock.Symbol,
		)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("HTTP error:", err)
			continue
		}

		// ✅ Always close body safely
		defer resp.Body.Close()

		// ✅ Check API response status
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

		for dateStr, values := range data.TimeSeries {

			// ✅ Parse date
			date, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				continue
			}

			// ✅ Parse close price
			closePrice, err := strconv.ParseFloat(values["4. close"], 64)
			if err != nil {
				continue
			}

			// ✅ Insert history
			err = s.StockRepo.UpdateHistory(stock.ID, closePrice, date)
			if err != nil {
				fmt.Println("History insert error:", err)
			}

			// ✅ Track latest
			if date.After(latestDate) {
				latestDate = date
				latestClose = closePrice
			}
		}

		// ✅ Update latest stock price
		if latestClose > 0 {
			err = s.StockRepo.UpdateStockPrice(stock.ID, latestClose)
			if err != nil {
				fmt.Println("Update price error:", err)
			}
		}

		// ✅ Rate limit protection
		time.Sleep(12 * time.Second)
	}

	fmt.Println("✅ Prices & History Updated")
}
