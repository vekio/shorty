package config

import (
	"fmt"
	"os"

	vekconfig "github.com/vekio/config"
)

const applicationName = "shorty"

// New creates one Shorty configuration file in the operating system's
// configuration directory. Process-specific files share the Shorty namespace.
func New[T vekconfig.Validatable](fileName string) (*vekconfig.ConfigFile[T], error) {
	return vekconfig.NewYAMLConfigFile[T](applicationName, fileName)
}

// SetPathFromEnv overrides file's conventional path when environmentName
// contains a path. An unset or empty variable keeps the conventional path.
func SetPathFromEnv[T vekconfig.Validatable](
	file *vekconfig.ConfigFile[T],
	environmentName string,
) error {
	path, exists := os.LookupEnv(environmentName)
	if !exists || path == "" {
		return nil
	}
	if err := file.SetPath(path); err != nil {
		return fmt.Errorf("%s: %w", environmentName, err)
	}
	return nil
}
