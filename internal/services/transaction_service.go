package services

import (
	"errors"
	"karkki-hub/Stock-Portfolio-Manager/internal/models"
	"karkki-hub/Stock-Portfolio-Manager/internal/repository"
)

type TransactionService struct {
	Repo          *repository.TransactionRepository
	StockRepo     *repository.StockRepository
	Portfolioserv *PortfolioService
}

func NewTransactionService(r *repository.TransactionRepository, s *repository.StockRepository, p *PortfolioService) *TransactionService {
	return &TransactionService{Repo: r, StockRepo: s, Portfolioserv: p}
}

func (s *TransactionService) Buy(userID uint, Symbol string, quantity float64, price float64) error {
	stock, err := s.StockRepo.GetBySymbol(Symbol)
	if err != nil {
		return err
	}

	if stock == nil {
		return errors.New("stock not found")
	}

	tx := &models.Transaction{
		UserID:   userID,
		StockID:  stock.ID,
		Type:     "buy",
		Quantity: quantity,
		Price:    price,
	}
	err = s.Repo.Create(tx)
	if err != nil {
		return err
	}
	return s.Portfolioserv.Buy(userID, stock.ID, quantity, price)
}

func (s *TransactionService) Sell(userID uint, Symbol string, quantity float64, price float64) error {
	stock, err := s.StockRepo.GetBySymbol(Symbol)
	if err != nil {
		return err
	}

	if stock == nil {
		return errors.New("stock not found")
	}

	holding, err := s.Repo.GetHolding(userID, stock.ID)
	if err != nil {
		return err
	}

	if holding < int(quantity) {
		return errors.New("insufficient holdings")
	}

	tx := &models.Transaction{
		UserID:   userID,
		StockID:  stock.ID,
		Type:     "sell",
		Quantity: quantity,
		Price:    price,
	}
	err = s.Repo.Create(tx)
	if err != nil {
		return err
	}
	return s.Portfolioserv.Sell(userID, stock.ID, quantity)
}

func (s *TransactionService) History(userID uint) ([]*models.Transaction, error) {
	return s.Repo.GetByUserID(userID)
}
