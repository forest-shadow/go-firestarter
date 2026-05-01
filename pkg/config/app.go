package config

import (
	"fmt"
	"strings"

	"go-starter/pkg/env"
)

type App struct {
	Name    string     `mapstructure:"name"		required:"true"`
	Version string     `mapstructure:"version"	required:"true"`
	Env     env.AppEnv `mapstructure:"env"		required:"true"`
}

func (a App) Validate() error {
	if strings.TrimSpace(a.Name) == "" {
		return fmt.Errorf("app name is required")
	}

	if strings.TrimSpace(a.Version) == "" {
		return fmt.Errorf("app version is required")
	}

	if !a.Env.IsValid() {
		return fmt.Errorf("unsupported app env %q", a.Env)
	}

	return nil
}
