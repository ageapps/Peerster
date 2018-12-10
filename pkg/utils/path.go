package utils

import (
	"path"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

// GetRootPath get
func GetRootPath() string {
	return path.Join(basepath, "..", "..")
}

// GetFilesPath get
func GetFilesPath() string {
	// return path.Join(GetRootPath(), "files")
	return GetRootPath()
}
