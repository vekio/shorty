package cli

import (
	"context"
	"fmt"

	urfavecli "github.com/urfave/cli/v3"
	vekconfig "github.com/vekio/config"
)

func newSetupCommand(
	configFile *vekconfig.ConfigFile[Config],
	defaults Config,
) *urfavecli.Command {
	return &urfavecli.Command{
		Name:  "setup",
		Usage: "Initialize and validate the CLI configuration",
		Action: func(_ context.Context, command *urfavecli.Command) error {
			config, err := configFile.LoadOrCreate(defaults)
			if err != nil {
				return fmt.Errorf("setup CLI configuration: %w", err)
			}

			if _, err := fmt.Fprintf(
				command.Root().Writer,
				"Configuration: %s\n",
				configFile.Path(),
			); err != nil {
				return err
			}
			if config.APIKey == "" {
				_, err = fmt.Fprintln(
					command.Root().Writer,
					"Review server_url and set api_key before using the CLI.",
				)
				return err
			}
			_, err = fmt.Fprintln(command.Root().Writer, "Configuration is valid and ready.")
			return err
		},
	}
}
