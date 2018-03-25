package models

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func encryptPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)

	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(hash), nil
}
