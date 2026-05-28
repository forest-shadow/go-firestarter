package config

import (
	"strings"
	"testing"

	"github.com/spf13/viper"

	"github.com/forest-shadow/go-firestarter/pkg/env"
	"github.com/forest-shadow/go-firestarter/pkg/logger"
)

func TestLoggerWithDefaults(t *testing.T) {
	t.Parallel()

	cfg := Logger{}.WithDefaults(env.AppEnvLocal)

	if cfg.Level != logger.LogLevelDebug {
		t.Fatalf("expected default level %q, got %q", logger.LogLevelDebug, cfg.Level)
	}

	if cfg.Format != logger.LogFormatConsole {
		t.Fatalf("expected default format %q, got %q",logger.LogFormatConsole, cfg.Format)
	}
}

func TestLoggerWithDefaultsProductionFormat(t *testing.T) {
	t.Parallel()

	cfg := Logger{}.WithDefaults(env.AppEnvProduction)

	if cfg.Format != logger.LogFormatJSON {
		t.Fatalf("expected production default format %q, got %q", logger.LogFormatJSON, cfg.Format)
	}
}

func TestLoggerValidate(t *testing.T) {
	t.Parallel()

	valid := Logger{
		Level:  logger.LogLevelInfo,
		Format: logger.LogFormatJSON,
	}

	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid config, got error: %v", err)
	}

	invalidLevel := Logger{
		Level:  "trace",
		Format: logger.LogFormatJSON,
	}

	if err := invalidLevel.Validate(); err == nil {
		t.Fatal("expected invalid level error")
	}

	invalidFormat := Logger{
		Level:  logger.LogLevelInfo,
		Format: "pretty",
	}

	if err := invalidFormat.Validate(); err == nil {
		t.Fatal("expected invalid format error")
	}
}

func TestLogLevelUnmarshalText(t *testing.T) {
	t.Parallel()

	var level logger.LogLevel
	if err := level.UnmarshalText([]byte("info")); err != nil {
		t.Fatalf("expected valid log level, got error: %v", err)
	}

	if level != logger.LogLevelInfo {
		t.Fatalf("expected log level %q, got %q", logger.LogLevelInfo, level)
	}

	if err := level.UnmarshalText([]byte("trace")); err == nil {
		t.Fatal("expected invalid log level error")
	}
}

func TestLogFormatUnmarshalText(t *testing.T) {
	t.Parallel()

	var format logger.LogFormat
	if err := format.UnmarshalText([]byte("json")); err != nil {
		t.Fatalf("expected valid log format, got error: %v", err)
	}

	if format != logger.LogFormatJSON {
		t.Fatalf("expected log format %q, got %q", logger.LogFormatJSON, format)
	}

	if err := format.UnmarshalText([]byte("pretty")); err == nil {
		t.Fatal("expected invalid log format error")
	}
}

func TestViperUnmarshalLoggerValidation(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Logger Logger `mapstructure:"logger"`
	}

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		v := viper.New()
		v.SetConfigType("yaml")

		if err := v.ReadConfig(strings.NewReader("logger:\n  level: info\n  format: json\n")); err != nil {
			t.Fatalf("read config: %v", err)
		}

		var got cfg
		if err := v.Unmarshal(&got, viper.DecodeHook(DecodeHook())); err != nil {
			t.Fatalf("expected successful unmarshal, got error: %v", err)
		}
	})

	t.Run("invalid level", func(t *testing.T) {
		t.Parallel()

		v := viper.New()
		v.SetConfigType("yaml")

		if err := v.ReadConfig(strings.NewReader("logger:\n  level: trace\n  format: json\n")); err != nil {
			t.Fatalf("read config: %v", err)
		}

		var got cfg
		if err := v.Unmarshal(&got, viper.DecodeHook(DecodeHook())); err == nil {
			t.Fatal("expected unmarshal error for invalid logger level")
		}
	})

	t.Run("invalid format", func(t *testing.T) {
		t.Parallel()

		v := viper.New()
		v.SetConfigType("yaml")

		if err := v.ReadConfig(strings.NewReader("logger:\n  level: info\n  format: pretty\n")); err != nil {
			t.Fatalf("read config: %v", err)
		}

		var got cfg
		if err := v.Unmarshal(&got, viper.DecodeHook(DecodeHook())); err == nil {
			t.Fatal("expected unmarshal error for invalid logger format")
		}
	})
}
