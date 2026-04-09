package services

import (
	"karkki-hub/Stock-Portfolio-Manager/internal/models"
	"karkki-hub/Stock-Portfolio-Manager/internal/repository"
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

func (s *ReportService) LogReport(filename string, action string, status string) error {
	return s.Repo.LogReport(filename, action, status)
}

func (s *ReportService) ListReports(userID uint) ([][]string, error) {
	return s.Repo.ListReports(userID)
}
