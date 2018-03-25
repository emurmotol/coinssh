package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// WebGetLogin default implementation.
func WebGetLogin(c buffalo.Context) error {
	return c.Render(200, r.HTML("web/get_login.html"))
}

// WebPostLogin default implementation.
func WebPostLogin(c buffalo.Context) error {
	return c.Render(200, r.HTML("web/post_login.html"))
}

// WebGetLogout default implementation.
func WebGetLogout(c buffalo.Context) error {
	return c.Render(200, r.HTML("web/get_logout.html"))
}

// WebGetHome is a default handler to serve up
// a home page.
func WebGetHome(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("index.html"))
}

// WebGetDashboard default implementation.
func WebGetDashboard(c buffalo.Context) error {
	return c.Render(200, r.HTML("web/get_dashboard.html"))
}
