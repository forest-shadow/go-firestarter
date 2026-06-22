package httpserver

import (
	"fmt"
	"time"
)

type Config struct {
	Port              string        `mapstructure:"port"`
	ReadTimeout       time.Duration `mapstructure:"read_timeout"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
	WriteTimeout      time.Duration `mapstructure:"write_timeout"`
	IdleTimeout       time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdown_timeout"`
}

func (c Config) WithDefaults() Config {
	if c.Port == "" {
		c.Port = "8080"
	}

	if c.ReadTimeout == 0 {
		c.ReadTimeout = 20 * time.Second
	}

	if c.ReadHeaderTimeout == 0 {
		c.ReadHeaderTimeout = 10 * time.Second
	}

	if c.WriteTimeout == 0 {
		c.WriteTimeout = 20 * time.Second
	}

	if c.IdleTimeout == 0 {
		c.IdleTimeout = 60 * time.Second
	}

	if c.ShutdownTimeout == 0 {
		c.ShutdownTimeout = 25 * time.Second
	}

	return c
}

func (c Config) Validate() error {
	if c.Port == "" {
		return fmt.Errorf("missing required config field: http_server.port")
	}

	if c.ReadTimeout <= 0 {
		return fmt.Errorf("invalid config field: http_server.read_timeout=%q", c.ReadTimeout)
	}

	if c.ReadHeaderTimeout <= 0 {
		return fmt.Errorf("invalid config field: http_server.read_header_timeout=%q", c.ReadHeaderTimeout)
	}

	if c.WriteTimeout <= 0 {
		return fmt.Errorf("invalid config field: http_server.write_timeout=%q", c.WriteTimeout)
	}

	if c.IdleTimeout <= 0 {
		return fmt.Errorf("invalid config field: http_server.idle_timeout=%q", c.IdleTimeout)
	}

	if c.ShutdownTimeout <= 0 {
		return fmt.Errorf("invalid config field: http_server.shutdown_timeout=%q", c.ShutdownTimeout)
	}

	return nil
}
