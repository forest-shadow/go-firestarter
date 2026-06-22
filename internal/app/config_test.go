package app

import (
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"

	"github.com/forest-shadow/go-firestarter/pkg/config"
)

func TestAutomaticEnvOverridesHTTPServerDuration(t *testing.T) {
	t.Setenv("APP_HTTP_SERVER_READ_TIMEOUT", "7s")
	t.Setenv("APP_HTTP_SERVER_READ_HEADER_TIMEOUT", "3s")
	t.Setenv("APP_HTTP_SERVER_IDLE_TIMEOUT", "45s")

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetEnvPrefix("app")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadConfig(strings.NewReader(`app:
  name: starter
  version: 0.1.0
  env: local
logger:
  level: debug
  format: console
http_server:
  port: "8080"
  read_timeout: 20s
  read_header_timeout: 10s
  write_timeout: 20s
  idle_timeout: 60s
  shutdown_timeout: 25s
`)); err != nil {
		t.Fatalf("read config: %v", err)
	}

	var got Config
	if err := v.Unmarshal(&got, viper.DecodeHook(config.DecodeHook())); err != nil {
		t.Fatalf("expected successful unmarshal, got error: %v", err)
	}

	if got.HTTPServer.ReadTimeout != 7*time.Second {
		t.Fatalf("expected env read timeout override %s, got %s", 7*time.Second, got.HTTPServer.ReadTimeout)
	}

	if got.HTTPServer.ReadHeaderTimeout != 3*time.Second {
		t.Fatalf("expected env read header timeout override %s, got %s", 3*time.Second, got.HTTPServer.ReadHeaderTimeout)
	}

	if got.HTTPServer.IdleTimeout != 45*time.Second {
		t.Fatalf("expected env idle timeout override %s, got %s", 45*time.Second, got.HTTPServer.IdleTimeout)
	}
}
