package main

import (
	"context"

	cfg "github.com/conductorone/baton-panther/pkg/config"
	"github.com/conductorone/baton-panther/pkg/connector"
	"github.com/conductorone/baton-sdk/pkg/cli"
	"github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/connectorrunner"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

var version = "dev"

func main() {
	ctx := context.Background()

	config.RunConnector(
		ctx,
		"baton-panther",
		version,
		cfg.Config,
		newConnector,
		connectorrunner.WithDefaultCapabilitiesConnectorBuilderV2(&connector.Panther{}),
	)
}

func newConnector(ctx context.Context, c *cfg.Panther, opts *cli.ConnectorOpts) (connectorbuilder.ConnectorBuilderV2, []connectorbuilder.Opt, error) {
	l := ctxzap.Extract(ctx)

	cb, err := connector.New(ctx, c.Token, c.Url, opts.WillSyncResourceType("role"))
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, nil, err
	}

	return cb, nil, nil
}
