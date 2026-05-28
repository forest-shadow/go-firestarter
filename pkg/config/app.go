package config

import (
	"strings"

	"github.com/forest-shadow/go-firestarter/pkg/env"
)

type App struct {
	Name    string     `mapstructure:"name"`
	Version string     `mapstructure:"version"`
	Env     env.AppEnv `mapstructure:"env"`
}

func (a App) Validate() error {
	if strings.TrimSpace(a.Name) == "" {
		return missingRequiredField("app.name")
	}

	if strings.TrimSpace(a.Version) == "" {
		return missingRequiredField("app.version")
	}

	if !a.Env.IsValid() {
		return invalidField("app.env", a.Env)
	}

	return nil
}
