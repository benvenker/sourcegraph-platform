package database

import (
	"context"
	"reflect"
	"testing"

	"github.com/sourcegraph/sourcegraph/internal/database/dbtest"
)

func TestNamespaces(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	t.Parallel()
	db := NewDB(dbtest.NewDB(t))
	ctx := context.Background()

	// Create user and organization to test lookups.
	user, err := db.Users().Create(ctx, NewUser{Username: "alice"})
	if err != nil {
		t.Fatal(err)
	}
	org, err := db.Orgs().Create(ctx, "Acme", nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("GetByID", func(t *testing.T) {
		t.Run("no ID", func(t *testing.T) {
			ns, err := db.Namespaces().GetByID(ctx, 0, 0)
			if ns != nil {
				t.Errorf("unexpected non-nil namespace: %v", ns)
			}
			if want := ErrNamespaceNoID; err != want {
				t.Errorf("unexpected error: have=%v want=%v", err, want)
			}
		})

		t.Run("multiple IDs", func(t *testing.T) {
			ns, err := db.Namespaces().GetByID(ctx, 123, 456)
			if ns != nil {
				t.Errorf("unexpected non-nil namespace: %v", ns)
			}
			if want := ErrNamespaceMultipleIDs; err != want {
				t.Errorf("unexpected error: have=%v want=%v", err, want)
			}
		})

		t.Run("user not found", func(t *testing.T) {
			ns, err := db.Namespaces().GetByID(ctx, user.ID+1, 0)
			if ns != nil {
				t.Errorf("unexpected non-nil namespace: %v", ns)
			}
			if want := ErrNamespaceNotFound; err != want {
				t.Errorf("unexpected error: have=%v want=%v", err, want)
			}
		})

		t.Run("organization not found", func(t *testing.T) {
			ns, err := db.Namespaces().GetByID(ctx, 0, org.ID+1)
			if ns != nil {
				t.Errorf("unexpected non-nil namespace: %v", ns)
			}
			if want := ErrNamespaceNotFound; err != want {
				t.Errorf("unexpected error: have=%v want=%v", err, want)
			}
		})

		t.Run("user", func(t *testing.T) {
			ns, err := db.Namespaces().GetByID(ctx, 0, user.ID)
			if err != nil {
				t.Errorf("unexpected non-nil error: %v", err)
			}
			if want := (&Namespace{Name: "alice", User: user.ID}); !reflect.DeepEqual(ns, want) {
				t.Errorf("unexpected namespace: have=%v want=%v", ns, want)
			}
		})

		t.Run("organization", func(t *testing.T) {
			ns, err := db.Namespaces().GetByID(ctx, org.ID, 0)
			if err != nil {
				t.Errorf("unexpected non-nil error: %v", err)
			}
			if want := (&Namespace{Name: "Acme", Organization: org.ID}); !reflect.DeepEqual(ns, want) {
				t.Errorf("unexpected namespace: have=%v want=%v", ns, want)
			}
		})
	})

	t.Run("GetByName", func(t *testing.T) {
		t.Run("user", func(t *testing.T) {
			ns, err := db.Namespaces().GetByName(ctx, "Alice")
			if err != nil {
				t.Fatal(err)
			}
			if want := (&Namespace{Name: "alice", User: user.ID}); !reflect.DeepEqual(ns, want) {
				t.Errorf("got %+v, want %+v", ns, want)
			}
		})
		t.Run("organization", func(t *testing.T) {
			ns, err := db.Namespaces().GetByName(ctx, "acme")
			if err != nil {
				t.Fatal(err)
			}
			if want := (&Namespace{Name: "Acme", Organization: org.ID}); !reflect.DeepEqual(ns, want) {
				t.Errorf("got %+v, want %+v", ns, want)
			}
		})
		t.Run("not found", func(t *testing.T) {
			if _, err := db.Namespaces().GetByName(ctx, "doesntexist"); err != ErrNamespaceNotFound {
				t.Fatal(err)
			}
		})
	})
}
