package actions

import (
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
)

var r *render.Engine
var assetsBox = packr.NewBox("../public")

const ApplicationLayout = "application.html"
const AdminLayout = "admin/layout/admin.html"
const AdminAuthLayout = "admin/auth/layout/auth.html"
const WebLayout = "web/layout/admin.html"
const WebAuthLayout = "web/auth/layout/auth.html"

func init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: ApplicationLayout,

		// Box containing all of the templates:
		TemplatesBox: packr.NewBox("../templates"),
		AssetsBox:    assetsBox,

		// Add template helpers here:
		Helpers: render.Helpers{
			// uncomment for non-Bootstrap form helpers:
			// "form":     plush.FormHelper,
			// "form_for": plush.FormForHelper,
		},
	})
}
