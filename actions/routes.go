package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// GetRoutes is a default handler to serve up
// a routes page.
func GetRoutes(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("routes.html"))
}
