package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	Password  string    `json:"-"`
	APIToken  string    `json:"api_token"`
	CreatedAt time.Time `json:"created_at"`
}
