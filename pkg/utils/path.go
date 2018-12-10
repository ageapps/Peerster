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

func GetRootPath() string {
	return path.Join(basepath, "..", "..")
}

func GetFilesPath() string {
	// return path.Join(GetRootPath(), "files")
	return GetRootPath()
}
