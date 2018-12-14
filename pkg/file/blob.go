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

// SharedBlobsDir directory where shared files are stored
const SharedBlobsDir = "_SharedFiles"

// ChunksDir directory where Chunks are stored
const ChunksDir = "._Chunks"

// DownloadsDir directory where reconstructed files are stored
const DownloadsDir = "_Downloads"

const metafileDir = "._Metafiles"

// Blob struct
type Blob struct {
	fileName     string
	metadata     *Metadata
	MetafileHash []byte
}

// NewBlobFromLocalAsync create
func NewBlobFromLocalAsync(name string) (*Blob, error) {

	metadata, err := newMetadata(name, true)
	if err != nil {
		return nil, fmt.Errorf("error creating metadata %v", err)
	}

	go func() {
		if err = metadata.loadMetadata(); err != nil {
			logger.Logf("error loading metadata %v", err)
		}
	}()
	return &Blob{
		fileName: name,
		metadata: metadata,
	}, nil
}

// NewBlobFromLocalSync create
func NewBlobFromLocalSync(name string) (*Blob, error) {
	metadata, err := newMetadata(name, true)
	if err != nil {
		return nil, fmt.Errorf("error creating metadata %v", err)
	}
	if err = metadata.loadMetadata(); err != nil {
		return nil, fmt.Errorf("error loading metadata %v", err)
	}
	return &Blob{
		fileName:     name,
		metadata:     metadata,
		MetafileHash: metadata.metahash,
	}, nil
}

// NewDownloadingBlob create
func NewDownloadingBlob(name string) *Blob {
	metadata, err := newMetadata(name, false)
	if err != nil {
		log.Fatal(err)
	}
	return &Blob{
		fileName: name,
		metadata: metadata,
	}
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
func (file *Blob) GetChunkCount() uint64 {
	return uint64(len(file.metadata.fileHashes))
}

// MatchKeyword file name
func (blob *Blob) MatchKeyword(match string) bool {
	return strings.Contains(blob.fileName, match)
}

// AddChunk to metadata
func (blob *Blob) AddChunk(chunk []byte, hash utils.HashValue) error {
	return blob.metadata.addChunk(chunk, hash, false)
}

// AddMetafile to metadata
func (blob *Blob) AddMetafile(chunk []byte, hash utils.HashValue) error {
	blob.MetafileHash = hash
	return blob.metadata.addMetafile(chunk, hash)
}

// GetChunkHash get hash of chunk
func (blob *Blob) GetChunkHash(index int) utils.HashValue {
	return blob.metadata.fileHashes[index]
}

// GetChunkMap get hash of chunk
func (blob *Blob) GetChunkMap() []uint64 {
	return blob.metadata.chunkMap
}

// GetBlobSize get hash of chunk
func (blob *Blob) GetBlobSize() int64 {
	return blob.metadata.size
}

// GetName get hash of chunk
func (blob *Blob) GetName() string {
	return blob.fileName
}

// GetMetaHash of metadata
func (blob *Blob) GetMetaHash() utils.HashValue {
	return blob.MetafileHash
}

// Reconstruct file from metadata
func (blob *Blob) Reconstruct() (int64, error) {
	var buffer []byte
	logger.Logf("Reconstructing blob: %v, chunks: %v", blob.fileName, len(blob.metadata.fileHashes))
	for _, hash := range blob.metadata.fileHashes {
		chunkPath := path.Join(utils.GetBlobsPath(), ChunksDir, hash.String())
		//logger.Logf("Reading hash file: %v", chunkPath)
		b, err := ioutil.ReadFile(chunkPath) // just pass the file name
		if err != nil {
			return 0, fmt.Errorf("Error reading chunks to reconstruct blob : %v", err)
		}
		buffer = append(buffer, b...)
	}
	fileSize := int64(len(buffer))
	blob.metadata.size = fileSize
	blobPath := path.Join(utils.GetBlobsPath(), DownloadsDir, blob.fileName)
	if err := blob.metadata.saveBlob(blobPath, buffer); err != nil {
		return fileSize, fmt.Errorf("Error writing reconstructed file : %v", err)
	}
	logger.Logf("Reconstructed blob saved in: %v", blobPath)
	logger.LogReconstructed(blob.fileName)
	return fileSize, nil
}
