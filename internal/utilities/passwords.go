package utilities

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func IsValidPassword(password string) bool {
	// Simple password validation - at least 8 characters, one uppercase, one lowercase, one digit
	if len(password) < 8 {
		return false
	}
	hasUpper := false
	hasLower := false
	hasDigit := false
	for _, char := range password {
		hasUpper = hasUpper || (char >= 'A' && char <= 'Z')
		hasLower = hasLower || (char >= 'a' && char <= 'z')
		hasDigit = hasDigit || (char >= '0' && char <= '9')
	}
	return hasUpper && hasLower && hasDigit
}
