package backend

import (
	"fmt"
	"log"
	"path"

	"gopkg.in/inconshreveable/log15.v2"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"sourcegraph.com/sourcegraph/sourcegraph/api/sourcegraph"
	"sourcegraph.com/sourcegraph/sourcegraph/pkg/htmlutil"
	"sourcegraph.com/sourcegraph/sourcegraph/pkg/store"
	"sourcegraph.com/sourcegraph/sourcegraph/services/backend/accesscontrol"
	"sourcegraph.com/sourcegraph/sourcegraph/services/svc"
	"sourcegraph.com/sourcegraph/srclib/graph"
	srcstore "sourcegraph.com/sourcegraph/srclib/store"
	"sourcegraph.com/sourcegraph/srclib/unit"
)

var Defs sourcegraph.DefsServer = &defs{}

type defs struct{}

var _ sourcegraph.DefsServer = (*defs)(nil)

func (s *defs) Get(ctx context.Context, op *sourcegraph.DefsGetOp) (*sourcegraph.Def, error) {
	defSpec := op.Def

	if err := accesscontrol.VerifyUserHasReadAccess(ctx, "Defs.Get", defSpec.Repo); err != nil {
		return nil, err
	}

	if !isAbsCommitID(defSpec.CommitID) {
		return nil, grpc.Errorf(codes.InvalidArgument, "absolute commit ID required (got %q)", defSpec.CommitID)
	}

	rawDef, err := s.get(ctx, defSpec)
	if err != nil {
		return nil, err
	}
	def := &sourcegraph.Def{Def: *rawDef}
	if op.Opt == nil {
		op.Opt = &sourcegraph.DefGetOptions{}
	}
	if op.Opt.Doc {
		def.DocHTML = htmlutil.EmptyForPB()
		if len(def.Docs) > 0 {
			def.DocHTML = htmlutil.SanitizeForPB(def.Docs[0].Data)
		}
	}
	if op.Opt.ComputeLineRange {
		startLine, endLine, err := computeLineRange(ctx, sourcegraph.TreeEntrySpec{
			RepoRev: sourcegraph.RepoRevSpec{
				RepoSpec: sourcegraph.RepoSpec{URI: defSpec.Repo},
				CommitID: defSpec.CommitID,
			},
			Path: def.File,
		}, def.DefStart, def.DefEnd)
		if err != nil {
			log15.Warn("Defs.Get: failed to compute line range.", "err", err, "repo", defSpec.Repo, "commitID", defSpec.CommitID, "file", def.File)
		}
		def.StartLine = startLine
		def.EndLine = endLine
	}
	populateDefFormatStrings(def)
	return def, nil
}

func computeLineRange(ctx context.Context, entrySpec sourcegraph.TreeEntrySpec, startByte, endByte uint32) (startLine, endLine uint32, err error) {
	entry, err := (&repoTree{}).Get(ctx, &sourcegraph.RepoTreeGetOp{
		Entry: entrySpec,
	})
	if err != nil {
		return
	}

	const max = 1024 * 1024 // 1 MB max size
	if len(entry.Contents) > max {
		err = fmt.Errorf("file exceeds max size (%d bytes)", max)
		return
	}

	line := uint32(1)
	for i, c := range entry.Contents {
		if uint32(i) == startByte {
			startLine = line
		}
		if uint32(i) == endByte {
			endLine = line
			break
		}
		if c == '\n' {
			line++
		}
	}
	return
}

// get returns the def with the given def key (and no additional
// information, such as docs).
func (s *defs) get(ctx context.Context, def sourcegraph.DefSpec) (*graph.Def, error) {
	d, err := store.GraphFromContext(ctx).Defs(srcstore.ByDefKey(def.DefKey()))
	if err != nil {
		return nil, err
	}
	if len(d) == 0 {
		return nil, grpc.Errorf(codes.NotFound, "def %v not found", def)
	}
	return d[0], nil
}

