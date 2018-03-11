package actions

import (
	"github.com/gobuffalo/buffalo"
	"net/http"
)

// GetHome is a default handler to serve up
// a home page.
func GetHome(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("index.html"))
}
