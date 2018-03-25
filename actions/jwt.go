package actions

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"io/ioutil"
	"os"
	"github.com/pkg/errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/gobuffalo/pop"
	"github.com/emurmotol/coinssh/models"
	"github.com/gobuffalo/buffalo"
)

func makeToken(id string) (string, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		Id:        id,
	}

	signingKey, err := ioutil.ReadFile(os.Getenv("JWT_KEY_PATH"))

	if err != nil {
		return "", errors.WithStack(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(signingKey)
}

func authenticated(tokenName string, c buffalo.Context) (interface{}, error) {
	sessionToken := c.Session().Get(tokenName)
	emptyTokenErr := fmt.Errorf("No token set in session")

	if sessionToken == nil {
		return nil, emptyTokenErr
	}
	tokenString := sessionToken.(string)

	if len(tokenString) == 0 {
		return nil, emptyTokenErr
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

	if err != nil {
		return nil, fmt.Errorf("Could not parse the token: %v", err)
	}

	// Getting claims
	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {

		logrus.Errorf("Claims: %v", claims)

		// Get the DB connection from the context
		tx, ok := c.Value("tx").(*pop.Connection)

		if !ok {
			return nil, errors.New("No transaction found")
		}
		var model interface{}

		if tokenName == AdminTokenName {
			model = &models.User{}
		} else if tokenName == WebTokenName {
			model = &models.Account{}
		} else {
			return nil, fmt.Errorf("Could not identify the %s", tokenName)
		}

		// Retrieving user from db
		if err := tx.Find(model, claims["jti"].(string)); err != nil {
			return nil, fmt.Errorf("Could not identify the %s: %v", tokenName, err)
		}
		return model, nil
	}
	return nil, fmt.Errorf("Failed to validate token: %v", claims)
}
