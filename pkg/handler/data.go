package handler

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ageapps/Peerster/pkg/file"
	"github.com/ageapps/Peerster/pkg/utils"

	"github.com/ageapps/Peerster/pkg/logger"
	"github.com/ageapps/Peerster/pkg/router"

	"github.com/ageapps/Peerster/pkg/data"
)

// DataHandler is a handler that will be in
// charge of requesting data from other peers
// FileName            string
// MetaHash            data.HashValue
// stopped              bool
// connection          *ConnectionHandler
// router              *router.Router
// mux                 sync.Mutex
// quitChannel         chan bool
// resetChannel        chan bool
//
type DataHandler struct {
	Name          string
	file          *file.File
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
func NewDataHandler(name, filename, origin, destination string, hash utils.HashValue, peerConection *ConnectionHandler, router *router.Router) *DataHandler {
	return &DataHandler{
		Name:          name,
		file:          file.NewDownloadingFile(filename),
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

func (handler *DataHandler) resetTimer() {
	//logger.Log("Launching new timer")
	if handler.getTimer().C != nil {
		handler.getTimer().Stop()
	}
	handler.mux.Lock()
	handler.timer = time.NewTimer(5 * time.Second)
	handler.mux.Unlock()
}
func (handler *DataHandler) getTimer() *time.Timer {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	return handler.timer
}

// Start handler
func (handler *DataHandler) Start(onStopHandler func()) {
	go handler.resetTimer()
	handler.sendRequest()
	handler.handleTimeout()
	go func() {
		for chunk := range handler.ChunkChannel {
			go handler.resetTimer()
			handler.handleTimeout()

			// First chunk received is the metafile
			if !handler.gotMetafile && chunk.Data != nil {
				if err := handler.file.AddMetafile(chunk.Data, chunk.Hash); err != nil {
					log.Fatal(err)
				}
				logger.Logf("Saved Metafile: %v", chunk.Hash.String())
				handler.gotMetafile = true
				logger.LogMetafile(handler.file.Name, handler.currentPeer)
				handler.sendRequest()
				// Normal chunk received
			} else {
				if err := handler.file.AddChunk(chunk.Data, chunk.Hash); err != nil {
					log.Fatal(err)
				}
				logger.Logf("Added Chunk: %v", chunk.Hash.String())
				logger.LogChunk(handler.file.Name, handler.currentPeer, handler.chunk+1)
				handler.chunk++
				if int64(len(chunk.Data)) < file.ChunckSize {
					// last chunk of file
					handler.getTimer().Stop()
					go handler.file.Reconstruct()
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

func (handler *DataHandler) Stop() {
	logger.Log("Stopping data handler")
	if !handler.stopped {
		handler.stopped = true
		close(handler.ChunkChannel)
		return
	}
	logger.Log("Data Handler already stopped....")
}

// GetExpectingHashStr get
func (handler *DataHandler) GetExpectingHashStr() string {
	return handler.expectingHash.String()
}

// GetFile get
func (handler *DataHandler) GetFile() *file.File {
	return handler.file
}

// GetCurrentPeer get
func (handler *DataHandler) GetCurrentPeer() string {
	return handler.currentPeer
}

func (handler *DataHandler) sendRequest() {
	if handler.gotMetafile {
		handler.expectingHash = handler.file.GetChunkHash(handler.chunk)
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

func (handler *DataHandler) handleTimeout() {
	go func() {
		<-handler.getTimer().C
		fmt.Println("TIMEOUT requesting file")
		handler.sendRequest()
		handler.resetTimer()
	}()
}
