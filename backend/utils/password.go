package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword — bcrypt( password + pepper )
func HashPassword(plain, pepper string) (string, error) {
	sum := plain + pepper
	hash, err := bcrypt.GenerateFromPassword([]byte(sum), 12)
	return string(hash), err
}

// CheckPassword — сравнение ввода с хешем
func CheckPassword(plain, hash, pepper string) error {
	sum := plain + pepper
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(sum))
}
