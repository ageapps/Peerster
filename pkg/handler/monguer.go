package handler

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/ageapps/Peerster/pkg/data"
	"github.com/ageapps/Peerster/pkg/logger"
	"github.com/ageapps/Peerster/pkg/utils"
)

var usedPeers = make(map[string]bool)

// MongerHandler is a handler that will be in
// charge of the monguering process whenever the
// gossiper gets a message from a client:
// name                   name of the Handler
// originalMessage        original message that was being monguered
// currentMessage         message currently being monguered
// currentPeer            client currently being monguered
// active                 monguer handler active state
// currentlySynchronicing bool
// connection             *ConnectionHandler
// peers                  *utils.PeerAddresses
// mux                    sync.Mutex
// timer                  *time.Timer
// quitChannel            chan bool
// resetChannel           chan bool
//
type MongerHandler struct {
	Name                   string
	currentMessage         *data.RumorMessage
	currentPeer            string
	active                 bool
	routeMonguer           bool
	currentlySynchronicing bool
	connection             *ConnectionHandler
	peers                  *utils.PeerAddresses
	mux                    sync.Mutex
	timer                  *time.Timer
	quitChannel            chan bool
	resetChannel           chan bool
	StopChannel           chan bool
	usedPeers              *map[string]bool
}

// NewMongerHandler function
func NewMongerHandler(originPeer, nameStr string, isRouter bool, msg *data.RumorMessage, peerConection *ConnectionHandler, connectPeers *utils.PeerAddresses) *MongerHandler {
	used := make(map[string]bool)
	if originPeer != "" {
		used[originPeer] = true
	}
	return &MongerHandler{
		Name:                   nameStr,
		currentMessage:         msg,
		currentPeer:            "",
		active:                 false,
		routeMonguer:           isRouter,
		currentlySynchronicing: false,
		connection:             peerConection,
		peers:                  connectPeers,
		timer:                  &time.Timer{},
		quitChannel:            make(chan bool),
		StopChannel:            make(chan bool),
		resetChannel:           make(chan bool),
		usedPeers:              &used,
	}
}

// Start monguering process
func (handler *MongerHandler) Start() {
	go func() {
		go handler.setActive(true)
		handler.monguerWithPeer(false)
		for {
			select {
			case <-handler.resetChannel:
				logger.Log("Restarting monger handler - " + handler.Name)
				handler.monguerWithPeer(true)
			case <-handler.quitChannel:
				logger.Log("Finishing monger handler - " + handler.Name)
				if handler.timer.C != nil {
					handler.timer.Stop()
				}
				return
			case <-handler.timer.C:
				// Flip coin
				if !handler.isSynking() {
					logger.Log("TIMEOUT, FLIPPING COIN")
					if !keepRumorering() {
						handler.Stop()
					} else {
						handler.monguerWithPeer(true)
					}
				}
			}
		}
	}()
}

func newTimer() *time.Timer {
	// logger.Log("Launching new timer")
	return time.NewTimer(1 * time.Second)
}

func (handler *MongerHandler) monguerWithPeer(flipped bool) {
	if peer := handler.GetPeers().GetRandomPeer(*handler.usedPeers); peer != nil {
		handler.timer = newTimer()
		handler.setMonguerPeer(peer.String())
		handler.addUsedPeer(peer.String())
		// logger.Log(fmt.Sprint("Monguering with peer: ", peer.String()))
		packet := &data.GossipPacket{Rumor: handler.getMonguerMessage()}
		if !flipped {
			logger.LogMonguer(peer.String())
		} else {
			logger.LogCoin(peer.String())
		}
		handler.connection.SendPacketToPeer(peer.String(), packet)
	} else {
		logger.Log(fmt.Sprint("No peers to monger with"))
		handler.Stop()
	}
}

// Stop handler
func (handler *MongerHandler) Stop() {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	logger.Log("Stopping monger handler")
	go handler.setActive(false)
	go func() {
		close(handler.quitChannel)
		close(handler.StopChannel)
	}()
}

// Reset handler
func (handler *MongerHandler) Reset() {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	logger.Log("Restart monger handler")
	go func() {
		handler.resetChannel <- true
	}()
}

//SetMonguerMessage function
func (handler *MongerHandler) SetMonguerMessage(msg *data.RumorMessage) {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	handler.currentMessage = msg
}

//GetMonguerMessage function
func (handler *MongerHandler) getMonguerMessage() *data.RumorMessage {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	// logger.Log(fmt.Sprint("Monger message is ", handler.currentMessage))
	return handler.currentMessage
}

// GetMonguerPeer function
func (handler *MongerHandler) GetMonguerPeer() string {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	return handler.currentPeer
}

// GetPeers function
func (handler *MongerHandler) GetPeers() *utils.PeerAddresses {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	return handler.peers
}

func (handler *MongerHandler) setMonguerPeer(peer string) {
	handler.mux.Lock()
	handler.currentPeer = peer
	handler.mux.Unlock()
}

// IsActive gets handler status
func (handler *MongerHandler) IsActive() bool {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	return handler.active
}

// IsRouteMonguer gets handler status
func (handler *MongerHandler) IsRouteMonguer() bool {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	return handler.routeMonguer
}

func (handler *MongerHandler) isSynking() bool {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	return handler.currentlySynchronicing
}

// SetSynking state
func (handler *MongerHandler) SetSynking(value bool) {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	handler.timer.Stop()
	handler.currentlySynchronicing = value
}
func (handler *MongerHandler) setActive(value bool) {
	handler.mux.Lock()
	handler.active = value
	handler.mux.Unlock()
}
func (handler *MongerHandler) addUsedPeer(peer string) {
	handler.mux.Lock()
	(*handler.usedPeers)[peer] = true
	handler.mux.Unlock()
}

func (handler *MongerHandler) logMonguer(msg string) {
	logger.Log(fmt.Sprintf("[MONGER-%v]%v", handler.Name, msg))
}

func keepRumorering() bool {
	// flipCoin
	coin := rand.Int() % 2
	return coin != 0
}
