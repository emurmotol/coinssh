package grifts

import (
	"github.com/emurmotol/coinssh/models"
	"github.com/manveru/faker"
	"github.com/markbates/grift/grift"
	"github.com/pkg/errors"
)

var _ = grift.Namespace("db", func() {

	grift.Add("seed:admin", func(c *grift.Context) error {

		admin := &models.User{
			Name:     "Administrator",
			Email:    "su@coinssh.com",
			Password: "su@coinssh.com",
		}

		return models.DB.Create(admin)
	})

	grift.Add("seed:users", func(c *grift.Context) error {
		fake, err := faker.New("en")

		if err != nil {
			return err
		}

		for i := 0; i <= 5; i++ {
			user := &models.User{
				Name:     fake.Name(),
				Email:    fake.Email(),
				Password: "secret123",
			}

			if err := models.DB.Create(user); err != nil {
				return err
			}
		}

		return nil
	})

	grift.Add("seed:accounts", func(c *grift.Context) error {
		fake, err := faker.New("en")

		if err != nil {
			return err
		}

		for i := 0; i <= 5; i++ {
			account := &models.Account{
				Name:     fake.Name(),
				Email:    fake.Email(),
				Password: "secret123",
			}

			if err := models.DB.Create(account); err != nil {
				return err
			}
		}

		return nil
	})

	grift.Add("seed", func(c *grift.Context) error {

		if err := models.DB.TruncateAll(); err != nil {
			return errors.WithStack(err)
		}

		if err := grift.Run("db:seed:admin", c); err != nil {
			return errors.WithStack(err)
		}

		if err := grift.Run("db:seed:users", c); err != nil {
			return errors.WithStack(err)
		}

		if err := grift.Run("db:seed:accounts", c); err != nil {
			return errors.WithStack(err)
		}

		return nil
	})

})
