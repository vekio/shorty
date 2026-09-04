// Package cli defines the Shorty command-line client.
package cli

import (
	urfavecli "github.com/urfave/cli/v3"
	vekconfig "github.com/vekio/config"
	configurfave "github.com/vekio/config/urfave"
)

const (
	configApplicationName = "shorty"
	configFileName        = "config.yml"
)

// New creates the root Shorty command.
func New() (*urfavecli.Command, error) {
	configFile, err := vekconfig.NewYAMLConfigFile[Config](configApplicationName, configFileName)
	if err != nil {
		return nil, err
	}

	defaults := DefaultConfig()
	return newRootCommand(configFile, defaults, newClient), nil
}

func newRootCommand(
	configFile *vekconfig.ConfigFile[Config],
	defaults Config,
	newClient clientFactory,
) *urfavecli.Command {
	return &urfavecli.Command{
		Name:  "shorty",
		Usage: "Manage links on a Shorty server",
		Flags: []urfavecli.Flag{
			configurfave.NewConfigFlag(configFile),
		},
		Commands: []*urfavecli.Command{
			newSetupCommand(configFile, defaults),
			configurfave.NewConfigCommand(configFile, defaults),
			newLinksCommand(configFile, newClient),
		},
	}
}
