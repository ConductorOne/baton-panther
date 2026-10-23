package connector

import (
	"context"
	"fmt"

	cfg "github.com/conductorone/baton-panther/pkg/config"
	"github.com/conductorone/baton-panther/pkg/panther"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/cli"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

// RoleResourceTypeID is the ID of the role resource type.
const RoleResourceTypeID = "role"

var (
	resourceTypeUser = &v2.ResourceType{
		Id:          "user",
		DisplayName: "User",
		Traits: []v2.ResourceType_Trait{
			v2.ResourceType_TRAIT_USER,
		},
	}
	resourceTypeRole = &v2.ResourceType{
		Id:          RoleResourceTypeID,
		DisplayName: "Role",
		Traits: []v2.ResourceType_Trait{
			v2.ResourceType_TRAIT_ROLE,
		},
		Annotations: annotationsForRoleResourceType(),
	}
)

type Panther struct {
	client *panther.Client
	// skipRoleResourceType is true when the role resource type is excluded from the sync filter.
	// Named for the skip condition so the zero value used for capabilities syncs everything.
	skipRoleResourceType bool
}

func (pn *Panther) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncerV2 {
	return []connectorbuilder.ResourceSyncerV2{
		userBuilder(pn.client, pn.skipRoleResourceType),
		roleBuilder(pn.client),
	}
}

// Metadata returns metadata about the connector.
func (pn *Panther) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "Panther",
		Description: "Connector that syncs users and roles from Panther to Baton.",
	}, nil
}

// Validate hits the Panther API to ensure that the URL & API key passed are valid.
func (pn *Panther) Validate(ctx context.Context) (annotations.Annotations, error) {
	_, err := pn.client.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("baton-panther: failed to authenticate. Error: %w", err)
	}

	return nil, nil
}

// New returns the Panther connector. User grants are only emitted when the role resource type is synced.
func New(ctx context.Context, cfg *cfg.Panther, opts *cli.ConnectorOpts) (connectorbuilder.ConnectorBuilderV2, []connectorbuilder.Opt, error) {
	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	if err != nil {
		return nil, nil, err
	}

	// nil opts means no filter, so nothing is skipped.
	skipRoleResourceType := opts != nil && !opts.WillSyncResourceType(RoleResourceTypeID)

	return &Panther{
		client:               panther.NewClient(httpClient, cfg.Token, cfg.Url),
		skipRoleResourceType: skipRoleResourceType,
	}, nil, nil
}
