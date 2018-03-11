package models

import (
	"golang.org/x/crypto/bcrypt"
	"github.com/pkg/errors"
)

func encryptPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)

	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(hash), nil
}