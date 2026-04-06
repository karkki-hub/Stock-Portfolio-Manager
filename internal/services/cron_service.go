package services

import (
	// "karkki-hub/Stock-Portfolio-Manager/internal/models"
	"karkki-hub/Stock-Portfolio-Manager/internal/repository"
	// "karkki-hub/Stock-Portfolio-Manager/internal/utilities"
)

type CronService struct {
	Repo *repository.CronRepository
}

func NewCronService(r *repository.CronRepository) *CronService {
	return &CronService{Repo: r}
}

func (s *CronService) CreateLog(job, status, message string) error {
	return s.Repo.CreateLog(job, status, message)
}
