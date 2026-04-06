package repository

import (
	"database/sql"

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

func (r *ProfileRepository) Update(userID uint, phone string, name string, address string) error {
	query := `UPDATE users SET phone = ?, name = ?, Address = ?  WHERE user_id = ?`
	_, err := r.DB.Exec(query, phone, name, address, userID)
	return err
}

func (r *ProfileRepository) ExistingPassword(userID uint) (string, error) {
	query := `SELECT password FROM users WHERE user_id = ?`

	var Password string

	err := r.DB.QueryRow(query, userID).Scan(
		&Password,
	)

	if err != nil {
		return "", err
	}
	return Password, err
}

func (r *ProfileRepository) ResetPassword(userID uint, password string) error {
	query := `UPDATE users SET password = ?  WHERE user_id = ?`
	_, err := r.DB.Exec(query, password, userID)
	return err
}

func (r *ProfileRepository) GetAllUserId() ([]models.UserID, error) {
	query := `SELECT user_id, name FROM users`

	var userIDs []models.UserID

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u models.UserID

		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}

		userIDs = append(userIDs, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userIDs, nil
}
