package actions

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/emurmotol/coinssh/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

func makeToken(id string) (string, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		Id:        id,
	}

	signingKey, err := ioutil.ReadFile(envy.Get("JWT_KEY_PATH", ""))

	if err != nil {
		return "", errors.WithStack(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(signingKey)
}

func authenticated(c buffalo.Context, tokenName string) (interface{}, error) {
	sessionToken := c.Session().Get(tokenName)
	emptyTokenErr := fmt.Errorf(T.Translate(c, "jwt.token.empty"))

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
			return nil, fmt.Errorf(T.Translate(c, "jwt.signing.error", token.Header["alg"]))
		}

		// RSA key
		mySignedKey, err := ioutil.ReadFile(envy.Get("JWT_KEY_PATH", ""))

		if err != nil {
			return nil, fmt.Errorf(T.Translate(c, "jwt.open.error", err))
		}

		return mySignedKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf(T.Translate(c, "jwt.parse.error", err))
	}

	// Getting claims
	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		// Get the DB connection from the context
		tx, ok := c.Value("tx").(*pop.Connection)

		if !ok {
			return nil, errors.New(T.Translate(c, "tx.not.ok"))
		}
		var model interface{}

		if tokenName == AdminTokenName {
			model = &models.User{}
		} else if tokenName == WebTokenName {
			model = &models.Account{}
		} else {
			return nil, fmt.Errorf(T.Translate(c, "jwt.identify.error", tokenName))
		}

		// Retrieving user from db
		if err := tx.Find(model, claims["jti"].(string)); err != nil {
			return nil, fmt.Errorf(T.Translate(c, "jwt.model.not.found", err))
		}
		return model, nil
	}
	return nil, fmt.Errorf(T.Translate(c, "jwt.validate.error", claims))
}
