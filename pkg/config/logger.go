package config

import (
	"fmt"

	"github.com/go-viper/mapstructure/v2"

	"github.com/forest-shadow/go-firestarter/pkg/env"
)

type Logger struct {
	Level  LogLevel  `mapstructure:"level"`
	Format LogFormat `mapstructure:"format"`
}

type LogFormat string

const (
	LogFormatJSON    LogFormat = "json"
	LogFormatConsole LogFormat = "console"
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

func (c Logger) WithDefaults(appEnv env.AppEnv) Logger {
	if c.Level == "" {
		c.Level = LogLevelDebug
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

func (l LogLevel) IsValid() bool {
	switch l {
	case LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError:
		return true
	default:
		return false
	}
}

func (l *LogLevel) UnmarshalText(text []byte) error {
	value := LogLevel(text)
	if !value.IsValid() {
		return fmt.Errorf("unsupported logger level %q", text)
	}

	*l = value

	return nil
}

func (f LogFormat) IsValid() bool {
	switch f {
	case LogFormatJSON, LogFormatConsole:
		return true
	default:
		return false
	}
}

func (f *LogFormat) UnmarshalText(text []byte) error {
	value := LogFormat(text)
	if !value.IsValid() {
		return fmt.Errorf("unsupported logger format %q", text)
	}

	*f = value

	return nil
}

func defaultLogFormat(appEnv env.AppEnv) LogFormat {
	switch appEnv {
	case env.AppEnvLocal, env.AppEnvDevelopment:
		return LogFormatConsole
	default:
		return LogFormatJSON
	}
}

func DecodeHook() mapstructure.DecodeHookFunc {
	return mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToWeakSliceHookFunc(","),
		mapstructure.TextUnmarshallerHookFunc(),
	)
}
