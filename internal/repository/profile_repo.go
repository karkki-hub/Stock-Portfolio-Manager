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
		// &user.UpdatedAT,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *ProfileRepository) UpdatePhone(userID uint, phone string) error {
	query := `UPDATE users SET phone = ? WHERE user_id = ?`
	_, err := r.DB.Exec(query, phone, userID)
	return err
}

func (r *ProfileRepository) UpdateEmail(userID uint, email string) error {
	query := `UPDATE users SET email = ? WHERE user_id = ?`
	_, err := r.DB.Exec(query, email, userID)
	return err
}

func (r *ProfileRepository) UpdateAddress(userID uint, address string) error {
	query := `UPDATE users SET Address = ? WHERE user_id = ?`
	_, err := r.DB.Exec(query, address, userID)
	return err
}
