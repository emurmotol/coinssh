package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/emurmotol/coinssh/models"
	"github.com/pkg/errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
	"os"
	"io/ioutil"
	"net/http"
)

const AdminTokenName = "_admin_token"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// GetAdminLogin default implementation.
func GetAdminLogin(c buffalo.Context) error {
	if AdminIsLoggedIn(c.Session()) {
		return c.Redirect(http.StatusFound, "/admin/dashboard")
	}

	c.Set("loginRequest", &LoginRequest{})

	return c.Render(http.StatusOK, r.HTML("admin/auth/login.html", AdminAuthLayout))
}

func PostAdminLogin(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return errors.WithStack(errors.New("No transaction found"))
	}

	req := &LoginRequest{}

	if err := c.Bind(req); err != nil {
		return errors.WithStack(err)
	}

	q := tx.Where(fmt.Sprintf("email = '%s'", req.Email))
	user := &models.User{}

	if err := q.First(user); err != nil {
		return errors.WithStack(err)
	}

	if err := comparePassword(user.Password, req.Password); err != nil {
		return errors.WithStack(err)
	}

	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		Id:        user.ID.String(),
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

	c.Session().Set(AdminTokenName, tokenString)

	return c.Redirect(http.StatusFound, "/admin/dashboard")
}

func GetAdminLogout(c buffalo.Context) error {
	c.Session().Delete(AdminTokenName)
	c.Session().Clear()

	return c.Redirect(http.StatusFound, "/admin/login")
}

func AdminIsLoggedIn(s *buffalo.Session) bool {
	sessionToken := s.Get(AdminTokenName)

	if sessionToken == nil {
		return false
	}

	tokenString := s.Get(AdminTokenName).(string)

	if len(tokenString) == 0 {
		return false
	}

	return true
}
