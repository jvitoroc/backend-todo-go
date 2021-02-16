package util

import (
	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(passwordHash *string, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}
	*passwordHash = string(hash)
	return nil
}

func VerifyPasswordHash(passwordHash, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return false
	}

	return true
}
