package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
	"strings"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	Password     string    `json:"-" db:"-"`
	PasswordHash string    `json:"-" db:"password_hash"`
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
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Name, Name: "Name"},
		&validators.EmailIsPresent{Field: u.Email, Name: "Email"},
		&validators.StringIsPresent{Field: u.Password, Name: "Password"},
		&UserEmailNotTaken{Field: u.Email, Name: "Email", tx: tx},
	), nil
}

type UserEmailNotTaken struct {
	Field string
	Name  string
	tx    *pop.Connection
}

func (v *UserEmailNotTaken) IsValid(errors *validate.Errors) {
	q := v.tx.Where("email = ?", v.Field)
	m := User{}
	err := q.First(&m)
	if err == nil {
		// found a user with the same email
		errors.Add(validators.GenerateKey(v.Name), "A user with that email already exists.")
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
	err := tx.Where("email = ?", strings.ToLower(u.Email)).First(u)

	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			// couldn't find a user with that email address
			return errors.New("User not found.")
		}
		return errors.WithStack(err)
	}
	// confirm that the given password matches the hashed password from the db
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))
	if err != nil {
		return errors.New("Invalid password.")
	}
	return nil
}
