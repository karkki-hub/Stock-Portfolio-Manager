package services

import (
	"karkki-hub/Stock-Portfolio-Manager/internal/models"
	"karkki-hub/Stock-Portfolio-Manager/internal/repository"
	// "karkki-hub/Stock-Portfolio-Manager/internal/utilities"
)

type ReportService struct {
	Repo *repository.ReportRepository
}

func NewReportService(r *repository.ReportRepository) *ReportService {
	return &ReportService{Repo: r}
}

func (s *ReportService) GetReport(userID uint) (*models.Report, error) {
	return s.Repo.GetReportById(userID)
}
