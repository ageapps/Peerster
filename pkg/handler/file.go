package handler

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ageapps/Peerster/pkg/utils"

	"github.com/ageapps/Peerster/pkg/data"
	"github.com/ageapps/Peerster/pkg/file"
	"github.com/ageapps/Peerster/pkg/logger"
	"github.com/ageapps/Peerster/pkg/router"
)

// FileHandler is a handler that will be in
// charge of requesting data from other peers
// FileName            string
// MetaHash            utils.HashValue
// stopped              bool
// connection          *ConnectionHandler
// router              *router.Router
// mux                 sync.Mutex
// quitChannel         chan bool
// resetChannel        chan bool
//
type FileHandler struct {
	Name          string
	file          *file.File
	blob          *file.Blob
	origin        string
	chunk         int
	expectingHash utils.HashValue
	gotMetafile   bool
	currentPeer   string
	stopped       bool
	connection    *ConnectionHandler
	router        *router.Router
	mux           sync.Mutex
	timer         *time.Timer
	quitChannel   chan bool
	ChunkChannel  chan utils.Chunk
}

// NewDataHandler function
func NewDataHandler(name, filename, origin, destination string, hash utils.HashValue, peerConection *ConnectionHandler, router *router.Router) *FileHandler {
	return &FileHandler{
		Name:          name,
		file:          file.NewFile(filename, 0, []byte{}),
		blob:          file.NewDownloadingBlob(filename),
		origin:        origin,
		chunk:         0,
		currentPeer:   destination,
		expectingHash: hash,
		stopped:       false,
		connection:    peerConection,
		router:        router,
		timer:         &time.Timer{},
		quitChannel:   make(chan bool),
		ChunkChannel:  make(chan utils.Chunk),
	}
}

func (handler *FileHandler) resetTimer() {
	//logger.Log("Launching new timer")
	if handler.getTimer().C != nil {
		handler.getTimer().Stop()
	}
	handler.mux.Lock()
	handler.timer = time.NewTimer(5 * time.Second)
	handler.mux.Unlock()
}
func (handler *FileHandler) getTimer() *time.Timer {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	return handler.timer
}

// Start handler
func (handler *FileHandler) Start(onStopHandler func()) {
	go handler.resetTimer()
	handler.sendRequest()
	handler.handleTimeout()
	go func() {
		for chunk := range handler.ChunkChannel {
			go handler.resetTimer()
			handler.handleTimeout()

			// First chunk received is the metafile
			if !handler.gotMetafile && chunk.Data != nil {
				if err := handler.blob.AddMetafile(chunk.Data, chunk.Hash); err != nil {
					log.Fatal(err)
				}
				handler.file.SetMetaHash(chunk.Hash)
				logger.Logf("Saved Metafile: %v", chunk.Hash.String())
				handler.gotMetafile = true
				logger.LogMetafile(handler.file.Name, handler.currentPeer)
				handler.sendRequest()
				// Normal chunk received
			} else {
				if err := handler.blob.AddChunk(chunk.Data, chunk.Hash); err != nil {
					log.Fatal(err)
				}
				logger.Logf("Added Chunk: %v", chunk.Hash.String())
				logger.LogChunk(handler.file.Name, handler.currentPeer, handler.chunk+1)
				handler.chunk++
				if int64(len(chunk.Data)) < file.ChunckSize {
					// last chunk of file
					handler.getTimer().Stop()
					go func() {
						size, err := handler.blob.Reconstruct()
						if err != nil {
							log.Fatal(err)
						}
						handler.file.SetSize(size)
					}()
					break
				}
				handler.sendRequest()
			}
		}
		if !handler.stopped {
			close(handler.ChunkChannel)
		}
		onStopHandler()
	}()

}

// Stop func
func (handler *FileHandler) Stop() {
	logger.Log("Stopping data handler")
	if !handler.stopped {
		handler.stopped = true
		close(handler.ChunkChannel)
		return
	}
	logger.Log("Data Handler already stopped....")
}

// GetExpectingHashStr get
func (handler *FileHandler) GetExpectingHashStr() string {
	return handler.expectingHash.String()
}

// GetFile get
func (handler *FileHandler) GetFile() *file.File {
	return handler.file
}

// GetBlob get
func (handler *FileHandler) GetBlob() *file.Blob {
	return handler.blob
}

// GetCurrentPeer get
func (handler *FileHandler) GetCurrentPeer() string {
	return handler.currentPeer
}

func (handler *FileHandler) sendRequest() {
	if handler.gotMetafile {
		handler.expectingHash = handler.blob.GetChunkHash(handler.chunk)
	}
	requestHash := handler.expectingHash
	logger.Log(fmt.Sprintf("Sending DATA REQUEST Hash:%v", requestHash.String()))
	msg := data.NewDataRequest(handler.origin, handler.currentPeer, uint32(10), requestHash)
	packet := &data.GossipPacket{DataRequest: msg}
	if destinationAdress, ok := handler.router.GetDestination(msg.Destination); ok {
		logger.Logf("Sending DATA REQUEST Dest:%v", msg.Destination)
		handler.connection.SendPacketToPeer(destinationAdress.String(), packet)
	} else {
		logger.Logf("INVALID DATA REQUEST Dest:%v", msg.Destination)
	}
}

func (handler *FileHandler) handleTimeout() {
	go func() {
		<-handler.getTimer().C
		fmt.Println("TIMEOUT requesting file")
		handler.sendRequest()
		handler.resetTimer()
	}()
}
