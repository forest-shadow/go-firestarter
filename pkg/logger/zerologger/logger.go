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
	Logger logger.Config
}

type Logger struct {
	l zerolog.Logger
}

var _ logger.Logger = Logger{}

func New(c Config) (logger.Logger, error) {
	loggerConfig := c.Logger.WithDefaults(c.App.Env)
	if err := loggerConfig.Validate(); err != nil {
		return nil, fmt.Errorf("validate logger config: %w", err)
	}

	level, err := zerolog.ParseLevel(string(loggerConfig.Level))
	if err != nil {
		return nil, fmt.Errorf("parse logger level %q: %w", loggerConfig.Level, err)
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
		return nil, fmt.Errorf("unsupported logger format %q", loggerConfig.Format)
	}

	l := zerolog.New(out).
		Hook(timestampHook{format: time.RFC3339Nano}).
		Level(level).
		With().
		Str("app_name", c.App.Name).
		Str("app_version", c.App.Version).
		Logger()

	return Logger{l: l}, nil
}

func (l Logger) Debug(msg string, fields ...logger.Field) {
	event := l.l.Debug()
	addFields(event, fields)
	event.Msg(msg)
}

func (l Logger) Info(msg string, fields ...logger.Field) {
	event := l.l.Info()
	addFields(event, fields)
	event.Msg(msg)
}

func (l Logger) Warn(msg string, fields ...logger.Field) {
	event := l.l.Warn()
	addFields(event, fields)
	event.Msg(msg)
}

func (l Logger) Error(msg string, fields ...logger.Field) {
	event := l.l.Error()
	addFields(event, fields)
	event.Msg(msg)
}

func (l Logger) With(fields ...logger.Field) logger.Logger {
	ctx := l.l.With()
	for _, field := range fields {
		ctx = ctx.Interface(field.Key, field.Value)
	}

	return Logger{l: ctx.Logger()}
}

func (l Logger) Sync() error {
	return nil
}

func addFields(event *zerolog.Event, fields []logger.Field) {
	for _, field := range fields {
		event.Interface(field.Key, field.Value)
	}
}

type timestampHook struct {
	format string
}

func (h timestampHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	e.Str(zerolog.TimestampFieldName, time.Now().Format(h.format))
}
