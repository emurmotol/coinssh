package models

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"encoding/json"
	"time"
)

func encryptPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)

	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(hash), nil
}

type kickbox struct {
	IsDisposable bool `json:"disposable"`
}

func getJson(url string, target interface{}) error {
	var httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