func (s *defs) List(ctx context.Context, opt *sourcegraph.DefListOptions) (*sourcegraph.DefList, error) {
	if opt == nil {
		opt = &sourcegraph.DefListOptions{}
	}

	if err := accesscontrol.VerifyUserHasReadAccess(ctx, "Defs.List", ""); err != nil {
		return nil, err
	}

	// Eliminate repos that don't exist.
	origRepoRevs := opt.RepoRevs
	opt.RepoRevs = nil
	for _, repoRev := range origRepoRevs {
		repoURI, commitID := sourcegraph.ParseRepoAndCommitID(repoRev)

		// Dealias. This call also verifies that the repo is visible to the current user.
		rA, err := svc.Repos(ctx).Get(ctx, &sourcegraph.RepoSpec{URI: repoURI})
		if err != nil {
			log.Printf("Warning: dropping repo rev %q from defs list because repo or repo alias was not found: %s.", repoRev, err)
			continue
		}
		repoURI = rA.URI

		// Determine the commit ID to use, if it wasn't specified or
		// if it's a non-commit-ID revspec.
		if !isAbsCommitID(commitID) {
			return nil, grpc.Errorf(codes.InvalidArgument, "absolute commit ID required for repo %q to list defs (got %q)", repoURI, commitID)
		}

		// The repo exists and the commit ID is valid, so include it
		// in the query.
		opt.RepoRevs = append(opt.RepoRevs, repoRev)
	}
	if len(origRepoRevs) > 0 && len(opt.RepoRevs) == 0 {
		log.Printf("Warning: DefsService.List got a RepoRevs param %v but none of the specified repos exist. Returning empty defs list.", origRepoRevs)
		return &sourcegraph.DefList{}, nil
	}

	// TODO(merge-to-master): don't try to search ALL repos until we
	// have a global index. Add a CLI flag to switch this behavior.
	//
	// if len(opt.RepoRevs) == 0 && len(opt.DefKeys) == 0 {
	// 	log.Println("WARNING: Defs.List cancelled - def queries that are not scoped to specific repos are rejected temporarily until global index exists!")
	// 	return &sourcegraph.DefList{}, nil
	// }

	fs := defListOptionsFilters(opt)
	fs = append(fs, srcstore.DefsSortByKey{})
	defs0, err := store.GraphFromContext(ctx).Defs(fs...)
	if err != nil {
		return nil, err
	}

	var defs []*sourcegraph.Def
	for i, def0 := range defs0 {
		if i >= opt.Offset() && i < (opt.Offset()+opt.Limit()) {
			defs = append(defs, &sourcegraph.Def{Def: *def0})
		}
	}
	// End kludge
	total := len(defs0)

	if opt.Doc {
		for _, def := range defs {
			def.DocHTML = htmlutil.EmptyForPB()
			if len(def.Docs) > 0 {
				def.DocHTML = htmlutil.SanitizeForPB(def.Docs[0].Data)
			}
		}
	}

	for _, def := range defs {
		populateDefFormatStrings(def)
	}

	return &sourcegraph.DefList{
		Defs: defs,
		ListResponse: sourcegraph.ListResponse{
			Total: int32(total),
		},
	}, nil
}

