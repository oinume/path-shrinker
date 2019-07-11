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
	panic("implement me")
}

func (m *MockFileInfo) Mode() os.FileMode {
	panic("implement me")
}

func (m *MockFileInfo) ModTime() time.Time {
	panic("implement me")
}

func (m *MockFileInfo) IsDir() bool {
	return m.isDir
}

func (m *MockFileInfo) Sys() interface{} {
	panic("implement me")
}
