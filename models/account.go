package models

import (
	"encoding/json"
	"time"

	"database/sql"
	"strings"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware/i18n"
)

type Account struct {
	ID           uuid.UUID        `json:"id" db:"id"`
	CreatedAt    time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at" db:"updated_at"`
	Name         string           `json:"name" db:"name"`
	Username     string           `json:"username" db:"username"`
	Email        string           `json:"email" db:"email"`
	Password     string           `json:"-" db:"-"`
	PasswordHash string           `json:"-" db:"password_hash"`
	C            buffalo.Context  `json:"-" db:"-"`
	T            *i18n.Translator `json:"-" db:"-"`
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
	T := a.T
	c := a.C
	lang := &Lang{C: c, T: T}

	return validate.Validate(
		&validators.StringIsPresent{Field: a.Username, Name: "Username"},
		&validators.StringLengthInRange{
			Field:   a.Username,
			Name:    "Username",
			Min:     8,
			Message: T.Translate(c, "username.char.len", 8),
		},
		&validators.StringLengthInRange{
			Field:   a.Password,
			Name:    "Password",
			Min:     8,
			Message: T.Translate(c, "password.char.len", 8),
		},
		&validators.EmailIsPresent{Field: a.Email, Name: "Email"},
		&validators.StringIsPresent{Field: a.Password, Name: "Password"},
		&AccountUsernameIsTaken{Field: a.Username, Name: "Username", tx: tx, lang: lang},
		&AccountEmailIsTaken{Field: a.Email, Name: "Email", tx: tx, lang: lang},
		&EmailIsDisposable{Field: a.Email, Name: "Email", tx: tx, lang: lang},
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
	lang  *Lang
}

func (v *AccountEmailIsTaken) IsValid(errors *validate.Errors) {
	lang := v.lang
	T := lang.T
	c := lang.C
	q := v.tx.Where("email = ?", v.Field)
	m := Account{}
	err := q.First(&m)
	if err == nil {
		// found a account with the same email
		errors.Add(validators.GenerateKey(v.Name), T.Translate(c, "email.taken"))
	}
}

type AccountUsernameIsTaken struct {
	Field string
	Name  string
	tx    *pop.Connection
	lang  *Lang
}

func (v *AccountUsernameIsTaken) IsValid(errors *validate.Errors) {
	lang := v.lang
	T := lang.T
	c := lang.C
	q := v.tx.Where("username = ?", v.Field)
	m := Account{}
	err := q.First(&m)
	if err == nil {
		// found a account with the same username
		errors.Add(validators.GenerateKey(v.Name), T.Translate(c, "username.taken"))
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
	T := a.T
	c := a.C
	username := strings.ToLower(a.Username)
	err := tx.Where("email = ? or username = ?", username, username).First(a)

	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			// couldn't find an account with that email or username
			return errors.New(T.Translate(c, "account.not.found"))
		}
		return errors.WithStack(err)
	}
	// confirm that the given password matches the hashed password from the db
	err = bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(a.Password))
	if err != nil {
		return errors.New(T.Translate(c, "wrong.password"))
	}
	return nil
}
