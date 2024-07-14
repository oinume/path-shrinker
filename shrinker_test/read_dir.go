package shrinker_test

import (
	"os"
	"time"
)

func NewMockFileInfo(name string, size int64, mode os.FileMode, modTime time.Time, isDir bool) *MockFileInfo {
	return &MockFileInfo{
		name:    name,
		size:    size,
		mode:    mode,
		modTime: modTime,
		isDir:   isDir,
	}
}

type MockFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (m *MockFileInfo) Name() string {
	return m.name
}

func (m *MockFileInfo) Size() int64 {
	return m.size
}

func (m *MockFileInfo) Mode() os.FileMode {
	return m.mode
}

func (m *MockFileInfo) ModTime() time.Time {
	return m.modTime
}

func (m *MockFileInfo) IsDir() bool {
	return m.isDir
}

func (m *MockFileInfo) Sys() interface{} {
	panic("implement me")
}

/* TODO: Create DirEntry
https://pkg.go.dev/io/fs#DirEntry or os.DirEntry

type DirEntry interface {
	// Name returns the name of the file (or subdirectory) described by the entry.
	// This name is only the final element of the path (the base name), not the entire path.
	// For example, Name would return "hello.go" not "home/gopher/hello.go".
	Name() string

	// IsDir reports whether the entry describes a directory.
	IsDir() bool

	// Type returns the type bits for the entry.
	// The type bits are a subset of the usual FileMode bits, those returned by the FileMode.Type method.
	Type() FileMode

	// Info returns the FileInfo for the file or subdirectory described by the entry.
	// The returned FileInfo may be from the time of the original directory read
	// or from the time of the call to Info. If the file has been removed or renamed
	// since the directory read, Info may return an error satisfying errors.Is(err, ErrNotExist).
	// If the entry denotes a symbolic link, Info reports the information about the link itself,
	// not the link's target.
	Info() (FileInfo, error)
}
*/
