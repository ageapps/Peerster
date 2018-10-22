package gossiper

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/ageapps/Peerster/data"
	"github.com/ageapps/Peerster/logger"
	"github.com/ageapps/Peerster/utils"
)

// Gossiper struct
type Gossiper struct {
	peerConection   *ConnectionHandler
	clientConection *ConnectionHandler
	Name            string
	Address         utils.PeerAddress
	simpleMode      bool
	peers           *utils.PeerAddresses
	messageStack    MessageStack
	monguerPocesses map[string]*MongerHandler
	counter         *utils.Counter // [name] adress
	mux             sync.Mutex
	usedPeers       map[string]bool
}

// NewGossiper return new instance
func NewGossiper(addressStr, name string, simple bool) (*Gossiper, error) {
	address, err1 := utils.GetPeerAddress(addressStr)
	connection, err2 := NewConnectionHandler(addressStr, name)
	logger.Log(fmt.Sprint("Running simple mode: ", simple))

	switch {
	case err1 != nil:
		return nil, err1
	case err2 != nil:
		return nil, err2
	}
	return &Gossiper{
		peerConection:   connection,
		Name:            name,
		Address:         address,
		simpleMode:      simple,
		peers:           &utils.PeerAddresses{},
		messageStack:    MessageStack{Messages: make(map[string][]data.RumorMessage)},
		monguerPocesses: make(map[string]*MongerHandler),
		counter:         utils.NewCounter(uint32(0)),
		usedPeers:       make(map[string]bool),
	}, nil
}

// Set gossiper peers
func (gossiper *Gossiper) SetPeers(newPeers *utils.PeerAddresses) {
	gossiper.peers = newPeers
}

// ListenToPeers function
// Start listening to Packets from peers
func (gossiper *Gossiper) ListenToPeers() {
	for {
		packet := &data.GossipPacket{}
		address := gossiper.peerConection.readPacket(packet)
		gossiper.handlePeerPacket(packet, address)
	}
}

// ListenToClients function
// start to listen for client messages
// in desired port
func (gossiper *Gossiper) ListenToClients(port int) {
	address := gossiper.Address.IP
	clientAddress := fmt.Sprintf("%v:%v", address.String(), port)
	if connection, err := NewConnectionHandler(clientAddress, gossiper.Name+"-Client"); err != nil {
		log.Fatal(err)
	} else {
		logger.Log("Listening to client in " + clientAddress)
		for {
			msg := &data.Message{}
			_, err := connection.readMessage(msg)
			if err != nil {
				log.Fatal(err)
			}
			go gossiper.handleClientMessage(msg)
		}
	}
}

func (gossiper *Gossiper) handleClientMessage(msg *data.Message) {
	logger.Log("Message received from client")
	logger.LogClient(*msg)

	if gossiper.simpleMode {
		newMsg := data.NewSimpleMessage(gossiper.Name, msg.Text, gossiper.Address.String())
		gossiper.peerConection.broadcastPacket(gossiper.peers, &data.GossipPacket{Simple: newMsg}, gossiper.Address.String())
	} else {
		id := gossiper.counter.Increment()
		rumorMessage := data.NewRumorMessage(gossiper.Name, id, msg.Text)
		gossiper.messageStack.AddMessage(*rumorMessage)
		gossiper.mongerMessage(rumorMessage, "")
	}

}

func (gossiper *Gossiper) handlePeerPacket(packet *data.GossipPacket, origin string) {
	err := gossiper.peers.Set(origin)
	if err != nil {
		log.Fatal(err)
	}

	packetType := packet.GetPacketType()
	logger.Log("Peer packet received: " + packetType)

	switch packetType {
	case data.PACKET_STATUS:
		gossiper.handleStatusMessage(packet.Status, origin)
	case data.PACKET_RUMOR:
		logger.LogRumor(*packet.Rumor, origin)
		logger.LogPeers(gossiper.peers.String())
		gossiper.handleRumorMessage(packet.Rumor, origin)
	case data.PACKET_SIMPLE:
		logger.LogSimple(*packet.Simple)
		logger.LogPeers(gossiper.peers.String())
		gossiper.handleSimpleMessage(packet.Simple, origin)
	default:
		// if packet.Rumor != nil {
		// 	fmt.Println(*packet.Rumor)
		// }
		// if packet.Status != nil {
		// 	fmt.Println(*packet.Status)
		// }
		// fmt.Println(*packet)
		fmt.Println("Message not recognized")
		// log.Fatal(errors.New("Message not recognized"))
	}
}

func (gossiper *Gossiper) handleSimpleMessage(msg *data.SimpleMessage, from string) {
	if msg.OriginalName == gossiper.Name {
		logger.Log("Received own message")
		return
	}
	newMsg := data.NewSimpleMessage(msg.OriginalName, msg.Contents, gossiper.Address.String())
	gossiper.peerConection.broadcastPacket(gossiper.peers, &data.GossipPacket{Simple: newMsg}, msg.RelayPeerAddr)
}

