package actions

import (
	"net/http"
	"github.com/emurmotol/coinssh/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
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
		return errors.WithStack(errors.New("No transaction found"))
	}

	user := &models.User{}

	if err := c.Bind(user); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := user.ValidateLogin(tx)

	if err != nil {
		return errors.WithStack(err)
	}
	c.Set("user", user)

	if verrs.HasAny() {
		c.Set("errors", verrs.Errors)
		return c.Render(http.StatusUnprocessableEntity, r.HTML("admin/auth/login.html", AdminAuthLayout))
	}

	if err := user.Authorize(tx); err != nil {
		verrs.Add("loginErrors", "Invalid email or password.")
	}

	if verrs.HasAny() {
		c.Set("loginErrors", verrs.Errors)
		return c.Render(http.StatusUnprocessableEntity, r.HTML("admin/auth/login.html", AdminAuthLayout))
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
