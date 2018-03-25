package mailers

import (
	"github.com/gobuffalo/buffalo/mail"
	"github.com/gobuffalo/buffalo/render"
	"github.com/pkg/errors"
	"github.com/emurmotol/coinssh/models"
)

func SendRegisterActivation(a *models.Account) error {
	m := mail.NewMessage()

	// fill in with your stuff:
	m.Subject = "Register Activation"
	m.From = NoReplyEmail
	m.To = []string{a.Email}
	err := m.AddBody(r.HTML("register_activation.html"), render.Data{
		"account": a,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return smtp.Send(m)
}
