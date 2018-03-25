package actions

import (
	"net/http"
	"github.com/emurmotol/coinssh/models"
	"github.com/gobuffalo/buffalo"
)

func AdminMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		model, err := authenticated(c, AdminTokenName)

		if err != nil {
			if c.Request().Header.Get("X-Requested-With") == "xmlhttprequest" {
				return c.Error(http.StatusUnauthorized, err)
			}

			return c.Redirect(http.StatusFound, "/admin/logout")
		}

		if model != nil {
			c.Set("authUser", model.(*models.User))
		}
		return next(c)
	}
}

func WebMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		model, err := authenticated(c, WebTokenName)

		if err != nil {
			if c.Request().Header.Get("X-Requested-With") == "xmlhttprequest" {
				return c.Error(http.StatusUnauthorized, err)
			}

			return c.Redirect(http.StatusFound, "/logout")
		}

		if model != nil {
			c.Set("authAccount", model.(*models.Account))
		}
		return next(c)
	}
}
