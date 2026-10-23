package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-panther/pkg/panther"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	resource "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type roleResourceType struct {
	resourceType *v2.ResourceType
	client       *panther.Client
}

const membership = "member"

func (r *roleResourceType) ResourceType(_ context.Context) *v2.ResourceType {
	return r.resourceType
}

// Create a new connector resource for a Panther role.
func roleResource(ctx context.Context, role panther.Role) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"role_name": role.Name,
		"role_id":   role.ID,
	}

	roleTraitOptions := []resource.RoleTraitOption{}

	ret, err := resource.NewRoleResource(
		role.Name,
		resourceTypeRole,
		role.ID,
		roleTraitOptions,
		resource.WithResourceProfile(profile),
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *roleResourceType) List(ctx context.Context, _ *v2.ResourceId, _ resource.SyncOpAttrs) ([]*v2.Resource, *resource.SyncOpResults, error) {
	roles, err := r.client.GetRoles(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("baton-panther: failed to list roles: %w", err)
	}

	var rv []*v2.Resource
	for _, role := range roles {
		rr, err := roleResource(ctx, role)
		if err != nil {
			return nil, nil, err
		}
		rv = append(rv, rr)
	}

	return rv, nil, nil
}

func (r *roleResourceType) Entitlements(_ context.Context, res *v2.Resource, _ resource.SyncOpAttrs) ([]*v2.Entitlement, *resource.SyncOpResults, error) {
	assignmentOptions := []ent.EntitlementOption{
		ent.WithGrantableTo(resourceTypeUser),
		ent.WithDescription(fmt.Sprintf("%s Panther role", res.DisplayName)),
		ent.WithDisplayName(fmt.Sprintf("%s Role %s", res.DisplayName, membership)),
	}

	return []*v2.Entitlement{ent.NewAssignmentEntitlement(res, membership, assignmentOptions...)}, nil, nil
}

// Grants is a no-op. Role membership grants are emitted by the user resource type.
func (r *roleResourceType) Grants(_ context.Context, _ *v2.Resource, _ resource.SyncOpAttrs) ([]*v2.Grant, *resource.SyncOpResults, error) {
	return nil, nil, nil
}

func roleBuilder(client *panther.Client) *roleResourceType {
	return &roleResourceType{
		resourceType: resourceTypeRole,
		client:       client,
	}
}
