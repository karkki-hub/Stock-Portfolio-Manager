package services

import (
	"karkki-hub/Stock-Portfolio-Manager/internal/models"
	"karkki-hub/Stock-Portfolio-Manager/internal/repository"
	"karkki-hub/Stock-Portfolio-Manager/pkg/utilities"
)

type PortfolioService struct {
	Repo *repository.PortfolioRepository
}

func NewPortfolioService(r *repository.PortfolioRepository) *PortfolioService {
	return &PortfolioService{Repo: r}
}

func (s *PortfolioService) Buy(userID uint, stockID uint, qty float64, price float64) error {
	return s.Repo.Update(userID, stockID, qty, price)
}

func (s *PortfolioService) Sell(userID uint, stockID uint, qty float64) error {

	var sell error

	sell = s.Repo.Sell(userID, stockID, qty)

	s.Repo.CheckStock()

	return sell
}

func (s *PortfolioService) Get(userID uint) (*models.PortfolioSummary, error) {
	stocks, err := s.Repo.GetByUser(userID)

	if err != nil {
		return nil, err
	}

	var totalInvestment float64
	var currentval float64

	for _, p := range stocks {
		totalInvestment += p.TotalInvest
		currentval += p.CurrentValue
	}

	summary := &models.PortfolioSummary{
		TotalInvestment: utilities.RoundUp(totalInvestment),
		TotCurrentValue: utilities.RoundUp(currentval),
		TotalProfitLoss: utilities.RoundUp(currentval - totalInvestment),
		Stocks:          stocks,
	}

	return summary, nil
}
