package config

import (
	vekconfig "github.com/vekio/config"
)

const applicationName = "shorty"

func New[T vekconfig.Validatable](fileName string) (*vekconfig.ConfigFile[T], error) {
	return vekconfig.NewYAMLConfigFile[T](applicationName, fileName)
}
