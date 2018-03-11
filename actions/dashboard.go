package actions

import (
	"github.com/gobuffalo/buffalo"
	"net/http"
)

// AdminDashboard default implementation.
func AdminDashboard(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("admin/dashboard/index.html", AdminLayout))
}
