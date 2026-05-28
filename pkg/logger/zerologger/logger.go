package zerologger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/forest-shadow/go-firestarter/pkg/config"
	"github.com/forest-shadow/go-firestarter/pkg/logger"
)

type Config struct {
	App    config.App
	Logger config.Logger
}

func New(c Config) (zerolog.Logger, error) {
	loggerConfig := c.Logger.WithDefaults(c.App.Env)
	if err := loggerConfig.Validate(); err != nil {
		return zerolog.Logger{}, fmt.Errorf("validate logger config: %w", err)
	}

	level, err := zerolog.ParseLevel(string(loggerConfig.Level))
	if err != nil {
		return zerolog.Logger{}, fmt.Errorf("parse logger level %q: %w", loggerConfig.Level, err)
	}

	var out io.Writer

	switch loggerConfig.Format {
	case logger.LogFormatJSON:
		out = os.Stderr

	case logger.LogFormatConsole:
		out = zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.DateTime,
		}

	default:
		return zerolog.Logger{}, fmt.Errorf("unsupported logger format %q", loggerConfig.Format)
	}

	l := zerolog.New(out).
		Hook(timestampHook{format: time.RFC3339Nano}).
		Level(level).
		With().
		Str("app_name", c.App.Name).
		Str("app_version", c.App.Version).
		Logger()

	return l, nil
}

type timestampHook struct {
	format string
}

func (h timestampHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	e.Str(zerolog.TimestampFieldName, time.Now().Format(h.format))
}
