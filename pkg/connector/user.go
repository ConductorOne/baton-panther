package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-panther/pkg/panther"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	grant "github.com/conductorone/baton-sdk/pkg/types/grant"
	resource "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type userResourceType struct {
	resourceType *v2.ResourceType
	client       *panther.Client
	syncRoles    bool
}

func (u *userResourceType) ResourceType(_ context.Context) *v2.ResourceType {
	return u.resourceType
}

// Create a new connector resource for a Panther user.
func userResource(ctx context.Context, user *panther.User) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"first_name": user.GivenName,
		"last_name":  user.FamilyName,
		"login":      user.Email,
		"user_id":    user.ID,
		"role_id":    user.Role.ID,
	}

	userTraitOptions := []resource.UserTraitOption{
		resource.WithEmail(user.Email, true),
	}

	ret, err := resource.NewUserResource(
		fmt.Sprint(user.GivenName+" "+user.FamilyName),
		resourceTypeUser,
		user.ID,
		userTraitOptions,
		resource.WithResourceProfile(profile),
		resource.WithResourceStatus(v2.Status_RESOURCE_STATUS_ENABLED, ""),
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (u *userResourceType) List(ctx context.Context, _ *v2.ResourceId, _ resource.SyncOpAttrs) ([]*v2.Resource, *resource.SyncOpResults, error) {
	users, err := u.client.GetUsers(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("baton-panther: failed to list users: %w", err)
	}

	var rv []*v2.Resource
	for _, user := range users {
		userCopy := user
		ur, err := userResource(ctx, &userCopy)
		if err != nil {
			return nil, nil, err
		}
		rv = append(rv, ur)
	}

	return rv, nil, nil
}

func (u *userResourceType) Entitlements(_ context.Context, _ *v2.Resource, _ resource.SyncOpAttrs) ([]*v2.Entitlement, *resource.SyncOpResults, error) {
	return nil, nil, nil
}

// Grants emits the role membership grant for the user. Each Panther user has exactly one role.
func (u *userResourceType) Grants(ctx context.Context, res *v2.Resource, _ resource.SyncOpAttrs) ([]*v2.Grant, *resource.SyncOpResults, error) {
	if !u.syncRoles {
		return nil, nil, nil
	}

	roleID, ok := resource.GetProfileStringValue(res.GetProfile(), "role_id")
	if !ok || roleID == "" {
		return nil, nil, nil
	}

	rr, err := roleResource(ctx, panther.Role{ID: roleID})
	if err != nil {
		return nil, nil, err
	}

	return []*v2.Grant{grant.NewGrant(rr, membership, res.Id)}, nil, nil
}

func userBuilder(client *panther.Client, syncRoles bool) *userResourceType {
	return &userResourceType{
		resourceType: userResourceTypeFor(syncRoles),
		client:       client,
		syncRoles:    syncRoles,
	}
}
