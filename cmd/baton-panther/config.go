package main

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-sdk/pkg/cli"
	"github.com/spf13/cobra"
)

// config defines the external configuration required for the connector to run.
type config struct {
	cli.BaseConfig `mapstructure:",squash"` // Puts the base config options in the same place as the connector options

	Token string `mapstructure:"token"`
	URL   string `mapstructure:"url"`
}

// validateConfig is run after the configuration is loaded, and should return an error if it isn't valid.
func validateConfig(ctx context.Context, cfg *config) error {
	if cfg.Token == "" {
		return fmt.Errorf("api token is missing")
	}

	if cfg.URL == "" {
		return fmt.Errorf("account api url is missing")
	}

	return nil
}

// cmdFlags sets the cmdFlags required for the connector.
func cmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("token", "", "API token used to authenticate to the Panther API. ($BATON_TOKEN)")
	cmd.PersistentFlags().String("url", "", "API url of your panther account. ($BATON_URL)")
}
