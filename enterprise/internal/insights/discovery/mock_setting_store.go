// Code generated by go-mockgen 1.3.1; DO NOT EDIT.
//
// This file was generated by running `sg generate` (or `go-mockgen`) at the root of
// this repository. To add additional mocks to this or another package, add a new entry
// to the metadata.yaml file in the root of this repository.

package discovery

import (
	"context"
	"sync"

	api "github.com/sourcegraph/sourcegraph/internal/api"
	schema "github.com/sourcegraph/sourcegraph/schema"
)

// MockSettingStore is a mock implementation of the SettingStore interface
// (from the package
// github.com/sourcegraph/sourcegraph/enterprise/internal/insights/discovery)
// used for unit testing.
type MockSettingStore struct {
	// GetLastestSchemaSettingsFunc is an instance of a mock function object
	// controlling the behavior of the method GetLastestSchemaSettings.
	GetLastestSchemaSettingsFunc *SettingStoreGetLastestSchemaSettingsFunc
	// GetLatestFunc is an instance of a mock function object controlling
	// the behavior of the method GetLatest.
	GetLatestFunc *SettingStoreGetLatestFunc
}

// NewMockSettingStore creates a new mock of the SettingStore interface. All
// methods return zero values for all results, unless overwritten.
func NewMockSettingStore() *MockSettingStore {
	return &MockSettingStore{
		GetLastestSchemaSettingsFunc: &SettingStoreGetLastestSchemaSettingsFunc{
			defaultHook: func(context.Context, api.SettingsSubject) (r0 *schema.Settings, r1 error) {
				return
			},
		},
		GetLatestFunc: &SettingStoreGetLatestFunc{
			defaultHook: func(context.Context, api.SettingsSubject) (r0 *api.Settings, r1 error) {
				return
			},
		},
	}
}

// NewStrictMockSettingStore creates a new mock of the SettingStore
// interface. All methods panic on invocation, unless overwritten.
func NewStrictMockSettingStore() *MockSettingStore {
	return &MockSettingStore{
		GetLastestSchemaSettingsFunc: &SettingStoreGetLastestSchemaSettingsFunc{
			defaultHook: func(context.Context, api.SettingsSubject) (*schema.Settings, error) {
				panic("unexpected invocation of MockSettingStore.GetLastestSchemaSettings")
			},
		},
		GetLatestFunc: &SettingStoreGetLatestFunc{
			defaultHook: func(context.Context, api.SettingsSubject) (*api.Settings, error) {
				panic("unexpected invocation of MockSettingStore.GetLatest")
			},
		},
	}
}

// NewMockSettingStoreFrom creates a new mock of the MockSettingStore
// interface. All methods delegate to the given implementation, unless
// overwritten.
func NewMockSettingStoreFrom(i SettingStore) *MockSettingStore {
	return &MockSettingStore{
		GetLastestSchemaSettingsFunc: &SettingStoreGetLastestSchemaSettingsFunc{
			defaultHook: i.GetLastestSchemaSettings,
		},
		GetLatestFunc: &SettingStoreGetLatestFunc{
			defaultHook: i.GetLatest,
		},
	}
}

// SettingStoreGetLastestSchemaSettingsFunc describes the behavior when the
// GetLastestSchemaSettings method of the parent MockSettingStore instance
// is invoked.
type SettingStoreGetLastestSchemaSettingsFunc struct {
	defaultHook func(context.Context, api.SettingsSubject) (*schema.Settings, error)
	hooks       []func(context.Context, api.SettingsSubject) (*schema.Settings, error)
	history     []SettingStoreGetLastestSchemaSettingsFuncCall
	mutex       sync.Mutex
}

