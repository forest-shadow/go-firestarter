package main

import (
	"github.com/forest-shadow/go-firestarter/internal/app"
	"github.com/forest-shadow/go-firestarter/pkg/logger"
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

	defer l.Sync()

	l.Info("application started", logger.F("env", c.App.Env))
}
