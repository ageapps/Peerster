package file

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"sync"

	"github.com/ageapps/Peerster/pkg/logger"
	"github.com/ageapps/Peerster/pkg/utils"
)

// Metadata struct
type Metadata struct {
	filename string
	size     int64
	metafile MetaFile
	metahash MetaHash
	mux      sync.Mutex
}

// MetaFile is a file containing the SHA-256 hashes of each chunk
type MetaFile []byte

// MetaHash is a file containing the SHA-256 hashes of each chunk
type MetaHash []byte

var chunckSize int64 = 8192 // 8kb

func newMetadata(filename string) (*Metadata, error) {
	fileSize, err := getFileSize(filename)
	if err != nil {
		return nil, fmt.Errorf("getting file size: %v", err)
	}
	return &Metadata{
		filename: filename,
		size:     fileSize,
	}, nil
}

func (meta *Metadata) loadMetadata() error {
	if meta.filename == "" {
		return fmt.Errorf("error filename in empty")
	}

	// 1. Open file
	file, err := os.Open(path.Join(utils.GetFilesPath(), SharedFilesDir, meta.filename))
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	fileSize := meta.size
	chunkNumber := int(math.Ceil(float64(fileSize) / float64(chunckSize)))

	logger.Log(fmt.Sprintf("File of size %v / %v chuncks", fileSize, chunkNumber))
	// 2. Read file
	reader := bufio.NewReader(file)
	var metafile MetaFile
	bytesNotRead := fileSize

	for bytesNotRead > 0 {

		bufferSize := chunckSize
		if bytesNotRead < chunckSize {
			bufferSize = bytesNotRead
		}

		chunk := make([]byte, bufferSize)
		_, err := reader.Read(chunk)
		if err != nil {
			return fmt.Errorf("error reading file: %v", err)
		}

		hashArr := sha256.Sum256(chunk)
		hash := hashArr[:]
		metafile = append(metafile, hash...)
		bytesNotRead -= chunckSize
		chunkFileName := path.Join(utils.GetFilesPath(), chunksDir, hex.EncodeToString(hash))
		err = ioutil.WriteFile(chunkFileName, chunk, 0644)
		if err != nil {
			return fmt.Errorf("error during write of chunk: %v", err)
		}
		logger.Log(fmt.Sprintf("Chunk saved: %v", chunkFileName))

	}
	metahash := sha256.Sum256(metafile)
	metahashFileName := path.Join(utils.GetFilesPath(), metafileDir, hex.EncodeToString(metahash[:]))
	err = ioutil.WriteFile(metahashFileName, metafile, 0644)
	if err != nil {
		return fmt.Errorf("error during write of metafile: %v", err)
	}
	logger.Log(fmt.Sprintf("Metafile saved: %v", metahashFileName))

	meta.mux.Lock()
	meta.metafile = metafile
	meta.metahash = metahash[:]
	meta.mux.Unlock()

	return nil
}

func getFileSize(filename string) (int64, error) {
	fi, e := os.Stat(utils.GetFilesPath() + "/_SharedFiles/" + filename)
	if e != nil {
		return -1, fmt.Errorf("error estracting data from file %v, ", e)
	}
	return fi.Size(), nil
}
