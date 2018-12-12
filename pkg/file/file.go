package file

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"

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
	Name         string `json:"name"`
	metadata     *Metadata
	Size         int64
	MetafileHash []byte
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

// Exists func
func Exists(name string) (string, bool) {
	files, err := ioutil.ReadDir(path.Join(utils.GetRootPath(), ChunksDir))
	if err != nil {
		log.Fatal(err)
	}
	for _, fileName := range files {
		if fileName.Name() == name {
			return path.Join(utils.GetRootPath(), ChunksDir, name), true
		}
	}
	return "", false
}

// GetChunkCount file name
func (file *File) GetChunkCount() uint64 {
	return uint64(len(file.metadata.fileHashes))
}

// MatchKeyword file name
func (file *File) MatchKeyword(match string) bool {
	return strings.Contains(file.Name, match)
}

// AddChunk to metadata
func (file *File) AddChunk(chunk []byte, hash utils.HashValue) error {
	return file.metadata.addChunk(chunk, hash, false)
}

// AddMetafile to metadata
func (file *File) AddMetafile(chunk []byte, hash utils.HashValue) error {
	return file.metadata.addMetafile(chunk, hash)
}

// GetChunkHash get hash of chunk
func (file *File) GetChunkHash(index int) utils.HashValue {
	return file.metadata.fileHashes[index]
}

// GetChunkMap get hash of chunk
func (file *File) GetChunkMap() []uint64 {
	return file.metadata.chunkMap
}

// GetMetaHash of metadata
func (file *File) GetMetaHash() utils.HashValue {
	return file.metadata.metahash
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
	logger.LogReconstructed(file.Name)
	return nil
}
