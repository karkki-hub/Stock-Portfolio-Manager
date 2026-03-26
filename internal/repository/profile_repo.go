package repository

import (
	"database/sql"

	// "fmt"
	"karkki-hub/Stock-Portfolio-Manager/internal/models"
)

type ProfileRepository struct {
	DB *sql.DB
}

func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{DB: db}
}

func (r *ProfileRepository) GetUserById(userID uint) (*models.Profile, error) {
	query := `SELECT user_id, email, password, name, phone, address FROM users WHERE user_id = ?`

	var user models.Profile

	err := r.DB.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Phone,
		&user.Address,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *ProfileRepository) Update(userID uint, phone string, email string, address string) error {
	query := `UPDATE users SET phone = ?, email = ?, Address = ?  WHERE user_id = ?`
	_, err := r.DB.Exec(query, phone, email, address, userID)
	return err
}
