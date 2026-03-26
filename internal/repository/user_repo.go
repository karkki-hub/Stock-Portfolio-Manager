package repository

import (
	"database/sql"
	"fmt"
	"karkki-hub/Stock-Portfolio-Manager/internal/models"
	"strings"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
	INSERT INTO users (name, email, phone, address, password, api_token)
	VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.DB.Exec(
		query,
		user.Name,
		user.Email,
		user.Phone,
		user.Address,
		user.Password,
		user.APIToken,
	)
	return err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	email = strings.TrimSpace(email)
	fmt.Println("Searching user with email:", email)

	query := `
	SELECT user_id, name, email, phone, address, password, api_token, created_at
	FROM users WHERE LOWER(email) = LOWER(?)`

	row := r.DB.QueryRow(query, email)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Address,
		&user.Password,
		&user.APIToken,
		&user.CreatedAt,
	)

	if err != nil {
		fmt.Println("Error scanning user:", err)
		return nil, err
	}

	return &user, nil
}
