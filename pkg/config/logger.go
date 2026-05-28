package config

import (
	"fmt"

	"github.com/go-viper/mapstructure/v2"

	"github.com/forest-shadow/go-firestarter/pkg/env"
	"github.com/forest-shadow/go-firestarter/pkg/logger"
)

type Logger struct {
	Level  logger.LogLevel  `mapstructure:"level"`
	Format logger.LogFormat `mapstructure:"format"`
}

func (c Logger) WithDefaults(appEnv env.AppEnv) Logger {
	if c.Level == "" {
		c.Level = logger.LogLevelDebug
	}

	if c.Format == "" {
		c.Format = defaultLogFormat(appEnv)
	}

	return c
}

func (c Logger) Validate() error {
	if !c.Level.IsValid() {
		return fmt.Errorf("unsupported logger level %q", c.Level)
	}

	if !c.Format.IsValid() {
		return fmt.Errorf("unsupported logger format %q", c.Format)
	}

	return nil
}

func defaultLogFormat(appEnv env.AppEnv) logger.LogFormat {
	switch appEnv {
	case env.AppEnvLocal, env.AppEnvDevelopment:
		return logger.LogFormatConsole
	default:
		return logger.LogFormatJSON
	}
}

func DecodeHook() mapstructure.DecodeHookFunc {
	return mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToWeakSliceHookFunc(","),
		mapstructure.TextUnmarshallerHookFunc(),
	)
}
