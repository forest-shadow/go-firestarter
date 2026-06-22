package httpserver

import (
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"

	"github.com/forest-shadow/go-firestarter/pkg/config"
)

func TestConfigWithDefaults(t *testing.T) {
	t.Parallel()

	cfg := Config{}.WithDefaults()

	if cfg.Port != "8080" {
		t.Fatalf("expected default port %q, got %q", "8080", cfg.Port)
	}

	if cfg.ReadTimeout != 20*time.Second {
		t.Fatalf("expected default read timeout %s, got %s", 20*time.Second, cfg.ReadTimeout)
	}

	if cfg.ReadHeaderTimeout != 10*time.Second {
		t.Fatalf("expected default read header timeout %s, got %s", 10*time.Second, cfg.ReadHeaderTimeout)
	}

	if cfg.WriteTimeout != 20*time.Second {
		t.Fatalf("expected default write timeout %s, got %s", 20*time.Second, cfg.WriteTimeout)
	}

	if cfg.IdleTimeout != 60*time.Second {
		t.Fatalf("expected default idle timeout %s, got %s", 60*time.Second, cfg.IdleTimeout)
	}

	if cfg.ShutdownTimeout != 25*time.Second {
		t.Fatalf("expected default shutdown timeout %s, got %s", 25*time.Second, cfg.ShutdownTimeout)
	}
}

func TestConfigValidate(t *testing.T) {
	t.Parallel()

	valid := Config{
		Port:              "8080",
		ReadTimeout:       time.Second,
		ReadHeaderTimeout: time.Second,
		WriteTimeout:      time.Second,
		IdleTimeout:       time.Second,
		ShutdownTimeout:   time.Second,
	}

	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid config, got error: %v", err)
	}

	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "missing port",
			config: Config{
				ReadTimeout:       time.Second,
				ReadHeaderTimeout: time.Second,
				WriteTimeout:      time.Second,
				IdleTimeout:       time.Second,
				ShutdownTimeout:   time.Second,
			},
		},
		{
			name: "invalid read timeout",
			config: Config{
				Port:              "8080",
				ReadHeaderTimeout: time.Second,
				WriteTimeout:      time.Second,
				IdleTimeout:       time.Second,
				ShutdownTimeout:   time.Second,
			},
		},
		{
			name: "invalid read header timeout",
			config: Config{
				Port:            "8080",
				ReadTimeout:     time.Second,
				WriteTimeout:    time.Second,
				IdleTimeout:     time.Second,
				ShutdownTimeout: time.Second,
			},
		},
		{
			name: "invalid write timeout",
			config: Config{
				Port:              "8080",
				ReadTimeout:       time.Second,
				ReadHeaderTimeout: time.Second,
				IdleTimeout:       time.Second,
				ShutdownTimeout:   time.Second,
			},
		},
		{
			name: "invalid idle timeout",
			config: Config{
				Port:              "8080",
				ReadTimeout:       time.Second,
				ReadHeaderTimeout: time.Second,
				WriteTimeout:      time.Second,
				ShutdownTimeout:   time.Second,
			},
		},
		{
			name: "invalid shutdown timeout",
			config: Config{
				Port:              "8080",
				ReadTimeout:       time.Second,
				ReadHeaderTimeout: time.Second,
				WriteTimeout:      time.Second,
				IdleTimeout:       time.Second,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := tt.config.Validate(); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestViperUnmarshalHTTPServerDuration(t *testing.T) {
	t.Parallel()

	type cfg struct {
		HTTPServer Config `mapstructure:"http_server"`
	}

	v := viper.New()
	v.SetConfigType("yaml")

	if err := v.ReadConfig(strings.NewReader(`http_server:
  port: "9090"
  read_timeout: 5s
  read_header_timeout: 6s
  write_timeout: 10s
  idle_timeout: 12s
  shutdown_timeout: 15s
`)); err != nil {
		t.Fatalf("read config: %v", err)
	}

	var got cfg
	if err := v.Unmarshal(&got, viper.DecodeHook(config.DecodeHook())); err != nil {
		t.Fatalf("expected successful unmarshal, got error: %v", err)
	}

	if got.HTTPServer.ReadTimeout != 5*time.Second {
		t.Fatalf("expected read timeout %s, got %s", 5*time.Second, got.HTTPServer.ReadTimeout)
	}

	if got.HTTPServer.ReadHeaderTimeout != 6*time.Second {
		t.Fatalf("expected read header timeout %s, got %s", 6*time.Second, got.HTTPServer.ReadHeaderTimeout)
	}

	if got.HTTPServer.WriteTimeout != 10*time.Second {
		t.Fatalf("expected write timeout %s, got %s", 10*time.Second, got.HTTPServer.WriteTimeout)
	}

	if got.HTTPServer.IdleTimeout != 12*time.Second {
		t.Fatalf("expected idle timeout %s, got %s", 12*time.Second, got.HTTPServer.IdleTimeout)
	}

	if got.HTTPServer.ShutdownTimeout != 15*time.Second {
		t.Fatalf("expected shutdown timeout %s, got %s", 15*time.Second, got.HTTPServer.ShutdownTimeout)
	}
}
