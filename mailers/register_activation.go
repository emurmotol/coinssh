package mailers

import (
	"github.com/emurmotol/coinssh/models"
	"github.com/gobuffalo/buffalo/mail"
	"github.com/gobuffalo/buffalo/render"
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
		panic(err) // Must recover from this
	}

	if err := smtp.Send(m); err != nil {
		panic(err) // Must recover from this
	}
	return nil
}
