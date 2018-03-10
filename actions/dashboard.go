package actions

import "github.com/gobuffalo/buffalo"

// DashboardIndex default implementation.
func DashboardIndex(c buffalo.Context) error {
	return c.Render(200, r.HTML("admin/dashboard/index.html", AdminLayout))
}
