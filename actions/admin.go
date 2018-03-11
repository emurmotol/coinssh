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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AdminLogin default implementation.
func AdminGetLogin(c buffalo.Context) error {
	c.Set("loginRequest", &LoginRequest{})

	return c.Render(http.StatusOK, r.HTML("admin/auth/login.html", AdminAuthLayout))
}

func AdminPostLogin(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
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

	signingKey, err := ioutil.ReadFile(os.Getenv("ADMIN_JWT_KEY_PATH"))

	if err != nil {
		return errors.WithStack(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(signingKey)

	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("user", user)
	c.Set("token", tokenString)

	return c.Render(http.StatusOK, r.HTML("admin/auth/ok.html", AdminAuthLayout))
}
