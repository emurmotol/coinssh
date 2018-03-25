package actions

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/emurmotol/coinssh/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/sirupsen/logrus"
)

func AdminMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		emptySessionTokenErr := c.Error(http.StatusUnauthorized, fmt.Errorf("No token set in session"))

		if !AdminIsLoggedIn(c.Session()) {

			if c.Request().Header.Get("X-Requested-With") == "xmlhttprequest" {
				return emptySessionTokenErr
			}

			return c.Redirect(http.StatusFound, "/admin/login")
		}

		tokenString := c.Session().Get(AdminTokenName).(string)

		if len(tokenString) == 0 {
			return emptySessionTokenErr
		}

		// Parsing token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// RSA key
			mySignedKey, err := ioutil.ReadFile(os.Getenv("JWT_KEY_PATH"))

			if err != nil {
				return nil, fmt.Errorf("Could not open jwt key: %v", err)
			}

			return mySignedKey, nil
		})

		// Token expired
		if err != nil {
			if c.Request().Header.Get("X-Requested-With") == "xmlhttprequest" {
				return c.Error(http.StatusUnauthorized, fmt.Errorf("Could not parse the token: %v", err))
			}

			return c.Redirect(http.StatusFound, "/admin/logout")
		}

		// Getting claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			logrus.Errorf("Claims: %v", claims)

			// Get the DB connection from the context
			tx, ok := c.Value("tx").(*pop.Connection)

			if !ok {
				return c.Error(http.StatusInternalServerError, fmt.Errorf("No transaction found"))
			}

			// Allocate an empty User
			user := &models.User{}

			// Retrieving user from db
			if err := tx.Find(user, claims["jti"].(string)); err != nil {
				return c.Error(http.StatusNotFound, err)
			}

			if err != nil {
				return c.Error(http.StatusUnauthorized, fmt.Errorf("Could not identify the user: %v", err))
			}

			c.Set("user", user)

		} else {
			return c.Error(http.StatusUnauthorized, fmt.Errorf("Failed to validate token: %v", claims))
		}

		return next(c)
	}
}
