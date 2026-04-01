package services

import (
	"karkki-hub/Stock-Portfolio-Manager/internal/repository"
)

type WatchlistService struct {
	WatchRepo    *repository.WatchlistRepository
	StockRepo    *repository.StockRepository
	StockService *StockService
}

func NewWatchlistService(w *repository.WatchlistRepository, s *repository.StockRepository, ss *StockService) *WatchlistService {
	return &WatchlistService{WatchRepo: w, StockRepo: s, StockService: ss}
}

func (s *WatchlistService) Add(userID uint, symbol string) error {
	stock, err := s.StockService.GetOrCreateStock(symbol)
	if err != nil {
		return err
	}

	return s.WatchRepo.Add(userID, stock.ID)
}

func (s *WatchlistService) Remove(userID uint, symbol string) error {
	return s.WatchRepo.Remove(userID, symbol)
}

func (s *WatchlistService) Get(userID uint) (interface{}, error) {
	return s.WatchRepo.GetByUser(userID)
}
