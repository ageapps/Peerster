package file

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/ageapps/Peerster/pkg/utils"
)

// File struct
type File struct {
	Name         string `json:"name"`
	Size         int64  `json:"size"`
	MetafileHash []byte `json:"hash"`
}

// FileStatus struct
type FileStatus struct {
	Name         string `json:"name"`
	Size         int64  `json:"size"`
	MetafileHash []byte `json:"hash"`
	Blob         bool   `json:"blob"`
}

// GetMetaHash of metadata
func (file *File) GetMetaHash() utils.HashValue {
	return file.MetafileHash
}

// SetMetaHash of metadata
func (file *File) SetMetaHash(hash []byte) {
	file.MetafileHash = hash
}

// SetSize of metadata
func (file *File) SetSize(size int64) {
	file.Size = size
}

// NewFile create
func NewFile(name string, size int64, metahash []byte) *File {
	return &File{name, size, metahash}
}

// MatchKeyword func
func MatchKeyword(match string) (string, bool) {
	files, err := ioutil.ReadDir(path.Join(utils.GetRootPath(), DownloadsDir))
	if err != nil {
		log.Fatal(err)
	}
	// logger.Logf("Looking for %v", name)

	for _, fileName := range files {
		if strings.Contains(fileName.Name(), match) {
			return path.Join(utils.GetRootPath(), ChunksDir, match), true
		}
	}
	return "", false
}

// Copy file
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
