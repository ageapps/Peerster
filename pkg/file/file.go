package file

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/ageapps/Peerster/pkg/data"
	"github.com/ageapps/Peerster/pkg/logger"
	"github.com/ageapps/Peerster/pkg/utils"
)

// SharedFilesDir directory where shared files are stored
const SharedFilesDir = "_SharedFiles"

// ChunksDir directory where Chunks are stored
const ChunksDir = "._Chunks"

// DownloadsDir directory where reconstructed files are stored
const DownloadsDir = "_Downloads"

const metafileDir = "._Metafiles"

// File struct
type File struct {
	Name     string `json:"name"`
	metadata *Metadata
}

// NewFileFromLocalAsync create
func NewFileFromLocalAsync(name string) (*File, error) {

	metadata, err := newMetadata(name, true)
	if err != nil {
		return nil, fmt.Errorf("error creating metadata %v", err)
	}

	go func() {
		if err = metadata.loadMetadata(); err != nil {
			logger.Logf("error loading metadata %v", err)
		}
	}()
	return &File{
		Name:     name,
		metadata: metadata,
	}, nil
}

// NewFileFromLocalSync create
func NewFileFromLocalSync(name string) (*File, error) {
	metadata, err := newMetadata(name, true)
	if err != nil {
		return nil, fmt.Errorf("error creating metadata %v", err)
	}
	if err = metadata.loadMetadata(); err != nil {
		return nil, fmt.Errorf("error loading metadata %v", err)
	}
	return &File{
		Name:     name,
		metadata: metadata,
	}, nil
}

// NewDownloadingFile create
func NewDownloadingFile(name string) *File {
	metadata, err := newMetadata(name, false)
	if err != nil {
		log.Fatal(err)
	}

	return &File{
		Name:     name,
		metadata: metadata,
	}
}

// AddChunk to metadata
func (file *File) AddChunk(chunk []byte, hash data.HashValue) error {
	return file.metadata.addChunk(chunk, hash)
}

// AddMetafile to metadata
func (file *File) AddMetafile(chunk []byte, hash data.HashValue) error {
	return file.metadata.addMetafile(chunk, hash)
}

// GetChunkHash get hash of chunk
func (file *File) GetChunkHash(index int) data.HashValue {
	return file.metadata.fileHashes[index]
}

// GetMetaHash of metadata
func (file *File) GetMetaHash() string {
	return file.metadata.metahash.String()
}

// Reconstruct file from metadata
func (file *File) Reconstruct() error {
	var buffer []byte
	logger.Logf("Reconstructing file: %v, chunks: %v", file.Name, len(file.metadata.fileHashes))
	for _, hash := range file.metadata.fileHashes {
		chunkPath := path.Join(utils.GetFilesPath(), ChunksDir, hash.String())
		//logger.Logf("Reading hash file: %v", chunkPath)
		b, err := ioutil.ReadFile(chunkPath) // just pass the file name
		if err != nil {
			return fmt.Errorf("Error reading chunks to reconstruct file : %v", err)
		}
		buffer = append(buffer, b...)
	}

	filePath := path.Join(utils.GetFilesPath(), DownloadsDir, file.Name)
	if err := file.metadata.saveFile(filePath, buffer); err != nil {
		return fmt.Errorf("Error writing reconstructed file : %v", err)
	}
	logger.Logf("Reconstructed file saved in: %v", filePath)
	return nil
}
