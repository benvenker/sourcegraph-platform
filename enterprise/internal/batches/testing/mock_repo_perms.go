package testing

import (
	"context"
	"testing"
	"time"

	edb "github.com/sourcegraph/sourcegraph/enterprise/internal/database"
	"github.com/sourcegraph/sourcegraph/internal/api"
	"github.com/sourcegraph/sourcegraph/internal/authz"
	"github.com/sourcegraph/sourcegraph/internal/database"
)

// MockRepoPermissions mocks repository permissions to include
// repositories by IDs for the given user.
func MockRepoPermissions(t *testing.T, db database.DB, userID int32, repoIDs ...api.RepoID) {
	t.Helper()

	permsStore := edb.Perms(db, time.Now)

	userIDs := map[int32]struct{}{
		userID: {},
	}
	for _, id := range repoIDs {
		err := permsStore.SetRepoPermissions(context.Background(),
			&authz.RepoPermissions{
				RepoID:  int32(id),
				Perm:    authz.Read,
				UserIDs: userIDs,
			},
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	authz.SetProviders(false, nil)
	t.Cleanup(func() {
		authz.SetProviders(true, nil)
	})
}
