package logger

import (
	"fmt"
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

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

type LogFormat string

const (
	LogFormatJSON    LogFormat = "json"
	LogFormatConsole LogFormat = "console"
)

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

type Field struct {
	Key   string
	Value any
}

func F(key string, value any) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	With(fields ...Field) Logger
	Sync() error
}
