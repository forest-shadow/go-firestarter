package logger

import (
	"fmt"

	"github.com/forest-shadow/go-firestarter/pkg/env"
)

type Config struct {
	Level  LogLevel  `mapstructure:"level"`
	Format LogFormat `mapstructure:"format"`
}

func (c Config) WithDefaults(appEnv env.AppEnv) Config {
	if c.Level == "" {
		c.Level = LogLevelDebug
	}

	if c.Format == "" {
		c.Format = DefaultLogFormat(appEnv)
	}

	return c
}

func (c Config) Validate() error {
	if !c.Level.IsValid() {
		return fmt.Errorf("unsupported logger level %q", c.Level)
	}

	if !c.Format.IsValid() {
		return fmt.Errorf("unsupported logger format %q", c.Format)
	}

	return nil
}

func DefaultLogFormat(appEnv env.AppEnv) LogFormat {
	switch appEnv {
	case env.AppEnvLocal, env.AppEnvDevelopment:
		return LogFormatConsole
	default:
		return LogFormatJSON
	}
}
