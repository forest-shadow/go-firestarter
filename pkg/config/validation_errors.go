package config

import "fmt"

func missingRequiredField(path string) error {
	return fmt.Errorf("missing required config field: %s", path)
}

func invalidField(path string, value any) error {
	return fmt.Errorf("invalid config field: %s=%q", path, value)
}
