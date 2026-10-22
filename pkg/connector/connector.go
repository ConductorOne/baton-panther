package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-panther/pkg/panther"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

var (
	resourceTypeUser = &v2.ResourceType{
		Id:          "user",
		DisplayName: "User",
		Traits: []v2.ResourceType_Trait{
			v2.ResourceType_TRAIT_USER,
		},
	}
	resourceTypeRole = &v2.ResourceType{
		Id:          "role",
		DisplayName: "Role",
		Traits: []v2.ResourceType_Trait{
			v2.ResourceType_TRAIT_ROLE,
		},
		Annotations: annotationsForRoleResourceType(),
	}
)

type Panther struct {
	client *panther.Client
	// skipRoleGrants is true when the role resource type is not synced.
	// The zero value matches the default of syncing everything.
	skipRoleGrants bool
}

func (pn *Panther) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncerV2 {
	return []connectorbuilder.ResourceSyncerV2{
		userBuilder(pn.client, !pn.skipRoleGrants),
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

// New returns the Panther connector. User grants are only emitted when syncRoles is true.
func New(ctx context.Context, token, url string, syncRoles bool) (*Panther, error) {
	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	if err != nil {
		return nil, err
	}

	return &Panther{
		client:         panther.NewClient(httpClient, token, url),
		skipRoleGrants: !syncRoles,
	}, nil
}
