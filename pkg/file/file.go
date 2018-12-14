package file

import (
	"io/ioutil"
	"log"
	"path"
	"strings"

	"github.com/ageapps/Peerster/pkg/utils"
)

// File struct
type File struct {
	Name         string `json:"name"`
	Size         int64
	MetafileHash []byte
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