// GetLastestSchemaSettings delegates to the next hook function in the queue
// and stores the parameter and result values of this invocation.
func (m *MockSettingStore) GetLastestSchemaSettings(v0 context.Context, v1 api.SettingsSubject) (*schema.Settings, error) {
	r0, r1 := m.GetLastestSchemaSettingsFunc.nextHook()(v0, v1)
	m.GetLastestSchemaSettingsFunc.appendCall(SettingStoreGetLastestSchemaSettingsFuncCall{v0, v1, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the
// GetLastestSchemaSettings method of the parent MockSettingStore instance
// is invoked and the hook queue is empty.
func (f *SettingStoreGetLastestSchemaSettingsFunc) SetDefaultHook(hook func(context.Context, api.SettingsSubject) (*schema.Settings, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// GetLastestSchemaSettings method of the parent MockSettingStore instance
// invokes the hook at the front of the queue and discards it. After the
// queue is empty, the default hook function is invoked for any future
// action.
func (f *SettingStoreGetLastestSchemaSettingsFunc) PushHook(hook func(context.Context, api.SettingsSubject) (*schema.Settings, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *SettingStoreGetLastestSchemaSettingsFunc) SetDefaultReturn(r0 *schema.Settings, r1 error) {
	f.SetDefaultHook(func(context.Context, api.SettingsSubject) (*schema.Settings, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *SettingStoreGetLastestSchemaSettingsFunc) PushReturn(r0 *schema.Settings, r1 error) {
	f.PushHook(func(context.Context, api.SettingsSubject) (*schema.Settings, error) {
		return r0, r1
	})
}

func (f *SettingStoreGetLastestSchemaSettingsFunc) nextHook() func(context.Context, api.SettingsSubject) (*schema.Settings, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *SettingStoreGetLastestSchemaSettingsFunc) appendCall(r0 SettingStoreGetLastestSchemaSettingsFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of
// SettingStoreGetLastestSchemaSettingsFuncCall objects describing the
// invocations of this function.
func (f *SettingStoreGetLastestSchemaSettingsFunc) History() []SettingStoreGetLastestSchemaSettingsFuncCall {
	f.mutex.Lock()
	history := make([]SettingStoreGetLastestSchemaSettingsFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// SettingStoreGetLastestSchemaSettingsFuncCall is an object that describes
// an invocation of method GetLastestSchemaSettings on an instance of
// MockSettingStore.
type SettingStoreGetLastestSchemaSettingsFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 api.SettingsSubject
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 *schema.Settings
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c SettingStoreGetLastestSchemaSettingsFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c SettingStoreGetLastestSchemaSettingsFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// SettingStoreGetLatestFunc describes the behavior when the GetLatest
// method of the parent MockSettingStore instance is invoked.
type SettingStoreGetLatestFunc struct {
	defaultHook func(context.Context, api.SettingsSubject) (*api.Settings, error)
	hooks       []func(context.Context, api.SettingsSubject) (*api.Settings, error)
	history     []SettingStoreGetLatestFuncCall
	mutex       sync.Mutex
}

// GetLatest delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockSettingStore) GetLatest(v0 context.Context, v1 api.SettingsSubject) (*api.Settings, error) {
	r0, r1 := m.GetLatestFunc.nextHook()(v0, v1)
	m.GetLatestFunc.appendCall(SettingStoreGetLatestFuncCall{v0, v1, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the GetLatest method of
// the parent MockSettingStore instance is invoked and the hook queue is
// empty.
func (f *SettingStoreGetLatestFunc) SetDefaultHook(hook func(context.Context, api.SettingsSubject) (*api.Settings, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// GetLatest method of the parent MockSettingStore instance invokes the hook
// at the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *SettingStoreGetLatestFunc) PushHook(hook func(context.Context, api.SettingsSubject) (*api.Settings, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *SettingStoreGetLatestFunc) SetDefaultReturn(r0 *api.Settings, r1 error) {
	f.SetDefaultHook(func(context.Context, api.SettingsSubject) (*api.Settings, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *SettingStoreGetLatestFunc) PushReturn(r0 *api.Settings, r1 error) {
	f.PushHook(func(context.Context, api.SettingsSubject) (*api.Settings, error) {
		return r0, r1
	})
}

func (f *SettingStoreGetLatestFunc) nextHook() func(context.Context, api.SettingsSubject) (*api.Settings, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *SettingStoreGetLatestFunc) appendCall(r0 SettingStoreGetLatestFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of SettingStoreGetLatestFuncCall objects
// describing the invocations of this function.
func (f *SettingStoreGetLatestFunc) History() []SettingStoreGetLatestFuncCall {
	f.mutex.Lock()
	history := make([]SettingStoreGetLatestFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// SettingStoreGetLatestFuncCall is an object that describes an invocation
// of method GetLatest on an instance of MockSettingStore.
type SettingStoreGetLatestFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 api.SettingsSubject
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 *api.Settings
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c SettingStoreGetLatestFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c SettingStoreGetLatestFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}
