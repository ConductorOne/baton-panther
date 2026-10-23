package connector

import (
	"context"
	"testing"

	cfg "github.com/conductorone/baton-panther/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/cli"
)

func TestNewSyncRoles(t *testing.T) {
	c := &cfg.Panther{Token: "t", Url: "http://localhost"}

	tests := []struct {
		name          string
		opts          *cli.ConnectorOpts
		wantSkipGrant bool
	}{
		{"nil opts syncs roles", nil, false},
		{"no filter syncs roles", &cli.ConnectorOpts{}, false},
		{"filter with role syncs roles", &cli.ConnectorOpts{SyncResourceTypeIDs: []string{"role", "user"}}, false},
		{"filter without role skips role grants", &cli.ConnectorOpts{SyncResourceTypeIDs: []string{"user"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb, _, err := New(context.Background(), c, tt.opts)
			if err != nil {
				t.Fatal(err)
			}
			if got := cb.(*Panther).skipRoleGrants; got != tt.wantSkipGrant {
				t.Fatalf("skipRoleGrants = %v, want %v", got, tt.wantSkipGrant)
			}
		})
	}
}
