package zaplogger

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/term"

	"go-starter/pkg/config"
)

type Config struct {
	App    config.App
	Logger config.Logger
}

func New(c Config) (*zap.SugaredLogger, error) {
	loggerConfig := c.Logger.WithDefaults(c.App.Env)
	if err := loggerConfig.Validate(); err != nil {
		return nil, fmt.Errorf("validate logger config: %w", err)
	}

	level, err := zapcore.ParseLevel(string(loggerConfig.Level))
	if err != nil {
		return nil, fmt.Errorf("parse logger level %q: %w", loggerConfig.Level, err)
	}

	encoder, out, err := buildOutput(loggerConfig.Format)
	if err != nil {
		return nil, err
	}

	l := zap.New(
		zapcore.NewCore(
			encoder,
			zapcore.AddSync(out),
			zap.NewAtomicLevelAt(level),
		),
		zap.AddCaller(),
	).With(
		zap.String("app_name", c.App.Name),
		zap.String("app_version", c.App.Version),
	)

	return l.Sugar(), nil
}

func buildOutput(format config.LogFormat) (zapcore.Encoder, zapcore.WriteSyncer, error) {
	out := zapcore.Lock(os.Stderr)

	switch format {
	case config.LogFormatJSON:
		return zapcore.NewJSONEncoder(jsonEncoderConfig()), out, nil
	case config.LogFormatConsole:
		return zapcore.NewConsoleEncoder(consoleEncoderConfig(term.IsTerminal(int(os.Stderr.Fd())))), out, nil
	default:
		return nil, nil, fmt.Errorf("unsupported logger format %q", format)
	}
}

func jsonEncoderConfig() zapcore.EncoderConfig {
	cfg := baseEncoderConfig()
	cfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder

	return cfg
}

func consoleEncoderConfig(useColor bool) zapcore.EncoderConfig {
	cfg := baseEncoderConfig()
	cfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(applyColor(gruvboxGray, t.Format(time.DateTime), useColor))
	}
	cfg.EncodeCaller = func(c zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(applyColor(gruvboxBlue, c.TrimmedPath(), useColor))
	}
	cfg.EncodeLevel = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(applyColor(levelColor(level), level.CapitalString(), useColor))
	}

	return cfg
}

func baseEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		CallerKey:      "caller",
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

const (
	gruvboxGray   = 245
	gruvboxBlue   = 109
	gruvboxGreen  = 142
	gruvboxYellow = 214
	gruvboxOrange = 208
	gruvboxRed    = 167
	gruvboxReset  = "\x1b[0m"
)

func levelColor(level zapcore.Level) int {
	switch level {
	case zapcore.DebugLevel:
		return gruvboxBlue
	case zapcore.InfoLevel:
		return gruvboxGreen
	case zapcore.WarnLevel:
		return gruvboxYellow
	case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		return gruvboxRed
	default:
		return gruvboxOrange
	}
}

func applyColor(colorCode int, value string, enabled bool) string {
	if !enabled || value == "" {
		return value
	}

	return "\x1b[38;5;" + strconv.Itoa(colorCode) + "m" + value + gruvboxReset
}
