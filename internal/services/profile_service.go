package services

import (
	"karkki-hub/Stock-Portfolio-Manager/internal/models"
	"karkki-hub/Stock-Portfolio-Manager/internal/repository"
	utilities "karkki-hub/Stock-Portfolio-Manager/pkg/utilities"
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

func (s *ProfileService) ChangeProfile(userID uint, phone string, name string, address string) error {
	return s.Repo.Update(userID, phone, name, address)
}

func (s *ProfileService) ChangePassword(userID uint, password string) error {
	hashedPassword, err := utilities.HashPassword(password) // Hash the new password before storing it in the database
	if err != nil {
		return err
	}

	return s.Repo.ResetPassword(userID, hashedPassword)
}

func (s *ProfileService) GetAllUserIDs() ([]models.UserID, error) {
	return s.Repo.GetAllUserId()
}
