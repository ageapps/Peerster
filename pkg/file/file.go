package file

import (
	"github.com/ageapps/Peerster/pkg/logger"
)

const SharedFilesDir = "_SharedFiles"
const chunksDir = "._Chunks"
const metafileDir = "._Metafiles"

// File struct
type File struct {
	filename string
	metadata *Metadata
}

// NewFile create
func NewFileAsync(name string) *File {

	metadata, err := newMetadata(name)
	if err != nil {
		logger.Logf("error creating metadata %v", err)
	}

	go func() {
		if err = metadata.loadMetadata(); err != nil {
			logger.Logf("error loading metadata %v", err)
		}
	}()
	return &File{
		filename: name,
		metadata: metadata,
	}
}

func NewFileSync(name string) *File {
	metadata, err := newMetadata(name)
	if err != nil {
		logger.Logf("error creating metadata %v", err)
	}
	if err = metadata.loadMetadata(); err != nil {
		logger.Logf("error loading metadata %v", err)
	}
	return &File{
		filename: name,
		metadata: metadata,
	}
}
