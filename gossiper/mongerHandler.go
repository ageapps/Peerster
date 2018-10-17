package gossiper

import (
	"math/rand"
	"sync"

	"github.com/ageapps/Peerster/utils"

	"github.com/ageapps/Peerster/data"
)

// MongerHandler is a handler that will be in
// charge of the monguering process whenever the
// gossiper gets a message from a client
type MongerHandler struct {
	originalMessage *data.RumorMessage
	currentMessage  *data.RumorMessage
	isMongering     bool
	connection      *ConnectionHandler
	peers           *utils.PeerAddresses
	mux             sync.Mutex
}

func NewMongerHandler(msg *data.RumorMessage, peerConection *ConnectionHandler, connectPeers *utils.PeerAddresses) *MongerHandler {
	return &MongerHandler{
		originalMessage: msg,
		currentMessage:  msg,
		isMongering:     false,
		connection:      peerConection,
		peers:           connectPeers,
	}
}

func (handler *MongerHandler) start() {
	// usedPeers := make(map[string]bool)
	// index := gossiper.peers.GetRandomPeer(usedPeers)
	// peer := gossiper.peers.Addresses[index]

	for {
		// msg := gossiper.MongerHandler.GetMonguerMessage()
		// // gossiper.MongerHandler.SetPendingMessage(peer.String(), msg.ID)
		// // logger.LogMonguer(peer.String())
		// go gossiper.peerConection.SendPacketToPeer(peer.String(), &data.GossipPacket{Rumor: msg})
		// //Sleep while waiting for message
		// time.Sleep(1 * time.Second)
		// // Check if message was acknowledged
		// if gossiper.MongerHandler.IsMessagePending(peer.String()) {
		// 	go gossiper.sendPacketToPeer(peer.String(), &data.GossipPacket{Rumor: msg})
		// } else {
		// 	break
		// }
	}
}

//SetMonguerMessage function
func (handler *MongerHandler) SetMonguerMessage(msg *data.RumorMessage) {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	handler.currentMessage = msg
}

//GetMonguerMessage function
func (handler *MongerHandler) GetMonguerMessage() *data.RumorMessage {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	return handler.currentMessage
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

func (gossiper *Gossiper) keepRumorering() bool {
	// flipCoin
	coin := rand.Int() % 2
	return coin != 0
}

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
