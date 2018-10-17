package gossiper

import (
	"fmt"
	"log"
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
	peers           utils.PeerAddresses
	messageStack    *MessageStack
	monguerPocesses []*MongerHandler
	counter         *utils.Counter // [name] adress
	mux             sync.Mutex
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
	newCounter := utils.NewCounter(uint32(0))
	return &Gossiper{
		Name:          name,
		Address:       address,
		peerConection: connection,
		simpleMode:    simple,
		counter:       newCounter,
	}, nil
}

// Set gossiper peers
func (gossiper *Gossiper) SetPeers(newPeers utils.PeerAddresses) {
	gossiper.peers = newPeers
}

// ListenToPeers function
// Start listening to Packets from peers
func (gossiper *Gossiper) ListenToPeers() {
	packet := &data.GossipPacket{}
	for {
		if address, err := gossiper.peerConection.ReadPacket(packet); err != nil {
			log.Fatal(err)
		} else {
			go gossiper.handlePeerPacket(packet, address)
		}
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
		msg := &data.Message{}
		logger.Log("Listening to client in " + clientAddress)
		for {
			_, err := connection.ReadMessage(msg)
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
		go gossiper.peerConection.BroadcastPacket(&gossiper.peers, &data.GossipPacket{Simple: newMsg}, gossiper.Address.String())
	} else {
		id := gossiper.counter.Increment()
		rumorMessage := data.NewRumorMessage(gossiper.Name, id, msg.Text)
		gossiper.messageStack.AddMessage(*rumorMessage)
		go gossiper.mongerMessage(rumorMessage)
	}

}

func (gossiper *Gossiper) handlePeerPacket(packet *data.GossipPacket, origin string) {
	err := gossiper.peers.Set(origin)
	if err != nil {
		log.Fatal(err)
	}
	switch {
	case packet.IsStatusMessage():
		go gossiper.handleStatusMessage(packet.Status, origin)
	case packet.IsRumorMessage():
		logger.LogRumor(*packet.Rumor, origin)
		logger.LogPeers(gossiper.peers)
		go gossiper.handleRumorMessage(*packet.Rumor, origin)
	default:
		logger.LogSimple(*packet.Simple)
		logger.LogPeers(gossiper.peers)
		go gossiper.handleSimpleMessage(packet.Simple, origin)
	}
}

func (gossiper *Gossiper) handleSimpleMessage(msg *data.SimpleMessage, from string) {
	if msg.OriginalName == gossiper.Name {
		logger.Log("Received own message")
		return
	}
	newMsg := data.NewSimpleMessage(msg.OriginalName, msg.Contents, gossiper.Address.String())
	gossiper.peerConection.BroadcastPacket(&gossiper.peers, &data.GossipPacket{Simple: newMsg}, msg.RelayPeerAddr)
}

func (gossiper *Gossiper) handleRumorMessage(msg data.RumorMessage, from string) {
	gossiper.messageStack.AddMessage(msg)
	var message = gossiper.messageStack.getStatusMessage()
	packet := &data.GossipPacket{Status: message}
	logger.Log("Sending Status Message to <" + from + ">")
	go gossiper.peerConection.SendPacketToPeer(from, packet)
}

func (gossiper *Gossiper) handleStatusMessage(msg *data.StatusPacket, from string) {
	for _, status := range msg.Want {
		logger.LogStatus(status, from)
		logger.LogPeers(gossiper.peers)
		// gossiper.MongerHandler.ConfirmPendingMessage(from, status.NextID)
		if gossiper.messageStack.AreMessagesMissing(from, status.NextID) {
		}

	}
}

//Reset function
func (gossiper *Gossiper) StartEntropyTimer() {
	usedPeers := make(map[string]bool)
	for {
		time.Sleep(1 * time.Second)
		randPeerIndex := gossiper.peers.GetRandomPeer(usedPeers)
		newpeer := gossiper.peers.Addresses[randPeerIndex]
		usedPeers[newpeer.String()] = true
		var message = gossiper.messageStack.getStatusMessage()
		packet := &data.GossipPacket{Status: message}
		go gossiper.peerConection.SendPacketToPeer(newpeer.String(), packet)
		if len(usedPeers) >= len(gossiper.peers.Addresses) {
			break
		}
	}
}

func (gossiper *Gossiper) mongerMessage(msg *data.RumorMessage) {
	monguerProcess := NewMongerHandler(msg, gossiper.peerConection, &gossiper.peers)
	gossiper.monguerPocesses = append(gossiper.monguerPocesses, monguerProcess)
	gossiper.monguerPocesses[len(gossiper.monguerPocesses)-1].start()
}
