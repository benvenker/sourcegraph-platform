package lsifstore

import (
	"context"
	"database/sql"
	"runtime"

	"github.com/hashicorp/go-multierror"
	"github.com/keegancsmith/sqlf"
	"github.com/pkg/errors"
	"github.com/sourcegraph/sourcegraph/internal/db/basestore"
	"github.com/sourcegraph/sourcegraph/internal/db/batch"
	"github.com/sourcegraph/sourcegraph/internal/db/dbutil"
)

var ErrNoMetadata = errors.New("no rows in meta table")

type store struct {
	*basestore.Store
	serializer *serializer
}

var _ Store = &store{}

func NewStore(db dbutil.DB) Store {
	return &store{
		Store:      basestore.NewWithHandle(basestore.NewHandleWithDB(db, sql.TxOptions{})),
		serializer: newSerializer(),
	}
}

func (s *store) Transact(ctx context.Context) (Store, error) {
	tx, err := s.Store.Transact(ctx)
	if err != nil {
		return nil, err
	}

	return &store{
		Store:      tx,
		serializer: s.serializer,
	}, nil
}

func (s *store) Done(err error) error {
	return s.Store.Done(err)
}

func (s *store) ReadMeta(ctx context.Context, bundleID int) (MetaData, error) {
	numResultChunks, ok, err := basestore.ScanFirstInt(s.Store.Query(
		ctx,
		sqlf.Sprintf(
			`SELECT num_result_chunks FROM lsif_data_metadata WHERE dump_id = %s`,
			bundleID,
		),
	))
	if err != nil {
		return MetaData{}, err
	}
	if !ok {
		return MetaData{}, ErrNoMetadata
	}

	return MetaData{NumResultChunks: numResultChunks}, nil
}

func (s *store) PathsWithPrefix(ctx context.Context, bundleID int, prefix string) ([]string, error) {
	paths, err := basestore.ScanStrings(s.Store.Query(
		ctx,
		sqlf.Sprintf(
			`SELECT path FROM lsif_data_documents WHERE dump_id = %s AND path LIKE %s ORDER BY path`,
			bundleID,
			prefix+"%",
		),
	))
	if err != nil {
		return nil, err
	}

	return paths, nil
}

func (s *store) ReadDocument(ctx context.Context, bundleID int, path string) (DocumentData, bool, error) {
	data, ok, err := basestore.ScanFirstString(s.Store.Query(
		ctx,
		sqlf.Sprintf(
			`SELECT data FROM lsif_data_documents WHERE dump_id = %s AND path = %s LIMIT 1`,
			bundleID,
			path,
		),
	))
	if err != nil || !ok {
		return DocumentData{}, false, err
	}

	documentData, err := s.serializer.UnmarshalDocumentData([]byte(data))
	if err != nil {
		return DocumentData{}, false, err
	}

	return documentData, true, nil
}

func (s *store) ReadResultChunk(ctx context.Context, bundleID int, id int) (ResultChunkData, bool, error) {
	data, ok, err := basestore.ScanFirstString(s.Store.Query(
		ctx,
		sqlf.Sprintf(
			`SELECT data FROM lsif_data_result_chunks WHERE dump_id = %s AND idx = %s LIMIT 1`,
			bundleID,
			id,
		),
	))
	if err != nil || !ok {
		return ResultChunkData{}, false, err
	}

	resultChunkData, err := s.serializer.UnmarshalResultChunkData([]byte(data))
	if err != nil {
		return ResultChunkData{}, false, err
	}

	return resultChunkData, true, nil
}

func (s *store) ReadDefinitions(ctx context.Context, bundleID int, scheme, identifier string, skip, take int) ([]Location, int, error) {
	return s.readDefinitionReferences(ctx, bundleID, "lsif_data_definitions", scheme, identifier, skip, take)
}

func (s *store) ReadReferences(ctx context.Context, bundleID int, scheme, identifier string, skip, take int) ([]Location, int, error) {
	return s.readDefinitionReferences(ctx, bundleID, "lsif_data_references", scheme, identifier, skip, take)
}

