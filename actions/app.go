package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/ssl"
	"github.com/gobuffalo/envy"
	"github.com/unrolled/secure"

	"github.com/emurmotol/coinssh/models"
	"github.com/gobuffalo/buffalo/middleware/csrf"
	"github.com/gobuffalo/buffalo/middleware/i18n"
	"github.com/gobuffalo/packr"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator

const (
	WebTokenName = "_web_token"
	AdminTokenName = "_admin_token"
	CoinsshSessionName = "_coinssh_session"
)

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: CoinsshSessionName,
		})
		// Automatically redirect to SSL
		app.Use(ssl.ForceSSL(secure.Options{
			SSLRedirect:     ENV == "production",
			SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		}))

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		// Remove to disable this.
		app.Use(middleware.PopTransaction(models.DB))

		// Setup and use translations:
		var err error
		if T, err = i18n.New(packr.NewBox("../locales"), "en-US"); err != nil {
			app.Stop(err)
		}
		app.Use(T.Middleware())

		app.GET("/routes", GetRoutes)

		web := app.Group("/")
		web.Use(WebMiddleware)
		aR := AccountsResource{}
		web.Middleware.Skip(WebMiddleware, WebGetHome, WebGetLogin, WebPostLogin, WebGetLogout, WebGetRegister, WebPostRegister)
		web.GET("/", WebGetHome)
		web.GET("/login", WebGetLogin)
		web.POST("/login", WebPostLogin)
		web.GET("/logout", WebGetLogout)
		web.GET("/dashboard", WebGetDashboard)
		web.GET("/register", WebGetRegister)
		web.POST("/register", WebPostRegister)
		web.Resource("/accounts", aR)
		// Test if must clear middleware after above lines

		admin := app.Group("/admin")
		admin.Use(AdminMiddleware)
		uR := UsersResource{}
		admin.Middleware.Skip(AdminMiddleware, AdminGetLogin, AdminPostLogin, AdminGetLogout)
		admin.GET("/login", AdminGetLogin)
		admin.POST("/login", AdminPostLogin)
		admin.GET("/logout", AdminGetLogout)
		admin.GET("/dashboard", AdminGetDashboard)
		admin.Resource("/users", uR)

		app.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	return app
}
