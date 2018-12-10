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

	"github.com/ageapps/Peerster/pkg/data"
	"github.com/ageapps/Peerster/pkg/logger"
	"github.com/ageapps/Peerster/pkg/utils"
)

// Metadata struct
type Metadata struct {
	filename   string
	size       int64
	metafile   data.HashValue
	metahash   data.HashValue
	fileHashes []data.HashValue
	mux        sync.Mutex
}

// SHAhashSize size of SHA hash
const SHAhashSize = sha256.Size

// ChunckSize Size of chunks files are splitted to
var ChunckSize int64 = 8192 // 8kb

func newMetadata(filename string, local bool) (*Metadata, error) {
	meta := &Metadata{
		filename:   filename,
		metafile:   []byte{},
		fileHashes: []data.HashValue{},
	}
	if local {
		fileSize, err := getFileSize(filename)
		if err != nil {
			return nil, fmt.Errorf("error getting file size: %v", err)
		}
		meta.size = fileSize
	}
	return meta, nil
}

func (meta *Metadata) loadMetadata() error {
	if meta.filename == "" {
		return fmt.Errorf("error filename in empty")
	}
	if meta.size == 0 {
		return fmt.Errorf("error file size in empty")
	}

	// 1. Open file
	filePath := path.Join(utils.GetFilesPath(), SharedFilesDir, meta.filename)
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	fileSize := meta.size
	chunkNumber := int(math.Ceil(float64(fileSize) / float64(ChunckSize)))

	logger.Log(fmt.Sprintf("File of size %v / %v chuncks", fileSize, chunkNumber))
	// 2. Read file
	reader := bufio.NewReader(file)
	bytesNotRead := fileSize

	logger.Logf("Reading file in path %v", filePath)

	for bytesNotRead > 0 {

		bufferSize := ChunckSize
		if bytesNotRead < ChunckSize {
			bufferSize = bytesNotRead
		}

		chunk := make([]byte, bufferSize)
		_, err := reader.Read(chunk)
		if err != nil {
			return fmt.Errorf("error reading file: %v", err)
		}

		hashArr := sha256.Sum256(chunk)
		hash := hashArr[:]

		if err = meta.addChunk(chunk, hash); err != nil {
			return fmt.Errorf("error adding chunk: %v", err)
		}
		bytesNotRead -= ChunckSize
	}
	metahash := sha256.Sum256(meta.metafile)

	err = meta.addMetafile(meta.metafile, metahash[:])
	if err != nil {
		return fmt.Errorf("error during write of metafile: %v", err)
	}
	// err = meta.saveMetafile(meta.metafile, metahash[:])
	// if err != nil {
	// 	return fmt.Errorf("error during write of metafile: %v", err)
	// }
	logger.Log("Metafile saved")

	return nil
}

func getFileSize(filename string) (int64, error) {
	fi, e := os.Stat(utils.GetFilesPath() + "/_SharedFiles/" + filename)
	if e != nil {
		return -1, fmt.Errorf("error estracting data from file %v, ", e)
	}
	return fi.Size(), nil
}

func (meta *Metadata) saveFile(path string, chunk []byte) error {
	err := ioutil.WriteFile(path, chunk, 0644)
	if err != nil {
		return fmt.Errorf("error during write of chunk: %v", err)
	}
	//logger.Log(fmt.Sprintf("Chunk saved: %v", chunkFileName))
	return nil
}

func (meta *Metadata) addChunk(chunk []byte, hash data.HashValue) error {
	meta.mux.Lock()
	meta.metafile = append(meta.metafile, hash...)
	//logger.Logf("Saving hash v% : %v", len(hash), hex.EncodeToString(hash))
	//meta.fileHashes = append(meta.fileHashes, hash)
	meta.mux.Unlock()
	chunkFilePath := path.Join(utils.GetFilesPath(), ChunksDir, hex.EncodeToString(hash))
	return meta.saveFile(chunkFilePath, chunk)
}
func (meta *Metadata) saveMetafile(data []byte, hash data.HashValue) error {
	meta.mux.Lock()
	meta.metahash = hash
	meta.mux.Unlock()
	metahashFilePath := path.Join(utils.GetFilesPath(), ChunksDir, hex.EncodeToString(hash))
	metahashBackupFilePath := path.Join(utils.GetFilesPath(), metafileDir, hex.EncodeToString(hash))
	go meta.saveFile(metahashBackupFilePath, data)
	return meta.saveFile(metahashFilePath, data)
}

func (meta *Metadata) addMetafile(data []byte, hash data.HashValue) error {
	logger.Logf("Reading metadata file")
	hashNumber := len(data) / SHAhashSize
	logger.Logf("Reading metadata with %v hashes", hashNumber)
	for index := 0; index < hashNumber; index++ {
		meta.mux.Lock()
		startIndex := index * SHAhashSize
		endIndex := startIndex + SHAhashSize
		hash := data[startIndex:endIndex]
		//logger.Logf("Reading hash: %v", hex.EncodeToString(hash))
		meta.fileHashes = append(meta.fileHashes, hash)
		meta.mux.Unlock()
	}
	return meta.saveMetafile(data, hash)
}
