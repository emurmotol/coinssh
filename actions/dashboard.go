package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// AdminGetDashboard default implementation.
func AdminGetDashboard(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("admin/dashboard/index.html", AdminLayout))
}
