package app

import (
	"go.uber.org/zap"

	"go-starter/pkg/zaplogger"
)

func NewLogger(c *Config) (*zap.SugaredLogger, error) {
	return zaplogger.New(zaplogger.Config{
		App:    c.App,
		Logger: c.Logger,
	})
}
