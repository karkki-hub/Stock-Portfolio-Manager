package services

import (
	"errors"
	"fmt"
	"time"

	"karkki-hub/Stock-Portfolio-Manager/internal/models"
	"karkki-hub/Stock-Portfolio-Manager/internal/repository"
	utilities "karkki-hub/Stock-Portfolio-Manager/pkg/utilities"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	UserRepo  *repository.UserRepository
	JWTSecret []byte
}

func NewAuthService(repo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		UserRepo:  repo,
		JWTSecret: []byte(jwtSecret),
	}
}

func (s *AuthService) Register(name, email, phone, address, password, apiKey string) (*models.User, error) {

	existingUser, _ := s.UserRepo.GetByEmail(email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := utilities.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Phone:    phone,
		Address:  address,
		Password: string(hashedPassword),
		APIToken: apiKey,
	}

	err = s.UserRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (string, error) { // Authenticate a user and return a JWT token

	user, err := s.UserRepo.GetByEmail(email)
	if err != nil {
		fmt.Println("Login failed: user not found for email:", email)
		return "", errors.New("invalid credentials")
	}

	err = utilities.CheckPasswordHash(password, user.Password) // Compare the provided password with the stored hash
	if err != nil {
		fmt.Println("Login failed: password mismatch for email:", email)
		return "", errors.New("invalid credentials")
	}

	fmt.Println("Login successful for email:", email)

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(6 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //	Create a new JWT token with the user claims and sign it with the secret key

	return token.SignedString(s.JWTSecret)
}
