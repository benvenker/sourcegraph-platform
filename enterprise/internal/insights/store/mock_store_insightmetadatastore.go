// Code generated by go-mockgen 1.3.1; DO NOT EDIT.
//
// This file was generated by running `sg generate` (or `go-mockgen`) at the root of
// this repository. To add additional mocks to this or another package, add a new entry
// to the metadata.yaml file in the root of this repository.

package store

import (
	"context"
	"sync"

	types "github.com/sourcegraph/sourcegraph/enterprise/internal/insights/types"
)

// MockInsightMetadataStore is a mock implementation of the
// InsightMetadataStore interface (from the package
// github.com/sourcegraph/sourcegraph/enterprise/internal/insights/store)
// used for unit testing.
type MockInsightMetadataStore struct {
	// GetDirtyQueriesFunc is an instance of a mock function object
	// controlling the behavior of the method GetDirtyQueries.
	GetDirtyQueriesFunc *InsightMetadataStoreGetDirtyQueriesFunc
	// GetDirtyQueriesAggregatedFunc is an instance of a mock function
	// object controlling the behavior of the method
	// GetDirtyQueriesAggregated.
	GetDirtyQueriesAggregatedFunc *InsightMetadataStoreGetDirtyQueriesAggregatedFunc
	// GetMappedFunc is an instance of a mock function object controlling
	// the behavior of the method GetMapped.
	GetMappedFunc *InsightMetadataStoreGetMappedFunc
}

// NewMockInsightMetadataStore creates a new mock of the
// InsightMetadataStore interface. All methods return zero values for all
// results, unless overwritten.
func NewMockInsightMetadataStore() *MockInsightMetadataStore {
	return &MockInsightMetadataStore{
		GetDirtyQueriesFunc: &InsightMetadataStoreGetDirtyQueriesFunc{
			defaultHook: func(context.Context, *types.InsightSeries) (r0 []*types.DirtyQuery, r1 error) {
				return
			},
		},
		GetDirtyQueriesAggregatedFunc: &InsightMetadataStoreGetDirtyQueriesAggregatedFunc{
			defaultHook: func(context.Context, string) (r0 []*types.DirtyQueryAggregate, r1 error) {
				return
			},
		},
		GetMappedFunc: &InsightMetadataStoreGetMappedFunc{
			defaultHook: func(context.Context, InsightQueryArgs) (r0 []types.Insight, r1 error) {
				return
			},
		},
	}
}

// NewStrictMockInsightMetadataStore creates a new mock of the
// InsightMetadataStore interface. All methods panic on invocation, unless
// overwritten.
func NewStrictMockInsightMetadataStore() *MockInsightMetadataStore {
	return &MockInsightMetadataStore{
		GetDirtyQueriesFunc: &InsightMetadataStoreGetDirtyQueriesFunc{
			defaultHook: func(context.Context, *types.InsightSeries) ([]*types.DirtyQuery, error) {
				panic("unexpected invocation of MockInsightMetadataStore.GetDirtyQueries")
			},
		},
		GetDirtyQueriesAggregatedFunc: &InsightMetadataStoreGetDirtyQueriesAggregatedFunc{
			defaultHook: func(context.Context, string) ([]*types.DirtyQueryAggregate, error) {
				panic("unexpected invocation of MockInsightMetadataStore.GetDirtyQueriesAggregated")
			},
		},
		GetMappedFunc: &InsightMetadataStoreGetMappedFunc{
			defaultHook: func(context.Context, InsightQueryArgs) ([]types.Insight, error) {
				panic("unexpected invocation of MockInsightMetadataStore.GetMapped")
			},
		},
	}
}

// NewMockInsightMetadataStoreFrom creates a new mock of the
// MockInsightMetadataStore interface. All methods delegate to the given
// implementation, unless overwritten.
func NewMockInsightMetadataStoreFrom(i InsightMetadataStore) *MockInsightMetadataStore {
	return &MockInsightMetadataStore{
		GetDirtyQueriesFunc: &InsightMetadataStoreGetDirtyQueriesFunc{
			defaultHook: i.GetDirtyQueries,
		},
		GetDirtyQueriesAggregatedFunc: &InsightMetadataStoreGetDirtyQueriesAggregatedFunc{
			defaultHook: i.GetDirtyQueriesAggregated,
		},
		GetMappedFunc: &InsightMetadataStoreGetMappedFunc{
			defaultHook: i.GetMapped,
		},
	}
}

