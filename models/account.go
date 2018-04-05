package models

import (
	"encoding/json"
	"time"

	"database/sql"
	"strings"

	"github.com/emurmotol/coinssh/external"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID           uuid.UUID `json:"id" db:"id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	Name         string    `json:"name" db:"name"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	Password     string    `json:"-" db:"-"`
	PasswordHash string    `json:"-" db:"password_hash"`
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
		&validators.StringIsPresent{Field: a.Username, Name: "Username"},
		&validators.EmailIsPresent{Field: a.Email, Name: "Email"},
		&validators.StringIsPresent{Field: a.Password, Name: "Password"},
		&AccountUsernameIsTaken{Field: a.Username, Name: "Username", tx: tx},
		&AccountEmailIsTaken{Field: a.Email, Name: "Email", tx: tx},
		&AccountEmailIsDisposable{Field: a.Email, Name: "Email", tx: tx},
	), nil
}

func (a *Account) ValidateLogin(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: a.Username, Name: "Username"},
		&validators.StringIsPresent{Field: a.Password, Name: "Password"},
	), nil
}

type AccountEmailIsTaken struct {
	Field string
	Name  string
	tx    *pop.Connection
}

func (v *AccountEmailIsTaken) IsValid(errors *validate.Errors) {
	q := v.tx.Where("email = ?", v.Field)
	m := Account{}
	err := q.First(&m)
	if err == nil {
		// found a account with the same email
		errors.Add(validators.GenerateKey(v.Name), "Email already taken.")
	}
}

type AccountUsernameIsTaken struct {
	Field string
	Name  string
	tx    *pop.Connection
}

func (v *AccountUsernameIsTaken) IsValid(errors *validate.Errors) {
	q := v.tx.Where("username = ?", v.Field)
	m := Account{}
	err := q.First(&m)
	if err == nil {
		// found a account with the same email
		errors.Add(validators.GenerateKey(v.Name), "Username already taken.")
	}
}

type AccountEmailIsDisposable struct {
	Field string
	Name  string
	tx    *pop.Connection
}

func (v *AccountEmailIsDisposable) IsValid(errors *validate.Errors) {
	yes, _ := external.IsEmailDisposable(v.Field)

	if yes {
		errors.Add(validators.GenerateKey(v.Name), "Disposable email address are not allowed.")
	}
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
	a.Email = strings.ToLower(a.Email)
	a.Username = strings.ToLower(a.Username)
	password, err := encryptPassword(a.Password)

	if err != nil {
		return errors.WithStack(err)
	}

	a.PasswordHash = password

	return nil
}

func (a *Account) Authorize(tx *pop.Connection) error {
	username := strings.ToLower(a.Username)
	err := tx.Where("email = ? or username = ?", username, username).First(a)

	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			// couldn't find an account with that email address
			return errors.New("Couldn't find your account.")
		}
		return errors.WithStack(err)
	}
	// confirm that the given password matches the hashed password from the db
	err = bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(a.Password))
	if err != nil {
		return errors.New("Wrong password.")
	}
	return nil
}
