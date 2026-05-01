package env

import "fmt"

type AppEnv string

const (
	AppEnvLocal       AppEnv = "local"
	AppEnvDevelopment AppEnv = "development"
	AppEnvStaging     AppEnv = "staging"
	AppEnvProduction  AppEnv = "production"
)

func (e AppEnv) IsLocal() bool {
	return e == AppEnvLocal
}

func (e AppEnv) IsValid() bool {
	switch e {
	case AppEnvLocal, AppEnvDevelopment, AppEnvStaging, AppEnvProduction:
		return true
	default:
		return false
	}
}

func (e *AppEnv) UnmarshalText(text []byte) error {
	value := AppEnv(text)
	if !value.IsValid() {
		return fmt.Errorf("unsupported app env %q", text)
	}

	*e = value

	return nil
}
