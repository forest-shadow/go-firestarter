package app

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/forest-shadow/go-firestarter/pkg/config"
	e "github.com/forest-shadow/go-firestarter/pkg/env"
	"github.com/forest-shadow/go-firestarter/pkg/logger"
)

type Config struct {
	App        config.App        `mapstructure:"app"`
	Logger     logger.Config     `mapstructure:"logger"`
}

func GetConfig() (*Config, error) {
	appEnvName := os.Getenv("APP_ENV")
	if appEnvName == "" {
		log.Printf("APP_ENV is not set, defaulting to 'local'")
		appEnvName = string(e.AppEnvLocal)
	}

	v := viper.New()

	// config file setup
	v.SetConfigName("env." + appEnvName)
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	v.SetEnvPrefix("app")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg Config
	if err := v.Unmarshal(&cfg, viper.DecodeHook(config.DecodeHook())); err != nil {
		return nil, fmt.Errorf("decode config to struct: %w", err)
	}

	if err := cfg.App.Validate(); err != nil {
		return nil, fmt.Errorf("validate app config: %w", err)
	}

	cfg.Logger = cfg.Logger.WithDefaults(cfg.App.Env)
	if err := cfg.Logger.Validate(); err != nil {
		return nil, fmt.Errorf("validate logger config: %w", err)
	}

	return &cfg, nil
}
