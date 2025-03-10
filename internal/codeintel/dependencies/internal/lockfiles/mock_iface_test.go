// Code generated by go-mockgen 1.3.1; DO NOT EDIT.
//
// This file was generated by running `sg generate` (or `go-mockgen`) at the root of
// this repository. To add additional mocks to this or another package, add a new entry
// to the metadata.yaml file in the root of this repository.

package lockfiles

import (
	"context"
	"io"
	"sync"

	api "github.com/sourcegraph/sourcegraph/internal/api"
	gitserver "github.com/sourcegraph/sourcegraph/internal/gitserver"
)

// MockGitService is a mock implementation of the GitService interface (from
// the package
// github.com/sourcegraph/sourcegraph/internal/codeintel/dependencies/internal/lockfiles)
// used for unit testing.
type MockGitService struct {
	// ArchiveFunc is an instance of a mock function object controlling the
	// behavior of the method Archive.
	ArchiveFunc *GitServiceArchiveFunc
	// LsFilesFunc is an instance of a mock function object controlling the
	// behavior of the method LsFiles.
	LsFilesFunc *GitServiceLsFilesFunc
}

// NewMockGitService creates a new mock of the GitService interface. All
// methods return zero values for all results, unless overwritten.
func NewMockGitService() *MockGitService {
	return &MockGitService{
		ArchiveFunc: &GitServiceArchiveFunc{
			defaultHook: func(context.Context, api.RepoName, gitserver.ArchiveOptions) (r0 io.ReadCloser, r1 error) {
				return
			},
		},
		LsFilesFunc: &GitServiceLsFilesFunc{
			defaultHook: func(context.Context, api.RepoName, api.CommitID, ...gitserver.Pathspec) (r0 []string, r1 error) {
				return
			},
		},
	}
}

// NewStrictMockGitService creates a new mock of the GitService interface.
// All methods panic on invocation, unless overwritten.
func NewStrictMockGitService() *MockGitService {
	return &MockGitService{
		ArchiveFunc: &GitServiceArchiveFunc{
			defaultHook: func(context.Context, api.RepoName, gitserver.ArchiveOptions) (io.ReadCloser, error) {
				panic("unexpected invocation of MockGitService.Archive")
			},
		},
		LsFilesFunc: &GitServiceLsFilesFunc{
			defaultHook: func(context.Context, api.RepoName, api.CommitID, ...gitserver.Pathspec) ([]string, error) {
				panic("unexpected invocation of MockGitService.LsFiles")
			},
		},
	}
}

// NewMockGitServiceFrom creates a new mock of the MockGitService interface.
// All methods delegate to the given implementation, unless overwritten.
func NewMockGitServiceFrom(i GitService) *MockGitService {
	return &MockGitService{
		ArchiveFunc: &GitServiceArchiveFunc{
			defaultHook: i.Archive,
		},
		LsFilesFunc: &GitServiceLsFilesFunc{
			defaultHook: i.LsFiles,
		},
	}
}

// GitServiceArchiveFunc describes the behavior when the Archive method of
// the parent MockGitService instance is invoked.
type GitServiceArchiveFunc struct {
	defaultHook func(context.Context, api.RepoName, gitserver.ArchiveOptions) (io.ReadCloser, error)
	hooks       []func(context.Context, api.RepoName, gitserver.ArchiveOptions) (io.ReadCloser, error)
	history     []GitServiceArchiveFuncCall
	mutex       sync.Mutex
}

