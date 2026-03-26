package services

import (
	"karkki-hub/Stock-Portfolio-Manager/internal/models"
	"karkki-hub/Stock-Portfolio-Manager/internal/repository"
)

type ProfileService struct {
	Repo *repository.ProfileRepository
}

func NewProfileService(r *repository.ProfileRepository) *ProfileService {
	return &ProfileService{Repo: r}
}

func (s *ProfileService) GetProfile(userID uint) (*models.Profile, error) {
	return s.Repo.GetUserById(userID)
}

func (s *ProfileService) ChangeProfile(userID uint, phone string, email string, address string) error {
	return s.Repo.Update(userID, phone, email, address)
}
