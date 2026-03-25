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

func (s *ProfileService) ChangeEmail(userID uint, email string) error {
	return s.Repo.UpdateEmail(userID, email)
}

func (s *ProfileService) ChangePhone(userID uint, phone string) error {
	return s.Repo.UpdatePhone(userID, phone)
}

func (s *ProfileService) ChangeAddress(userID uint, address string) error {
	return s.Repo.UpdateAddress(userID, address)
}
