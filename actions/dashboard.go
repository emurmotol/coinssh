package actions

import (
	"github.com/gobuffalo/buffalo"
)

// AdminDashboard default implementation.
func AdminDashboard(c buffalo.Context) error {
	return c.Render(200, r.HTML("admin/dashboard/index.html", AdminLayout))
}
