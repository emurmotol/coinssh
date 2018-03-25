package actions

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/emurmotol/coinssh/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

const WebTokenName = "_web_token"

type WebLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// WebGetLogin default implementation.
func WebGetLogin(c buffalo.Context) error {
	if IsAccountLoggedIn(c.Session()) {
		return c.Redirect(http.StatusFound, "/dashboard")
	}

	c.Set("webLoginRequest", &WebLoginRequest{})

	return c.Render(http.StatusOK, r.HTML("web/auth/login.html", WebAuthLayout))
}

// WebPostLogin default implementation.
func WebPostLogin(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return errors.WithStack(errors.New("No transaction found"))
	}

	req := &WebLoginRequest{}

	if err := c.Bind(req); err != nil {
		return errors.WithStack(err)
	}

	q := tx.Where(fmt.Sprintf("email = '%s'", req.Email))
	account := &models.Account{}

	if err := q.First(account); err != nil {
		return errors.WithStack(err)
	}

	if err := comparePassword(account.Password, req.Password); err != nil {
		return errors.WithStack(err)
	}

	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		Id:        account.ID.String(),
	}

	signingKey, err := ioutil.ReadFile(os.Getenv("JWT_KEY_PATH"))

	if err != nil {
		return errors.WithStack(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(signingKey)

	if err != nil {
		return errors.WithStack(err)
	}

	c.Session().Set(WebTokenName, tokenString)

	return c.Redirect(http.StatusFound, "/dashboard")
}

// WebGetLogout default implementation.
func WebGetLogout(c buffalo.Context) error {
	c.Session().Delete(WebTokenName)
	c.Session().Clear()

	return c.Redirect(http.StatusFound, "/login")
}

// WebGetHome is a default handler to serve up
// a home page.
func WebGetHome(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("index.html"))
}

// WebGetDashboard default implementation.
func WebGetDashboard(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("web/dashboard/index.html", WebLayout))
}

func IsAccountLoggedIn(s *buffalo.Session) bool {
	sessionToken := s.Get(WebTokenName)

	if sessionToken == nil {
		return false
	}

	tokenString := s.Get(WebTokenName).(string)

	if len(tokenString) == 0 {
		return false
	}

	return true
}
