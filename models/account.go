package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
)

type Account struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
}

// String is not required by pop and may be deleted
func (a Account) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Accounts is not required by pop and may be deleted
type Accounts []Account

// String is not required by pop and may be deleted
func (a Accounts) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *Account) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: a.Name, Name: "Name"},
		&validators.EmailIsPresent{Field: a.Email, Name: "Email"},
		&validators.StringIsPresent{Field: a.Password, Name: "Password"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *Account) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *Account) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (a *Account) BeforeCreate(tx *pop.Connection) error {
	password, err := encryptPassword(a.Password)

	if err != nil {
		return errors.WithStack(err)
	}

	a.Password = password

	return nil
}
