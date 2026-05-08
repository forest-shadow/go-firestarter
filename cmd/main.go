package main

import (
	"github.com/forest-shadow/go-firestarter/internal/app"
)

func main() {
	c, err := app.GetConfig()
	if err != nil {
		panic(err)
	}

	l, err := app.NewLogger(c)
	if err != nil {
		panic(err)
	}

	l.Infow("application started", "env", c.App.Env)
}