func (gossiper *Gossiper) handleRumorMessage(msg *data.RumorMessage, from string) {
	if gossiper.messageStack.CompareMessage(msg.Origin, msg.ID) == NEW_MESSAGE {
		logger.Log("Received new message, appending...")
		// If I get own messages that i didn´t
		// have, set internal counter
		if msg.Origin == gossiper.Name {
			gossiper.counter.SetValue(msg.ID)
		}
		gossiper.messageStack.AddMessage(*msg)
		gossiper.mongerMessage(msg, from)
	}
	gossiper.sendStatusMessage(from)
}

func (gossiper *Gossiper) handleStatusMessage(msg *data.StatusPacket, from string) {
	handler := gossiper.findMonguerHandler(from)
	logger.Log(fmt.Sprint("Handler found:", handler != nil))
	if len(msg.Want) <= 0 {
		if handler != nil {
			handler.SetSynking(true)
		}
		// send my own first message
		for key, messages := range *gossiper.messageStack.getMessageStack() {
			lastMessageID := messages[0].ID
			logger.Log(fmt.Sprintf("Peer needs to update - ID:%v", lastMessageID))
			gossiper.sendRumrorMessage(from, key, lastMessageID)
			break
		}
	}
	logger.LogStatus(*msg, from)
	logger.LogPeers(gossiper.peers.String())
	inSync := true
	for _, status := range msg.Want {
		messageStatus := gossiper.messageStack.CompareMessage(status.Identifier, uint32(status.NextID-1))

		switch messageStatus {
		case NEW_MESSAGE:
			logger.Log("Gossiper needs to update")
			if handler != nil {
				handler.SetSynking(true)
			}
			gossiper.sendStatusMessage(from)
			break
		case IN_SYNC:
			logger.Log("Gossiper and Peer have same messages")
		case OLD_MESSAGE:
			logger.Log("Peer needs to update")
			if handler != nil {
				handler.SetSynking(true)
			}
			gossiper.sendRumrorMessage(from, status.Identifier, status.NextID)
			break
		}
		inSync = inSync && messageStatus==IN_SYNC
	}
	if inSync{
		logger.LogInSync(from)
			if handler != nil {
				handler.SetSynking(false)
				// Flip coin
				logger.Log("IN SYNC, FLIPPING COIN")
				if !keepRumorering() {
					handler.Stop()
					// delete handler from slice
					gossiper.deleteMongerProcess(handler.name)
				} else {
					handler.Reset()
				}
			}
	}
}

// StartEntropyTimer function
func (gossiper *Gossiper) StartEntropyTimer() {
	go func() {
		logger.Log("Starting Entropy timer")
		for {
			if len(gossiper.usedPeers) >= len(gossiper.peers.GetAdresses()) {
				// logger.Log("Entropy Timer - All peers where notified")
			}
			if newpeer := gossiper.peers.GetRandomPeer(gossiper.usedPeers); newpeer != nil {
				logger.Log("Entropy Timer - MESSAGE")
				gossiper.usedPeers[newpeer.String()] = true
				gossiper.sendStatusMessage(newpeer.String())
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

func (gossiper *Gossiper) sendStatusMessage(destination string) {
	var message = gossiper.messageStack.getStatusMessage()
	// fmt.Println(message)
	packet := &data.GossipPacket{Status: message}
	go gossiper.peerConection.sendPacketToPeer(destination, packet)
}

func (gossiper *Gossiper) sendRumrorMessage(destinationAdress, origin string, id uint32) {
	if message := gossiper.messageStack.GetRumorMessage(origin, id); message != nil {
		packet := &data.GossipPacket{Rumor: message}
		logger.Log(fmt.Sprintf("Sending RUMOR ID:%v", message.ID))
		go gossiper.peerConection.sendPacketToPeer(destinationAdress, packet)
	} else {
		logger.Log("Message to send not found")
	}
}

func (gossiper *Gossiper) mongerMessage(msg *data.RumorMessage, originPeer string) {
	gossiper.mux.Lock()
	processName := fmt.Sprint(len(gossiper.monguerPocesses))
	logger.Log("Starting monger process - " + processName)
	monguerProcess := NewMongerHandler(originPeer,processName, msg, gossiper.peerConection, gossiper.peers)
	gossiper.monguerPocesses[processName] = monguerProcess
	gossiper.mux.Unlock()
	gossiper.monguerPocesses[processName].start()
}

func (gossiper *Gossiper) findMonguerHandler(origin string) *MongerHandler {
	for _, process := range gossiper.getMongerProcesses() {
		// delete inactive peer
		if !process.IsActive() {
			gossiper.deleteMongerProcess(process.name)
			continue
		}
		if process.currentPeer == origin {
			return process
		}
	}
	return nil
}

func keepRumorering() bool {
	// flipCoin
	coin := rand.Int() % 2
	return coin != 0
}

func (gossiper *Gossiper) deleteMongerProcess(name string) {
	gossiper.mux.Lock()
	defer gossiper.mux.Unlock()
	delete(gossiper.monguerPocesses, name)
}
func (gossiper *Gossiper) getMongerProcesses() map[string]*MongerHandler{
	gossiper.mux.Lock()
	defer gossiper.mux.Unlock()
	return gossiper.monguerPocesses
}
