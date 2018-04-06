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
	"golang.org/x/tools/go/gcimporter15/testdata"
)

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	Password     string    `json:"-" db:"-"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Lang         *Lang     `json:"-" db:"-"`
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	lang := u.Lang

	return validate.Validate(
		&validators.StringIsPresent{Field: u.Name, Name: "Name"},
		&validators.EmailIsPresent{Field: u.Email, Name: "Email"},
		&validators.StringIsPresent{Field: u.Password, Name: "Password"},
		&UserEmailIsTaken{Field: u.Email, Name: "Email", tx: tx, lang: lang},
		&EmailIsDisposable{Field: u.Email, Name: "Email", tx: tx, lang: lang},
	), nil
}

func (u *User) ValidateLogin(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.EmailIsPresent{Field: u.Email, Name: "Email"},
		&validators.StringIsPresent{Field: u.Password, Name: "Password"},
	), nil
}

type UserEmailIsTaken struct {
	Field string
	Name  string
	tx    *pop.Connection
	lang  *Lang
}

func (v *UserEmailIsTaken) IsValid(errors *validate.Errors) {
	lang := v.lang
	T := lang.T
	c := lang.C
	q := v.tx.Where("email = ?", v.Field)
	m := User{}
	err := q.First(&m)
	if err == nil {
		// found a user with the same email
		errors.Add(validators.GenerateKey(v.Name), T.Translate(c, "email.taken"))
	}
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (u *User) BeforeCreate(tx *pop.Connection) error {
	u.Email = strings.ToLower(u.Email)
	password, err := encryptPassword(u.Password)

	if err != nil {
		return errors.WithStack(err)
	}
	u.PasswordHash = password

	return nil
}

func (u *User) Authorize(tx *pop.Connection) error {
	lang := u.Lang
	T := lang.T
	c := lang.C
	err := tx.Where("email = ?", strings.ToLower(u.Email)).First(u)

	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			// couldn't find an account with that email or username
			return errors.New(T.Translate(c, "user.not.found"))
		}
		return errors.WithStack(err)
	}
	// confirm that the given password matches the hashed password from the db
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))
	if err != nil {
		return errors.New(T.Translate(c, "wrong.password"))
	}
	return nil
}
