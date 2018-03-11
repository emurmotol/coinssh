package actions

import (
	"github.com/gobuffalo/buffalo"
	"net/http"
)

// GetAdminDashboard default implementation.
func GetAdminDashboard(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("admin/dashboard/index.html", AdminLayout))
}
