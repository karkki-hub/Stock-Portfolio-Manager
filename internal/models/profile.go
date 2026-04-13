package models

type Profile struct {
	ID       uint   `json:"-"` // Exclude ID from JSON responses
	Email    string `json:"email"`
	Password string `json:"-"` // Exclude password from JSON responses
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

type UserID struct {
	ID   uint
	Name string
}
