package connector

import (
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"google.golang.org/protobuf/proto"
)

// annotationsForUserResourceType returns the user resource type annotations.
// Users have no entitlements. They only have grants when roles are synced.
func annotationsForUserResourceType(syncRoles bool) annotations.Annotations {
	annos := annotations.Annotations{}
	if syncRoles {
		annos.Update(&v2.SkipEntitlements{})
	} else {
		annos.Update(&v2.SkipEntitlementsAndGrants{})
	}
	return annos
}

// annotationsForRoleResourceType returns the role resource type annotations.
// Role grants are emitted from the user side.
func annotationsForRoleResourceType() annotations.Annotations {
	annos := annotations.Annotations{}
	annos.Update(&v2.SkipGrants{})
	return annos
}

// userResourceTypeFor returns a copy of resourceTypeUser with annotations that depend on syncRoles.
func userResourceTypeFor(syncRoles bool) *v2.ResourceType {
	rt := proto.Clone(resourceTypeUser).(*v2.ResourceType)
	rt.Annotations = annotationsForUserResourceType(syncRoles)
	return rt
}
