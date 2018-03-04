package grifts

import (
	"github.com/emurmotol/coinssh/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