// Archive delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockGitService) Archive(v0 context.Context, v1 api.RepoName, v2 gitserver.ArchiveOptions) (io.ReadCloser, error) {
	r0, r1 := m.ArchiveFunc.nextHook()(v0, v1, v2)
	m.ArchiveFunc.appendCall(GitServiceArchiveFuncCall{v0, v1, v2, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the Archive method of
// the parent MockGitService instance is invoked and the hook queue is
// empty.
func (f *GitServiceArchiveFunc) SetDefaultHook(hook func(context.Context, api.RepoName, gitserver.ArchiveOptions) (io.ReadCloser, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Archive method of the parent MockGitService instance invokes the hook at
// the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *GitServiceArchiveFunc) PushHook(hook func(context.Context, api.RepoName, gitserver.ArchiveOptions) (io.ReadCloser, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *GitServiceArchiveFunc) SetDefaultReturn(r0 io.ReadCloser, r1 error) {
	f.SetDefaultHook(func(context.Context, api.RepoName, gitserver.ArchiveOptions) (io.ReadCloser, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *GitServiceArchiveFunc) PushReturn(r0 io.ReadCloser, r1 error) {
	f.PushHook(func(context.Context, api.RepoName, gitserver.ArchiveOptions) (io.ReadCloser, error) {
		return r0, r1
	})
}

func (f *GitServiceArchiveFunc) nextHook() func(context.Context, api.RepoName, gitserver.ArchiveOptions) (io.ReadCloser, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *GitServiceArchiveFunc) appendCall(r0 GitServiceArchiveFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of GitServiceArchiveFuncCall objects
// describing the invocations of this function.
func (f *GitServiceArchiveFunc) History() []GitServiceArchiveFuncCall {
	f.mutex.Lock()
	history := make([]GitServiceArchiveFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// GitServiceArchiveFuncCall is an object that describes an invocation of
// method Archive on an instance of MockGitService.
type GitServiceArchiveFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 api.RepoName
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 gitserver.ArchiveOptions
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 io.ReadCloser
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c GitServiceArchiveFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c GitServiceArchiveFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// GitServiceLsFilesFunc describes the behavior when the LsFiles method of
// the parent MockGitService instance is invoked.
type GitServiceLsFilesFunc struct {
	defaultHook func(context.Context, api.RepoName, api.CommitID, ...gitserver.Pathspec) ([]string, error)
	hooks       []func(context.Context, api.RepoName, api.CommitID, ...gitserver.Pathspec) ([]string, error)
	history     []GitServiceLsFilesFuncCall
	mutex       sync.Mutex
}

// LsFiles delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockGitService) LsFiles(v0 context.Context, v1 api.RepoName, v2 api.CommitID, v3 ...gitserver.Pathspec) ([]string, error) {
	r0, r1 := m.LsFilesFunc.nextHook()(v0, v1, v2, v3...)
	m.LsFilesFunc.appendCall(GitServiceLsFilesFuncCall{v0, v1, v2, v3, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the LsFiles method of
// the parent MockGitService instance is invoked and the hook queue is
// empty.
func (f *GitServiceLsFilesFunc) SetDefaultHook(hook func(context.Context, api.RepoName, api.CommitID, ...gitserver.Pathspec) ([]string, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// LsFiles method of the parent MockGitService instance invokes the hook at
// the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *GitServiceLsFilesFunc) PushHook(hook func(context.Context, api.RepoName, api.CommitID, ...gitserver.Pathspec) ([]string, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *GitServiceLsFilesFunc) SetDefaultReturn(r0 []string, r1 error) {
	f.SetDefaultHook(func(context.Context, api.RepoName, api.CommitID, ...gitserver.Pathspec) ([]string, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *GitServiceLsFilesFunc) PushReturn(r0 []string, r1 error) {
	f.PushHook(func(context.Context, api.RepoName, api.CommitID, ...gitserver.Pathspec) ([]string, error) {
		return r0, r1
	})
}

func (f *GitServiceLsFilesFunc) nextHook() func(context.Context, api.RepoName, api.CommitID, ...gitserver.Pathspec) ([]string, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *GitServiceLsFilesFunc) appendCall(r0 GitServiceLsFilesFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of GitServiceLsFilesFuncCall objects
// describing the invocations of this function.
func (f *GitServiceLsFilesFunc) History() []GitServiceLsFilesFuncCall {
	f.mutex.Lock()
	history := make([]GitServiceLsFilesFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// GitServiceLsFilesFuncCall is an object that describes an invocation of
// method LsFiles on an instance of MockGitService.
type GitServiceLsFilesFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 api.RepoName
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 api.CommitID
	// Arg3 is a slice containing the values of the variadic arguments
	// passed to this method invocation.
	Arg3 []gitserver.Pathspec
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []string
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation. The variadic slice argument is flattened in this array such
// that one positional argument and three variadic arguments would result in
// a slice of four, not two.
func (c GitServiceLsFilesFuncCall) Args() []interface{} {
	trailing := []interface{}{}
	for _, val := range c.Arg3 {
		trailing = append(trailing, val)
	}

	return append([]interface{}{c.Arg0, c.Arg1, c.Arg2}, trailing...)
}

// Results returns an interface slice containing the results of this
// invocation.
func (c GitServiceLsFilesFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}
