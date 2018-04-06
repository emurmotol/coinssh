package actions

import (
	"net/http"

	"github.com/emurmotol/coinssh/external"
	"github.com/emurmotol/coinssh/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"
)

// AdminGetLogin default implementation.
func AdminGetLogin(c buffalo.Context) error {
	if IsUserLoggedIn(c.Session()) {
		return c.Redirect(http.StatusFound, "/admin/dashboard")
	}

	c.Set("user", &models.User{})

	return c.Render(http.StatusOK, r.HTML("admin/auth/login.html", AdminAuthLayout))
}

func AdminPostLogin(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return errors.WithStack(errors.New(T.Translate(c, "tx.not.ok")))
	}

	user := &models.User{}

	if err := c.Bind(user); err != nil {
		return errors.WithStack(err)
	}
	c.Set("user", user)

	vErrs := validate.NewErrors()
	errKey := "loginErrors"
	back := func(key string, with map[string][]string) error {
		c.Set(key, with)
		return c.Render(http.StatusUnprocessableEntity, r.HTML("admin/auth/login.html", AdminAuthLayout))
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

	vErrs, err = user.ValidateLogin(tx)

	if err != nil {
		return errors.WithStack(err)
	}

	if vErrs.HasAny() {
		return back("errors", vErrs.Errors)
	}

	if err := user.Authorize(tx); err != nil {
		vErrs.Add(errKey, err.Error())
	}

	if vErrs.HasAny() {
		return back(errKey, vErrs.Errors)
	}
	tokenString, err := makeToken(user.ID.String())

	if err != nil {
		return errors.WithStack(err)
	}

	c.Session().Set(AdminTokenName, tokenString)

	return c.Redirect(http.StatusFound, "/admin/dashboard")
}

func AdminGetLogout(c buffalo.Context) error {
	c.Session().Delete(AdminTokenName)
	c.Session().Clear()

	return c.Redirect(http.StatusFound, "/admin/login")
}

func IsUserLoggedIn(s *buffalo.Session) bool {
	sessionToken := s.Get(AdminTokenName)

	if sessionToken == nil {
		return false
	}
	tokenString := sessionToken.(string)

	if len(tokenString) == 0 {
		return false
	}
	return true
}

// AdminGetDashboard default implementation.
func AdminGetDashboard(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("admin/dashboard/index.html", AdminLayout))
}
