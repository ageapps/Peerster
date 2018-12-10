package handler

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ageapps/Peerster/pkg/file"

	"github.com/ageapps/Peerster/pkg/logger"
	"github.com/ageapps/Peerster/pkg/router"

	"github.com/ageapps/Peerster/pkg/data"
)

const MaxBudget = 32
const DefaultBudget = 2

// SearchHandler is a handler that will be in
// charge of requesting data from other peers
// FileName            string
// MetaHash            data.HashValue
// active              bool
// currentlyRequesting bool
// connection          *ConnectionHandler
// router              *router.Router
// mux                 sync.Mutex
// quitChannel         chan bool
// resetChannel        chan bool
// StopChannel         chan bool
//
type SearchHandler struct {
	file                *file.File
	origin              string
	chunk               int
	metaHash            data.HashValue
	gotMetafile         bool
	currentPeer         string
	active              bool
	currentlyRequesting bool
	connection          *ConnectionHandler
	router              *router.Router
	mux                 sync.Mutex
	timer               *time.Timer
	quitChannel         chan bool
	StopChannel         chan bool
	ChunkChannel        chan data.Chunk
}

// NewSearchHandler function
func NewSearchHandler(filename, origin, destination string, hash data.HashValue, peerConection *ConnectionHandler, router *router.Router) *DataHandler {
	return &DataHandler{
		file:                file.NewDownloadingFile(filename),
		origin:              origin,
		chunk:               0,
		currentPeer:         destination,
		metaHash:            hash,
		active:              false,
		currentlyRequesting: false,
		connection:          peerConection,
		router:              router,
		timer:               &time.Timer{},
		quitChannel:         make(chan bool),
		StopChannel:         make(chan bool),
		ChunkChannel:        make(chan data.Chunk),
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
func (handler *DataHandler) Start() {
	go handler.resetTimer()
	handler.sendRequest()
	handler.handleTimeout()
	go func() {
		for chunk := range handler.ChunkChannel {
			go handler.resetTimer()
			handler.handleTimeout()
			logger.Logf("RECEIVED CHUNK")

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
		close(handler.ChunkChannel)
		close(handler.StopChannel)
	}()

}

// GetMetaHashStr get
func (handler *DataHandler) GetMetaHashStr() string {
	return handler.metaHash.String()
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
	requestHash := handler.metaHash
	if handler.gotMetafile {
		requestHash = handler.file.GetChunkHash(handler.chunk)
	}
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
