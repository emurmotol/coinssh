package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"net/http"
	"fmt"
	"os"
	"io/ioutil"
	"github.com/gobuffalo/pop"
	"github.com/emurmotol/coinssh/models"
)

func AdminMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {

		if !AdminIsLoggedIn(c.Session()) {
			return c.Redirect(http.StatusFound, "/admin/login")
		}

		tokenString := c.Session().Get(AdminTokenName).(string)

		if len(tokenString) == 0 {
			return c.Error(http.StatusUnauthorized, fmt.Errorf("No token set in session"))
		}

		// parsing token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// key
			mySignedKey, err := ioutil.ReadFile(os.Getenv("ADMIN_JWT_KEY_PATH"))

			if err != nil {
				return nil, fmt.Errorf("could not open jwt key, %v", err)
			}

			return mySignedKey, nil
		})

		if err != nil {
			return c.Error(http.StatusUnauthorized, fmt.Errorf("Could not parse the token, %v", err))
		}

		// getting claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			logrus.Errorf("claims: %v", claims)

			// Get the DB connection from the context
			tx, ok := c.Value("tx").(*pop.Connection)

			if !ok {
				return c.Error(http.StatusInternalServerError, fmt.Errorf("no transaction found"))
			}

			// Allocate an empty User
			user := &models.User{}

			// retrieving user from db
			if err := tx.Find(user, claims["jti"].(string)); err != nil {
				return c.Error(http.StatusNotFound, err)
			}

			if err != nil {
				return c.Error(http.StatusUnauthorized, fmt.Errorf("Could not identify the user"))
			}

			c.Set("user", user)

		} else {
			return c.Error(http.StatusUnauthorized, fmt.Errorf("Failed to validate token: %v", claims))
		}

		return next(c)
	}
}