func defListOptionsFilters(o *sourcegraph.DefListOptions) []srcstore.DefFilter {
	var fs []srcstore.DefFilter
	if o.DefKeys != nil {
		fs = append(fs, srcstore.DefFilterFunc(func(def *graph.Def) bool {
			for _, dk := range o.DefKeys {
				if (def.Repo == "" || def.Repo == dk.Repo) && (def.CommitID == "" || def.CommitID == dk.CommitID) &&
					(def.UnitType == "" || def.UnitType == dk.UnitType) && (def.Unit == "" || def.Unit == dk.Unit) &&
					def.Path == dk.Path {
					return true
				}
			}
			return false
		}))
	}
	if o.Name != "" {
		fs = append(fs, srcstore.DefFilterFunc(func(def *graph.Def) bool {
			return def.Name == o.Name
		}))
	}
	if o.ByteEnd != 0 {
		fs = append(fs, srcstore.DefFilterFunc(func(d *graph.Def) bool {
			return d.DefStart == o.ByteStart && d.DefEnd == o.ByteEnd
		}))
	}
	if o.Query != "" {
		fs = append(fs, srcstore.ByDefQuery(o.Query))
	}
	if len(o.RepoRevs) > 0 {
		vs := make([]srcstore.Version, len(o.RepoRevs))
		for i, repoRev := range o.RepoRevs {
			repo, commitID := sourcegraph.ParseRepoAndCommitID(repoRev)
			if len(commitID) != 40 {
				log.Printf("WARNING: In DefListOptions.DefFilters, o.RepoRevs[%d]==%q has no commit ID or a non-absolute commit ID. No defs will match it.", i, repoRev)
			}
			vs[i] = srcstore.Version{Repo: repo, CommitID: commitID}
		}
		fs = append(fs, srcstore.ByRepoCommitIDs(vs...))
	}
	if o.Unit != "" && o.UnitType != "" {
		fs = append(fs, srcstore.ByUnits(unit.ID2{Type: o.UnitType, Name: o.Unit}))
	}
	if (o.UnitType != "" && o.Unit == "") || (o.UnitType == "" && o.Unit != "") {
		log.Println("WARNING: DefListOptions.DefFilter: must specify either both or neither of --type and --name (to filter by source unit)")
	}

	var files []string
	for _, f := range o.Files {
		if f != "" {
			files = append(files, path.Clean(f))
		}
	}
	if len(files) > 0 {
		fs = append(fs, srcstore.ByFiles(true, files...))
	}

	if o.FilePathPrefix != "" {
		fs = append(fs, srcstore.ByFiles(false, path.Clean(o.FilePathPrefix)))
	}
	if len(o.Kinds) > 0 {
		fs = append(fs, srcstore.DefFilterFunc(func(def *graph.Def) bool {
			for _, kind := range o.Kinds {
				if def.Kind == kind {
					return true
				}
			}
			return false
		}))
	}
	if o.Exported {
		fs = append(fs, srcstore.DefFilterFunc(func(def *graph.Def) bool {
			return def.Exported
		}))
	}
	if o.Nonlocal {
		fs = append(fs, srcstore.DefFilterFunc(func(def *graph.Def) bool {
			return !def.Local
		}))
	}
	if !o.IncludeTest {
		fs = append(fs, srcstore.DefFilterFunc(func(def *graph.Def) bool {
			return !def.Test
		}))
	}
	switch o.Sort {
	case "key":
		fs = append(fs, srcstore.DefsSortByKey{})
	case "name":
		fs = append(fs, srcstore.DefsSortByName{})
	}
	return fs
}

func populateDefFormatStrings(def *sourcegraph.Def) {
	if _, present := graph.MakeDefFormatters[def.UnitType]; !present {
		return
	}
	f := def.Fmt()
	quals := func(fn func(graph.Qualification) string) graph.QualFormatStrings {
		return graph.QualFormatStrings{
			Unqualified:             fn(graph.Unqualified),
			ScopeQualified:          fn(graph.ScopeQualified),
			DepQualified:            fn(graph.DepQualified),
			RepositoryWideQualified: fn(graph.RepositoryWideQualified),
			LanguageWideQualified:   fn(graph.LanguageWideQualified),
		}
	}
	def.FmtStrings = &graph.DefFormatStrings{
		Name:                 quals(f.Name),
		Type:                 quals(f.Type),
		NameAndTypeSeparator: f.NameAndTypeSeparator(),
		Language:             f.Language(),
		DefKeyword:           f.DefKeyword(),
		Kind:                 f.Kind(),
	}
}
