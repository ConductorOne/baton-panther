package connector

import (
	"context"
	"testing"

	"github.com/conductorone/baton-panther/pkg/panther"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
)

func TestUserGrants(t *testing.T) {
	ctx := context.Background()
	ur, err := userResource(ctx, &panther.User{ID: "u1", Email: "a@example.com", Role: panther.Role{ID: "r1", Name: "Admin"}})
	if err != nil {
		t.Fatal(err)
	}

	grants, _, err := userBuilder(nil, false).Grants(ctx, ur, resource.SyncOpAttrs{})
	if err != nil {
		t.Fatal(err)
	}
	if len(grants) != 1 {
		t.Fatalf("expected 1 grant, got %d", len(grants))
	}
	role := grants[0].GetEntitlement().GetResource().GetId()
	if role.GetResourceType() != "role" || role.GetResource() != "r1" {
		t.Fatalf("unexpected role resource: %v", role)
	}
	if grants[0].GetPrincipal().GetId().GetResource() != "u1" {
		t.Fatalf("unexpected principal: %v", grants[0].GetPrincipal())
	}

	grants, _, err = userBuilder(nil, true).Grants(ctx, ur, resource.SyncOpAttrs{})
	if err != nil {
		t.Fatal(err)
	}
	if len(grants) != 0 {
		t.Fatalf("expected no grants when roles are not synced, got %d", len(grants))
	}
}

func TestUserResourceTypeAnnotations(t *testing.T) {
	ctx := context.Background()
	with := annotations.Annotations(userBuilder(nil, false).ResourceType(ctx).GetAnnotations())
	without := annotations.Annotations(userBuilder(nil, true).ResourceType(ctx).GetAnnotations())

	if !with.Contains(&v2.SkipEntitlements{}) || with.Contains(&v2.SkipEntitlementsAndGrants{}) {
		t.Fatal("with roles: expected SkipEntitlements only")
	}
	if !without.Contains(&v2.SkipEntitlementsAndGrants{}) {
		t.Fatal("without roles: expected SkipEntitlementsAndGrants")
	}
	if len(resourceTypeUser.GetAnnotations()) != 0 {
		t.Fatal("shared user resource type must not be mutated")
	}
}

func TestRoleGrantsNoop(t *testing.T) {
	grants, _, err := roleBuilder(nil).Grants(context.Background(), nil, resource.SyncOpAttrs{})
	if err != nil || len(grants) != 0 {
		t.Fatalf("expected no grants, got %v, %v", grants, err)
	}
}

func TestZeroValueConnector_DoesNotSkipGrants(t *testing.T) {
	syncers := (&Panther{}).ResourceSyncers(context.Background())

	var userRT *v2.ResourceType
	for _, rs := range syncers {
		if rt := rs.ResourceType(context.Background()); rt.GetId() == resourceTypeUser.GetId() {
			userRT = rt
		}
	}
	if userRT == nil {
		t.Fatal("user resource syncer not found")
	}

	annos := annotations.Annotations(userRT.GetAnnotations())
	if !annos.Contains(&v2.SkipEntitlements{}) || annos.Contains(&v2.SkipEntitlementsAndGrants{}) {
		t.Fatal("zero-value Panther must advertise user grants (SkipEntitlements only)")
	}
}
