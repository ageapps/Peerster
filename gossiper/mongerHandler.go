package gossiper

import (
	"fmt"
	"sync"
	"time"

	"github.com/ageapps/Peerster/logger"
	"github.com/ageapps/Peerster/utils"

	"github.com/ageapps/Peerster/data"
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
	name                   string
	originalMessage        *data.RumorMessage
	currentMessage         *data.RumorMessage
	currentPeer            string
	active                 bool
	currentlySynchronicing bool
	connection             *ConnectionHandler
	peers                  *utils.PeerAddresses
	mux                    sync.Mutex
	timer                  *time.Timer
	quitChannel            chan bool
	resetChannel           chan bool
}

// NewMongerHandler function
func NewMongerHandler(nameStr string, msg *data.RumorMessage, peerConection *ConnectionHandler, connectPeers *utils.PeerAddresses) *MongerHandler {
	return &MongerHandler{
		name:                   nameStr,
		originalMessage:        msg,
		currentMessage:         msg,
		active:                 false,
		currentlySynchronicing: false,
		connection:             peerConection,
		peers:                  connectPeers,
		quitChannel:            make(chan bool),
		resetChannel:           make(chan bool),
	}
}

func (handler *MongerHandler) start() {
	go func() {
		handler.setActive(true)
		handler.monguerWithPeer(false)
		for {
			select {
			case <-handler.resetChannel:
				handler.mux.Lock()
				logger.Log("Restarting monger handler - " + handler.name)
				go handler.monguerWithPeer(true)
				handler.mux.Unlock()
			case <-handler.quitChannel:
				handler.mux.Lock()
				logger.Log("Finishing monger handler - " + handler.name)
				handler.timer.Stop()
				handler.setActive(false)
				handler.mux.Unlock()
				return
			case <-handler.timer.C:
				// Flip coin
				if !handler.isSynking() {
					logger.Log("TIMEOUT, FLIPPING COIN")
					if !keepRumorering() {
						handler.Stop()
					} else {
						go handler.monguerWithPeer(true)
					}
				}
			}
		}
	}()
}
func newTimer() *time.Timer {
	logger.Log("Launching new timer")
	return time.NewTimer(1 * time.Second)
}
func (handler *MongerHandler) monguerWithPeer(flipped bool) {
	handler.mux.Lock()
	handler.timer = newTimer()
	peer := handler.peers.GetRandomPeer(usedPeers)
	handler.currentPeer = peer.String()
	handler.mux.Unlock()
	logger.Log(fmt.Sprint("Monguering with peer: ", peer.String()))
	packet := &data.GossipPacket{Rumor: handler.getMonguerMessage()}
	if !flipped {
		logger.LogMonguer(peer.String())
	} else {
		logger.LogCoin(peer.String())
	}
	go handler.connection.sendPacketToPeer(peer.String(), packet)
}

func (handler *MongerHandler) Stop() {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	logger.Log("Stopping monger handler")
	go func() {
		handler.quitChannel <- true
	}()
}
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
	logger.Log(fmt.Sprint("Monger message is ", handler.currentMessage))
	return handler.currentMessage
}

func (handler *MongerHandler) IsActive() bool {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	return handler.active
}
func (handler *MongerHandler) isSynking() bool {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	return handler.currentlySynchronicing
}
func (handler *MongerHandler) SetSynking(value bool) {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	handler.timer.Stop()
	handler.currentlySynchronicing = value
}
func (handler *MongerHandler) setActive(value bool) {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	handler.active = value
}

//SetPendingMessage function
func (handler *MongerHandler) SetPendingMessage(address string, id uint32) string {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	//handler.pendingMessages[address] = id
	return address
}

//ConfirmPendingMessage function
func (handler *MongerHandler) ConfirmPendingMessage(address string, id uint32) {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	// if handler.pendingMessages[address] < id {
	// 	logger.Log("Confirming message from <" + address + "> with ID: " + fmt.Sprint(handler.pendingMessages[address]))
	// 	delete(handler.pendingMessages, address)
	// }
}

// IsMessagePending func
// func (handler *MongerHandler) IsMessagePending(address string) bool {
// 	handler.mux.Lock()
// 	defer handler.mux.Unlock()
// 	_, ok := handler.pendingMessages[address]
// 	return ok
// }

// GetPendingMessage func
// func (handler *MongerHandler) GetPendingMessage(address string) uint32 {
// 	handler.mux.Lock()
// 	defer handler.mux.Unlock()
// 	value, _ := handler.pendingMessages[address]
// 	return value
// }

// func (g *Gossiper) rumourMongering(msg data.RumorMessage) {
// 	addr := msg.Addr
// 	gp := msg.Msg
// 	usedNeighbours := make(map[string]bool)
// 	if addr != "" {
// 		usedNeighbours[addr] = true
// 	}
// 	randPeer := g.Neighbours.RandomIndexOutOfNeighbours(usedNeighbours)
// 	fmt.Printf("\n Mongering with %v \n", randPeer)
// 	usedNeighbours[randPeer] = true
// 	g.sendMessageToNeighbour(gp, randPeer)
// 	status := g.Status
// 	status.ChangeStatus(randPeer)
// 	var brk bool
// 	for {
// 		if randPeer == "" {
// 			return
// 		}
// 		//Sleep while waiting for message
// 		time.Sleep(1 * time.Second)
// 		if brk {
// 			break
// 		}
// 		select {
// 		case msg := <-status.StatusChannel:
// 			fmt.Println(msg.Msg.Status)
// 			needMsgs := g.Messages.NeedMsgs(*msg.Msg.Status)
// 			gp := &data.GossipPacket{
// 				Status: &needMsgs,
// 			}
// 			g.sendMessageToNeighbour(gp, randPeer)
// 			g.RetrieveMongerMessages()
// 			randPeer = g.Neighbours.RandomIndexOutOfNeighbours(usedNeighbours)
// 			usedNeighbours[randPeer] = true
// 		default:
// 			coin := rand.Int() % 2
// 			if coin == 0 {
// 				g.Status.ChangeStatus("")
// 				g.Status.StopMongering()
// 				brk = true
// 			} else {
// 				randPeer = g.Neighbours.RandomIndexOutOfNeighbours(usedNeighbours)
// 				usedNeighbours[randPeer] = true
// 				fmt.Printf("FLIPPED COIN sending rumor to %v \n", randPeer)
// 				if randPeer != "" {
// 					g.sendMessageToNeighbour(gp, randPeer)
// 				}
// 			}
// 		}
// 	}
// }
