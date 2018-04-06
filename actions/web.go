package actions

import (
	"net/http"

	"github.com/emurmotol/coinssh/external"
	"github.com/emurmotol/coinssh/mailers"
	"github.com/emurmotol/coinssh/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"
)

// WebGetLogin default implementation.
func WebGetLogin(c buffalo.Context) error {
	if IsAccountLoggedIn(c.Session()) {
		return c.Redirect(http.StatusFound, "/dashboard")
	}

	c.Set("account", &models.Account{})

	return c.Render(http.StatusOK, r.HTML("web/auth/login.html", WebAuthLayout))
}

// WebPostLogin default implementation.
func WebPostLogin(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return errors.WithStack(errors.New(T.Translate(c, "tx.not.ok")))
	}

	account := &models.Account{}

	if err := c.Bind(account); err != nil {
		return errors.WithStack(err)
	}
	c.Set("account", account)

	vErrs := validate.NewErrors()
	errKey := "loginErrors"
	back := func(key string, with map[string][]string) error {
		c.Set(key, with)
		return c.Render(http.StatusUnprocessableEntity, r.HTML("web/auth/login.html", WebAuthLayout))
	}

	isHuman, err := external.IsHuman(c.Request())

	if err != nil {
		return errors.WithStack(err)
	}

	if !isHuman {
		vErrs.Add(errKey, T.Translate(c, "verify.human"))
	}

	if vErrs.HasAny() {
		return back(errKey, vErrs.Errors)
	}

	vErrs, err = account.ValidateLogin(tx)

	if err != nil {
		return errors.WithStack(err)
	}

	if vErrs.HasAny() {
		return back("errors", vErrs.Errors)
	}

	if err := account.Authorize(tx); err != nil {
		vErrs.Add(errKey, err.Error())
	}

	if vErrs.HasAny() {
		return back(errKey, vErrs.Errors)
	}
	tokenString, err := makeToken(account.ID.String())

	if err != nil {
		return errors.WithStack(err)
	}

	c.Session().Set(WebTokenName, tokenString)

	return c.Redirect(http.StatusFound, "/dashboard")
}

// WebGetLogout default implementation.
func WebGetLogout(c buffalo.Context) error {
	c.Session().Delete(WebTokenName)
	c.Session().Clear()

	return c.Redirect(http.StatusFound, "/login")
}

// WebGetHome is a default handler to serve up
// a home page.
func WebGetHome(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("index.html"))
}

// WebGetDashboard default implementation.
func WebGetDashboard(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("web/dashboard/index.html", WebLayout))
}

func IsAccountLoggedIn(s *buffalo.Session) bool {
	sessionToken := s.Get(WebTokenName)

	if sessionToken == nil {
		return false
	}
	tokenString := sessionToken.(string)

	if len(tokenString) == 0 {
		return false
	}
	return true
}

// WebGetRegister default implementation.
func WebGetRegister(c buffalo.Context) error {
	if IsAccountLoggedIn(c.Session()) {
		return c.Redirect(http.StatusFound, "/dashboard")
	}

	c.Set("account", &models.Account{})

	return c.Render(http.StatusOK, r.HTML("web/auth/register.html", WebAuthLayout))
}

// WebPostRegister default implementation.
func WebPostRegister(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return errors.WithStack(errors.New(T.Translate(c, "tx.not.ok")))
	}

	// Allocate an empty User
	account := &models.Account{}
	// Bind user to the html form elements
	if err := c.Bind(account); err != nil {
		return errors.WithStack(err)
	}
	account.Lang = &models.Lang{C: c, T: T}
	c.Set("account", account)

	vErrs := validate.NewErrors()
	errKey := "registerErrors"
	back := func(key string, with map[string][]string) error {
		c.Set(key, with)
		return c.Render(http.StatusUnprocessableEntity, r.HTML("web/auth/register.html", WebAuthLayout))
	}

	isHuman, err := external.IsHuman(c.Request())

	if err != nil {
		return errors.WithStack(err)
	}

	if !isHuman {
		vErrs.Add(errKey, T.Translate(c, "verify.human"))
	}

	if vErrs.HasAny() {
		return back(errKey, vErrs.Errors)
	}

	// Validate the data from the html form
	vErrs, err = tx.ValidateAndCreate(account)
	if err != nil {
		return errors.WithStack(err)
	}

	if err != nil {
		return errors.WithStack(err)
	}

	if vErrs.HasAny() {
		return back("errors", vErrs.Errors)
	}
	go mailers.SendRegisterActivation(account)

	// If there are no errors set a success message
	c.Flash().Add("success", T.Translate(c, "register.activation.sent", account))
	// and redirect to the home page
	return c.Redirect(http.StatusFound, "/login")
}
