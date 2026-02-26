package services

import (
	"errors"

	"karkki-hub/Stock-Portfolio-Manager/internal/models"
	"karkki-hub/Stock-Portfolio-Manager/internal/repository"
	utils "karkki-hub/Stock-Portfolio-Manager/internal/utilitis"
)

func Register(email, password string) error {
	hashed, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := models.User{
		ID:       uint(len(email)), // simple test ID
		Email:    email,
		Password: hashed,
	}

	repository.CreateUser(user)
	return nil
}

func Login(email, password string) (string, error) {
	user := repository.GetUserByEmail(email)
	if user == nil {
		return "", errors.New("user not found")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid password")
	}

	return utils.GenerateJWT(user.ID)
}
