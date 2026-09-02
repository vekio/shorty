package config

import (
	"fmt"
	"os"

	vekconfig "github.com/vekio/config"
)

const applicationName = "shorty"

func New[T vekconfig.Validatable](fileName string) (*vekconfig.ConfigFile[T], error) {
	return vekconfig.NewYAMLConfigFile[T](applicationName, fileName)
}

func SetPathFromEnv[T vekconfig.Validatable](file *vekconfig.ConfigFile[T], environmentName string) error {
	path, exists := os.LookupEnv(environmentName)
	if !exists || path == "" {
		return nil
	}
	if err := file.SetPath(path); err != nil {
		return fmt.Errorf("%s: %w", environmentName, err)
	}
	return nil
}
