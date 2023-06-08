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
		Annotations: annotationsForUserResourceType(),
	}
	resourceTypeRole = &v2.ResourceType{
		Id:          "role",
		DisplayName: "Role",
		Traits: []v2.ResourceType_Trait{
			v2.ResourceType_TRAIT_ROLE,
		},
	}
)

type Panther struct {
	client *panther.Client
}

func (pn *Panther) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		userBuilder(pn.client),
		roleBuilder(pn.client),
	}
}

// Metadata returns metadata about the connector.
func (pn *Panther) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "Panther",
	}, nil
}

// Validate hits the Panther API to ensure that the URL & API key passed are valid.
func (pn *Panther) Validate(ctx context.Context) (annotations.Annotations, error) {
	_, err := pn.client.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("panther-connector: failed to authenticate. Error: %w", err)
	}

	return nil, nil
}

// New returns the Panther connector.
func New(ctx context.Context, token, url string) (*Panther, error) {
	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	if err != nil {
		return nil, err
	}

	return &Panther{
		client: panther.NewClient(httpClient, token, url),
	}, nil
}