// InsightMetadataStoreGetDirtyQueriesFunc describes the behavior when the
// GetDirtyQueries method of the parent MockInsightMetadataStore instance is
// invoked.
type InsightMetadataStoreGetDirtyQueriesFunc struct {
	defaultHook func(context.Context, *types.InsightSeries) ([]*types.DirtyQuery, error)
	hooks       []func(context.Context, *types.InsightSeries) ([]*types.DirtyQuery, error)
	history     []InsightMetadataStoreGetDirtyQueriesFuncCall
	mutex       sync.Mutex
}

// GetDirtyQueries delegates to the next hook function in the queue and
// stores the parameter and result values of this invocation.
func (m *MockInsightMetadataStore) GetDirtyQueries(v0 context.Context, v1 *types.InsightSeries) ([]*types.DirtyQuery, error) {
	r0, r1 := m.GetDirtyQueriesFunc.nextHook()(v0, v1)
	m.GetDirtyQueriesFunc.appendCall(InsightMetadataStoreGetDirtyQueriesFuncCall{v0, v1, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the GetDirtyQueries
// method of the parent MockInsightMetadataStore instance is invoked and the
// hook queue is empty.
func (f *InsightMetadataStoreGetDirtyQueriesFunc) SetDefaultHook(hook func(context.Context, *types.InsightSeries) ([]*types.DirtyQuery, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// GetDirtyQueries method of the parent MockInsightMetadataStore instance
// invokes the hook at the front of the queue and discards it. After the
// queue is empty, the default hook function is invoked for any future
// action.
func (f *InsightMetadataStoreGetDirtyQueriesFunc) PushHook(hook func(context.Context, *types.InsightSeries) ([]*types.DirtyQuery, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *InsightMetadataStoreGetDirtyQueriesFunc) SetDefaultReturn(r0 []*types.DirtyQuery, r1 error) {
	f.SetDefaultHook(func(context.Context, *types.InsightSeries) ([]*types.DirtyQuery, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *InsightMetadataStoreGetDirtyQueriesFunc) PushReturn(r0 []*types.DirtyQuery, r1 error) {
	f.PushHook(func(context.Context, *types.InsightSeries) ([]*types.DirtyQuery, error) {
		return r0, r1
	})
}

func (f *InsightMetadataStoreGetDirtyQueriesFunc) nextHook() func(context.Context, *types.InsightSeries) ([]*types.DirtyQuery, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *InsightMetadataStoreGetDirtyQueriesFunc) appendCall(r0 InsightMetadataStoreGetDirtyQueriesFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of InsightMetadataStoreGetDirtyQueriesFuncCall
// objects describing the invocations of this function.
func (f *InsightMetadataStoreGetDirtyQueriesFunc) History() []InsightMetadataStoreGetDirtyQueriesFuncCall {
	f.mutex.Lock()
	history := make([]InsightMetadataStoreGetDirtyQueriesFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// InsightMetadataStoreGetDirtyQueriesFuncCall is an object that describes
// an invocation of method GetDirtyQueries on an instance of
// MockInsightMetadataStore.
type InsightMetadataStoreGetDirtyQueriesFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 *types.InsightSeries
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []*types.DirtyQuery
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c InsightMetadataStoreGetDirtyQueriesFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c InsightMetadataStoreGetDirtyQueriesFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// InsightMetadataStoreGetDirtyQueriesAggregatedFunc describes the behavior
// when the GetDirtyQueriesAggregated method of the parent
// MockInsightMetadataStore instance is invoked.
type InsightMetadataStoreGetDirtyQueriesAggregatedFunc struct {
	defaultHook func(context.Context, string) ([]*types.DirtyQueryAggregate, error)
	hooks       []func(context.Context, string) ([]*types.DirtyQueryAggregate, error)
	history     []InsightMetadataStoreGetDirtyQueriesAggregatedFuncCall
	mutex       sync.Mutex
}

// GetDirtyQueriesAggregated delegates to the next hook function in the
// queue and stores the parameter and result values of this invocation.
func (m *MockInsightMetadataStore) GetDirtyQueriesAggregated(v0 context.Context, v1 string) ([]*types.DirtyQueryAggregate, error) {
	r0, r1 := m.GetDirtyQueriesAggregatedFunc.nextHook()(v0, v1)
	m.GetDirtyQueriesAggregatedFunc.appendCall(InsightMetadataStoreGetDirtyQueriesAggregatedFuncCall{v0, v1, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the
// GetDirtyQueriesAggregated method of the parent MockInsightMetadataStore
// instance is invoked and the hook queue is empty.
func (f *InsightMetadataStoreGetDirtyQueriesAggregatedFunc) SetDefaultHook(hook func(context.Context, string) ([]*types.DirtyQueryAggregate, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// GetDirtyQueriesAggregated method of the parent MockInsightMetadataStore
// instance invokes the hook at the front of the queue and discards it.
// After the queue is empty, the default hook function is invoked for any
// future action.
func (f *InsightMetadataStoreGetDirtyQueriesAggregatedFunc) PushHook(hook func(context.Context, string) ([]*types.DirtyQueryAggregate, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *InsightMetadataStoreGetDirtyQueriesAggregatedFunc) SetDefaultReturn(r0 []*types.DirtyQueryAggregate, r1 error) {
	f.SetDefaultHook(func(context.Context, string) ([]*types.DirtyQueryAggregate, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *InsightMetadataStoreGetDirtyQueriesAggregatedFunc) PushReturn(r0 []*types.DirtyQueryAggregate, r1 error) {
	f.PushHook(func(context.Context, string) ([]*types.DirtyQueryAggregate, error) {
		return r0, r1
	})
}

func (f *InsightMetadataStoreGetDirtyQueriesAggregatedFunc) nextHook() func(context.Context, string) ([]*types.DirtyQueryAggregate, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *InsightMetadataStoreGetDirtyQueriesAggregatedFunc) appendCall(r0 InsightMetadataStoreGetDirtyQueriesAggregatedFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of
// InsightMetadataStoreGetDirtyQueriesAggregatedFuncCall objects describing
// the invocations of this function.
func (f *InsightMetadataStoreGetDirtyQueriesAggregatedFunc) History() []InsightMetadataStoreGetDirtyQueriesAggregatedFuncCall {
	f.mutex.Lock()
	history := make([]InsightMetadataStoreGetDirtyQueriesAggregatedFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// InsightMetadataStoreGetDirtyQueriesAggregatedFuncCall is an object that
// describes an invocation of method GetDirtyQueriesAggregated on an
// instance of MockInsightMetadataStore.
type InsightMetadataStoreGetDirtyQueriesAggregatedFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 string
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []*types.DirtyQueryAggregate
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c InsightMetadataStoreGetDirtyQueriesAggregatedFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c InsightMetadataStoreGetDirtyQueriesAggregatedFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// InsightMetadataStoreGetMappedFunc describes the behavior when the
// GetMapped method of the parent MockInsightMetadataStore instance is
// invoked.
type InsightMetadataStoreGetMappedFunc struct {
	defaultHook func(context.Context, InsightQueryArgs) ([]types.Insight, error)
	hooks       []func(context.Context, InsightQueryArgs) ([]types.Insight, error)
	history     []InsightMetadataStoreGetMappedFuncCall
	mutex       sync.Mutex
}

// GetMapped delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockInsightMetadataStore) GetMapped(v0 context.Context, v1 InsightQueryArgs) ([]types.Insight, error) {
	r0, r1 := m.GetMappedFunc.nextHook()(v0, v1)
	m.GetMappedFunc.appendCall(InsightMetadataStoreGetMappedFuncCall{v0, v1, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the GetMapped method of
// the parent MockInsightMetadataStore instance is invoked and the hook
// queue is empty.
func (f *InsightMetadataStoreGetMappedFunc) SetDefaultHook(hook func(context.Context, InsightQueryArgs) ([]types.Insight, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// GetMapped method of the parent MockInsightMetadataStore instance invokes
// the hook at the front of the queue and discards it. After the queue is
// empty, the default hook function is invoked for any future action.
func (f *InsightMetadataStoreGetMappedFunc) PushHook(hook func(context.Context, InsightQueryArgs) ([]types.Insight, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *InsightMetadataStoreGetMappedFunc) SetDefaultReturn(r0 []types.Insight, r1 error) {
	f.SetDefaultHook(func(context.Context, InsightQueryArgs) ([]types.Insight, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *InsightMetadataStoreGetMappedFunc) PushReturn(r0 []types.Insight, r1 error) {
	f.PushHook(func(context.Context, InsightQueryArgs) ([]types.Insight, error) {
		return r0, r1
	})
}

func (f *InsightMetadataStoreGetMappedFunc) nextHook() func(context.Context, InsightQueryArgs) ([]types.Insight, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *InsightMetadataStoreGetMappedFunc) appendCall(r0 InsightMetadataStoreGetMappedFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of InsightMetadataStoreGetMappedFuncCall
// objects describing the invocations of this function.
func (f *InsightMetadataStoreGetMappedFunc) History() []InsightMetadataStoreGetMappedFuncCall {
	f.mutex.Lock()
	history := make([]InsightMetadataStoreGetMappedFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// InsightMetadataStoreGetMappedFuncCall is an object that describes an
// invocation of method GetMapped on an instance of
// MockInsightMetadataStore.
type InsightMetadataStoreGetMappedFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 InsightQueryArgs
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []types.Insight
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c InsightMetadataStoreGetMappedFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c InsightMetadataStoreGetMappedFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}
