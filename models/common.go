package models

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"github.com/gobuffalo/pop"
	"github.com/emurmotol/coinssh/external"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

func encryptPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)

	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(hash), nil
}

type EmailIsDisposable struct {
	Field string
	Name  string
	tx    *pop.Connection
	lang  *Lang
}

func (v *EmailIsDisposable) IsValid(errors *validate.Errors) {
	lang := v.lang
	T := lang.T
	c := lang.C
	yes, _ := external.IsEmailDisposable(v.Field)

	if yes {
		errors.Add(validators.GenerateKey(v.Name), T.Translate(c, "email.is.disposable"))
	}
}
