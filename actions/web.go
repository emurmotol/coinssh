package actions

import (
	"net/http"
	"github.com/emurmotol/coinssh/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"github.com/emurmotol/coinssh/mailers"
	"github.com/gobuffalo/validate"
	"fmt"
)

const WebTokenName = "_web_token"

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
		return errors.WithStack(errors.New("No transaction found"))
	}

	account := &models.Account{}

	if err := c.Bind(account); err != nil {
		return errors.WithStack(err)
	}

	if err := account.Authorize(tx); err != nil {
		c.Set("account", account)
		verrs := validate.NewErrors()
		verrs.Add("Login", "Invalid email or password.")
		c.Set("errors", verrs.Errors)
		return c.Render(http.StatusUnprocessableEntity, r.HTML("web/auth/login.html", WebAuthLayout))
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
		return errors.WithStack(errors.New("No transaction found"))
	}

	// Allocate an empty User
	account := &models.Account{}
	// Bind user to the html form elements
	if err := c.Bind(account); err != nil {
		return errors.WithStack(err)
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(account)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		c.Set("account", account)
		// Make the errors available inside the html template
		c.Set("errors", verrs.Errors)
		// Render again the register.html template that the user can
		// correct the input.
		return c.Render(http.StatusUnprocessableEntity, r.HTML("web/auth/register.html", WebAuthLayout))
	}
	go mailers.SendRegisterActivation(account)

	// If there are no errors set a success message
	c.Flash().Add("success", fmt.Sprintf("Hello, %s! We sent you an email. Please activate your account.", account.Name))
	// and redirect to the home page
	return c.Redirect(http.StatusFound, "/register")
}