func (s *store) readDefinitionReferences(ctx context.Context, bundleID int, tableName, scheme, identifier string, skip, take int) ([]Location, int, error) {
	data, ok, err := basestore.ScanFirstString(s.Store.Query(
		ctx,
		sqlf.Sprintf(
			`SELECT data FROM "`+tableName+`" WHERE dump_id = %s AND scheme = %s AND identifier = %s LIMIT 1`,
			bundleID,
			scheme,
			identifier,
		),
	))
	if err != nil || !ok {
		return nil, 0, err
	}

	locations, err := s.serializer.UnmarshalLocations([]byte(data))
	if err != nil {
		return nil, 0, err
	}

	if skip == 0 && take == 0 {
		// Pagination is disabled, return full result set
		return locations, len(locations), nil
	}

	lo := skip
	if lo >= len(locations) {
		// Skip lands past result set, return nothing
		return nil, len(locations), nil
	}

	hi := skip + take
	if hi >= len(locations) {
		hi = len(locations)
	}

	return locations[lo:hi], len(locations), nil
}

func (s *store) WriteMeta(ctx context.Context, bundleID int, meta MetaData) (err error) {
	inserter := batch.NewBatchInserter(ctx, s.Handle().DB(), "lsif_data_metadata", "dump_id", "num_result_chunks")

	defer func() {
		if flushErr := inserter.Flush(ctx); flushErr != nil {
			err = multierror.Append(err, errors.Wrap(flushErr, "inserter.Flush"))
		}
	}()

	return inserter.Insert(ctx, bundleID, meta.NumResultChunks)
}

func (s *store) WriteDocuments(ctx context.Context, bundleID int, documents chan KeyedDocumentData) error {
	inserter := func(inserter *batch.BatchInserter) error {
		for v := range documents {
			data, err := s.serializer.MarshalDocumentData(v.Document)
			if err != nil {
				return err
			}

			if err := inserter.Insert(ctx, bundleID, v.Path, data); err != nil {
				return err
			}
		}

		return nil
	}

	return withBatchInserter(ctx, s.Handle().DB(), "lsif_data_documents", []string{"dump_id", "path", "data"}, inserter)
}

func (s *store) WriteResultChunks(ctx context.Context, bundleID int, resultChunks chan IndexedResultChunkData) error {
	inserter := func(inserter *batch.BatchInserter) error {
		for v := range resultChunks {
			data, err := s.serializer.MarshalResultChunkData(v.ResultChunk)
			if err != nil {
				return err
			}

			if err := inserter.Insert(ctx, bundleID, v.Index, data); err != nil {
				return err
			}
		}

		return nil
	}

	return withBatchInserter(ctx, s.Handle().DB(), "lsif_data_result_chunks", []string{"dump_id", "idx", "data"}, inserter)
}

func (s *store) WriteDefinitions(ctx context.Context, bundleID int, monikerLocations chan MonikerLocations) error {
	return s.writeDefinitionReferences(ctx, bundleID, "lsif_data_definitions", monikerLocations)
}

func (s *store) WriteReferences(ctx context.Context, bundleID int, monikerLocations chan MonikerLocations) error {
	return s.writeDefinitionReferences(ctx, bundleID, "lsif_data_references", monikerLocations)
}

func (s *store) writeDefinitionReferences(ctx context.Context, bundleID int, tableName string, monikerLocations chan MonikerLocations) error {
	inserter := func(inserter *batch.BatchInserter) error {
		for v := range monikerLocations {
			data, err := s.serializer.MarshalLocations(v.Locations)
			if err != nil {
				return err
			}

			if err := inserter.Insert(ctx, bundleID, v.Scheme, v.Identifier, data); err != nil {
				return err
			}
		}

		return nil
	}

	return withBatchInserter(ctx, s.Handle().DB(), tableName, []string{"dump_id", "scheme", "identifier", "data"}, inserter)
}

var numWriterRoutines = runtime.GOMAXPROCS(0)

func withBatchInserter(ctx context.Context, db dbutil.DB, tableName string, columns []string, f func(inserter *batch.BatchInserter) error) error {
	return invokeN(numWriterRoutines, func() (err error) {
		inserter := batch.NewBatchInserter(ctx, db, tableName, columns...)

		defer func() {
			if flushErr := inserter.Flush(ctx); flushErr != nil {
				err = multierror.Append(err, errors.Wrap(flushErr, "inserter.Flush"))
			}
		}()

		return f(inserter)
	})
}
