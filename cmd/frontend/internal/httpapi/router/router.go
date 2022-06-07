// Package router contains the URL router for the HTTP API.
package router

import (
	"github.com/gorilla/mux"

	"github.com/sourcegraph/sourcegraph/cmd/frontend/internal/routevar"
)

const (
	LSIFUpload = "lsif.upload"
	GraphQL    = "graphql"

	GetStarted = "get-started"

	SearchStream  = "search.stream"
	ComputeStream = "compute.stream"

	SrcCliVersion  = "src-cli.version"
	SrcCliDownload = "src-cli.download"

	Registry = "registry"

	RepoShield  = "repo.shield"
	RepoRefresh = "repo.refresh"
	Telemetry   = "telemetry"

	GitHubWebhooks          = "github.webhooks"
	GitLabWebhooks          = "gitlab.webhooks"
	BitbucketServerWebhooks = "bitbucketServer.webhooks"
	BitbucketCloudWebhooks  = "bitbucketCloud.webhooks"

	SettingsGetForSubject  = "internal.settings.get-for-subject"
	OrgsListUsers          = "internal.orgs.list-users"
	OrgsGetByName          = "internal.orgs.get-by-name"
	UserEmailsGetEmail     = "internal.user-emails.get-email"
	ExternalURL            = "internal.app-url"
	SendEmail              = "internal.send-email"
	Extension              = "internal.extension"
	GitExec                = "internal.git.exec"
	GitInfoRefs            = "internal.git.info-refs"
	GitResolveRevision     = "internal.git.resolve-revision"
	GitUploadPack          = "internal.git.upload-pack"
	PhabricatorRepoCreate  = "internal.phabricator.repo.create"
	ReposGetByName         = "internal.repos.get-by-name"
	ReposInventoryUncached = "internal.repos.inventory-uncached"
	ReposInventory         = "internal.repos.inventory"
	ReposList              = "internal.repos.list"
	ReposIndex             = "internal.repos.index"
	Configuration          = "internal.configuration"
	SearchConfiguration    = "internal.search-configuration"
	ExternalServiceConfigs = "internal.external-services.configs"
	ExternalServicesList   = "internal.external-services.list"
	StreamingSearch        = "internal.stream-search"
)

// New creates a new API router with route URL pattern definitions but
// no handlers attached to the routes.
func New(base *mux.Router) *mux.Router {
	if base == nil {
		panic("base == nil")
	}

	base.StrictSlash(true)

	addRegistryRoute(base)
	addGraphQLRoute(base)
	base.Path("/github-webhooks").Methods("POST").Name(GitHubWebhooks)
	base.Path("/gitlab-webhooks").Methods("POST").Name(GitLabWebhooks)
	base.Path("/bitbucket-server-webhooks").Methods("POST").Name(BitbucketServerWebhooks)
	base.Path("/bitbucket-cloud-webhooks").Methods("POST").Name(BitbucketCloudWebhooks)
	base.Path("/lsif/upload").Methods("POST").Name(LSIFUpload)
	base.Path("/search/stream").Methods("GET").Name(SearchStream)
	base.Path("/compute/stream").Methods("GET", "POST").Name(ComputeStream)
	base.Path("/src-cli/version").Methods("GET").Name(SrcCliVersion)
	base.Path("/src-cli/{rest:.*}").Methods("GET").Name(SrcCliDownload)
	base.Path("/get-started").Methods("GET").Name(GetStarted)

	// repo contains routes that are NOT specific to a revision. In these routes, the URL may not contain a revspec after the repo (that is, no "github.com/foo/bar@myrevspec").
	repoPath := `/repos/` + routevar.Repo

	// Additional paths added will be treated as a repo. To add a new path that should not be treated as a repo
	// add above repo paths.
	repo := base.PathPrefix(repoPath + "/" + routevar.RepoPathDelim + "/").Subrouter()
	repo.Path("/shield").Methods("GET").Name(RepoShield)
	repo.Path("/refresh").Methods("POST").Name(RepoRefresh)

	return base
}

// NewInternal creates a new API router for internal endpoints.
func NewInternal(base *mux.Router) *mux.Router {
	if base == nil {
		panic("base == nil")
	}

	base.StrictSlash(true)
	// Internal API endpoints should only be served on the internal Handler
	base.Path("/settings/get-for-subject").Methods("POST").Name(SettingsGetForSubject)
	base.Path("/orgs/list-users").Methods("POST").Name(OrgsListUsers)
	base.Path("/orgs/get-by-name").Methods("POST").Name(OrgsGetByName)
	base.Path("/user-emails/get-email").Methods("POST").Name(UserEmailsGetEmail)
	base.Path("/app-url").Methods("POST").Name(ExternalURL)
	base.Path("/send-email").Methods("POST").Name(SendEmail)
	base.Path("/extension").Methods("POST").Name(Extension)
	base.Path("/git/{RepoID:[0-9]+}/exec").Methods("POST").Name(GitExec)
	base.Path("/git/{RepoName:.*}/info/refs").Methods("GET").Name(GitInfoRefs)
	base.Path("/git/{RepoName:.*}/resolve-revision/{Spec}").Methods("GET").Name(GitResolveRevision)
	base.Path("/git/{RepoName:.*}/git-upload-pack").Methods("GET", "POST").Name(GitUploadPack)
	base.Path("/phabricator/repo-create").Methods("POST").Name(PhabricatorRepoCreate)
	base.Path("/external-services/configs").Methods("POST").Name(ExternalServiceConfigs)
	base.Path("/external-services/list").Methods("POST").Name(ExternalServicesList)
	base.Path("/repos/inventory-uncached").Methods("POST").Name(ReposInventoryUncached)
	base.Path("/repos/inventory").Methods("POST").Name(ReposInventory)
	base.Path("/repos/list").Methods("POST").Name(ReposList)
	base.Path("/repos/index").Methods("POST").Name(ReposIndex)
	base.Path("/repos/{RepoName:.*}").Methods("POST").Name(ReposGetByName)
	base.Path("/configuration").Methods("POST").Name(Configuration)
	base.Path("/search/configuration").Methods("GET", "POST").Name(SearchConfiguration)
	base.Path("/telemetry").Methods("POST").Name(Telemetry)
	base.Path("/lsif/upload").Methods("POST").Name(LSIFUpload)
	base.Path("/search/stream").Methods("GET").Name(StreamingSearch)
	base.Path("/compute/stream").Methods("GET", "POST").Name(ComputeStream)
	addRegistryRoute(base)
	addGraphQLRoute(base)

	return base
}

func addRegistryRoute(m *mux.Router) {
	m.PathPrefix("/registry").Methods("GET").Name(Registry)
}

func addGraphQLRoute(m *mux.Router) {
	m.Path("/graphql").Methods("POST").Name(GraphQL)
}
