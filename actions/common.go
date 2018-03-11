package actions

import (
	"golang.org/x/crypto/bcrypt"
)

func comparePassword(hash string, pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
}