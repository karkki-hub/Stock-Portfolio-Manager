package models

type Profile struct {
	ID       uint   `json:"-"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

type UserID struct {
	ID   uint
	Name string
}
