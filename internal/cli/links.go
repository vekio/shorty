package cli

import (
	"context"
	"fmt"

	urfavecli "github.com/urfave/cli/v3"
	vekconfig "github.com/vekio/config"
	"github.com/vekio/shorty/pkg/shorty"
)

type linkCreator interface {
	CreateLink(context.Context, shorty.CreateLinkRequest) (shorty.CreateLinkResult, error)
}

type clientFactory func(Config) (linkCreator, error)

func newClient(config Config) (linkCreator, error) {
	return shorty.NewClient(config.ServerURL, config.APIKey)
}

func newLinksCommand(
	configFile *vekconfig.ConfigFile[Config],
	newClient clientFactory,
) *urfavecli.Command {
	return &urfavecli.Command{
		Name:  "links",
		Usage: "Manage short links",
		Commands: []*urfavecli.Command{
			newCreateLinkCommand(configFile, newClient),
		},
	}
}

func newCreateLinkCommand(
	configFile *vekconfig.ConfigFile[Config],
	newClient clientFactory,
) *urfavecli.Command {
	return &urfavecli.Command{
		Name:      "create",
		Usage:     "Create a short link",
		ArgsUsage: "<URL>",
		Arguments: []urfavecli.Argument{
			&urfavecli.StringArg{
				Name:      "url",
				UsageText: "destination URL",
				Config:    urfavecli.StringConfig{TrimSpace: true},
			},
		},
		Action: func(ctx context.Context, command *urfavecli.Command) error {
			originURL := command.StringArg("url")
			if originURL == "" {
				return fmt.Errorf("destination URL is required")
			}
			config, err := configFile.Load()
			if err != nil {
				return fmt.Errorf("load CLI configuration: %w", err)
			}
			client, err := newClient(config)
			if err != nil {
				return fmt.Errorf("configure Shorty client: %w", err)
			}
			result, err := client.CreateLink(ctx, shorty.CreateLinkRequest{OriginURL: originURL})
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(command.Root().Writer, result.Code)
			return err
		},
	}
}
