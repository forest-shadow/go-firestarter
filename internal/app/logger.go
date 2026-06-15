package app

import (
	"github.com/forest-shadow/go-firestarter/pkg/logger"
	"github.com/forest-shadow/go-firestarter/pkg/logger/zaplogger"
)

func NewLogger(c *Config) (logger.Logger, error) {
	return zaplogger.New(zaplogger.Config{
		App:    c.App,
		Logger: c.Logger,
	})
}
