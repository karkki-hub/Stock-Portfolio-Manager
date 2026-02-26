package repository

import "karkki-hub/Stock-Portfolio-Manager/internal/models"

var users = []models.User{}

func CreateUser(user models.User) {
	users = append(users, user)
}

func GetUserByEmail(email string) *models.User {
	for _, u := range users {
		if u.Email == email {
			return &u
		}
	}
	return nil
}
