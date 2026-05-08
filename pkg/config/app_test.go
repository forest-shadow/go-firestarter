package config

import (
	"strings"
	"testing"

	"github.com/spf13/viper"

	"github.com/forest-shadow/go-firestarter/pkg/env"
)

func TestAppValidate(t *testing.T) {
	t.Parallel()

	valid := App{
		Name:    "starter",
		Version: "0.1.0",
		Env:     env.AppEnvLocal,
	}

	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid app config, got error: %v", err)
	}

	tests := []struct {
		name string
		app  App
	}{
		{
			name: "missing name",
			app: App{
				Version: "0.1.0",
				Env:     env.AppEnvLocal,
			},
		},
		{
			name: "missing version",
			app: App{
				Name: "starter",
				Env:  env.AppEnvLocal,
			},
		},
		{
			name: "invalid env",
			app: App{
				Name:    "starter",
				Version: "0.1.0",
				Env:     "qa",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := tt.app.Validate(); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestAppEnvUnmarshalText(t *testing.T) {
	t.Parallel()

	var appEnv env.AppEnv
	if err := appEnv.UnmarshalText([]byte("local")); err != nil {
		t.Fatalf("expected valid app env, got error: %v", err)
	}

	if appEnv != env.AppEnvLocal {
		t.Fatalf("expected app env %q, got %q", env.AppEnvLocal, appEnv)
	}

	if err := appEnv.UnmarshalText([]byte("qa")); err == nil {
		t.Fatal("expected invalid app env error")
	}
}

func TestViperUnmarshalAppEnvValidation(t *testing.T) {
	t.Parallel()

	type cfg struct {
		App App `mapstructure:"app"`
	}

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		v := viper.New()
		v.SetConfigType("yaml")

		if err := v.ReadConfig(strings.NewReader("app:\n  name: starter\n  version: 0.1.0\n  env: local\n")); err != nil {
			t.Fatalf("read config: %v", err)
		}

		var got cfg
		if err := v.Unmarshal(&got, viper.DecodeHook(DecodeHook())); err != nil {
			t.Fatalf("expected successful unmarshal, got error: %v", err)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		v := viper.New()
		v.SetConfigType("yaml")

		if err := v.ReadConfig(strings.NewReader("app:\n  name: starter\n  version: 0.1.0\n  env: qa\n")); err != nil {
			t.Fatalf("read config: %v", err)
		}

		var got cfg
		if err := v.Unmarshal(&got, viper.DecodeHook(DecodeHook())); err == nil {
			t.Fatal("expected unmarshal error for invalid app env")
		}
	})
}
