package gossiper

import (
	"fmt"
	"log"
	"math/rand"
	"sync"

	"github.com/ageapps/Peerster/data"
	"github.com/ageapps/Peerster/logger"
	"github.com/ageapps/Peerster/router"
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
	rumorStack      RumorStack
	privateStack    PrivateStack
	router          *router.Router
	monguerPocesses map[string]*MongerHandler
	rumorCounter    *utils.Counter // [name] adress
	privateCounter  *utils.Counter // [name] adress
	mux             sync.Mutex
	usedPeers       map[string]bool
}

// NewGossiper return new instance
func NewGossiper(addressStr, name string, simple bool) (*Gossiper, error) {
	address, err1 := utils.GetPeerAddress(addressStr)
	connection, err2 := NewConnectionHandler(addressStr, name)
	switch {
	case err1 != nil:
		return nil, err1
	case err2 != nil:
		return nil, err2
	}
	logger.Log(fmt.Sprint("Running simple mode: ", simple))

	return &Gossiper{
		peerConection:   connection,
		Name:            name,
		Address:         address,
		simpleMode:      simple,
		peers:           &utils.PeerAddresses{},
		rumorStack:      RumorStack{Messages: make(map[string][]data.RumorMessage)},
		privateStack:    PrivateStack{Messages: make(map[string][]data.PrivateMessage)},
		router:          router.NewRouter(),
		monguerPocesses: make(map[string]*MongerHandler),
		rumorCounter:    utils.NewCounter(uint32(0)),
		privateCounter:  utils.NewCounter(uint32(0)),
		usedPeers:       make(map[string]bool),
	}, nil
}

// Kill func
func (gossiper *Gossiper) Kill() {
	if gossiper == nil {
		return
	}
	logger.Log("Finishing Gossiper " + gossiper.Name)
	for _, process := range gossiper.getMongerProcesses() {
		process.Stop()
	}
	gossiper.peerConection.Close()
}

// SetPeers peers
func (gossiper *Gossiper) SetPeers(newPeers *utils.PeerAddresses) {
	gossiper.peers = newPeers
}

// AddPeer func
func (gossiper *Gossiper) AddPeer(newPeer string) {
	gossiper.peers.Set(newPeer)
	gossiper.sendStatusMessage(newPeer)
}

// ListenToPeers function
// Start listening to Packets from peers
func (gossiper *Gossiper) ListenToPeers() {
	for {
		packet := &data.GossipPacket{}
		address, err := gossiper.peerConection.readPacket(packet)
		if err != nil {
			logger.Log("Error reading packet")
			break
		}
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
				logger.Log("Error reading message")
				break
			}
			go gossiper.HandleClientMessage(msg)
		}
	}
}

func (gossiper *Gossiper) HandleClientMessage(msg *data.Message) {
	logger.Log("Message received from client")
	logger.LogClient(*msg)

	if gossiper.simpleMode {
		newMsg := data.NewSimpleMessage(gossiper.Name, msg.Text, gossiper.Address.String())
		gossiper.peerConection.broadcastPacket(gossiper.peers, &data.GossipPacket{Simple: newMsg}, gossiper.Address.String())
	} else if msg.IsPrivate() {
		// Message is private
		id := gossiper.privateCounter.Increment()
		privateMessage := data.NewPrivateMessage(gossiper.Name, id, msg.Destination, msg.Text, uint32(10))
		gossiper.privateStack.AddMessage(*privateMessage)
		gossiper.sendPrivateMessage(privateMessage)
	} else {
		id := gossiper.rumorCounter.Increment()
		rumorMessage := data.NewRumorMessage(gossiper.Name, id, msg.Text)
		gossiper.rumorStack.AddMessage(*rumorMessage)
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
		gossiper.handleRumorMessage(packet.Rumor, origin)
	case data.PACKET_PRIVATE:
		gossiper.handlePrivateMessage(packet.Private, origin)
	case data.PACKET_SIMPLE:
		logger.LogSimple(*packet.Simple)
		logger.LogPeers(gossiper.peers.String())
		gossiper.handleSimpleMessage(packet.Simple, origin)
	default:
		logger.Log("Message not recognized")
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

func (gossiper *Gossiper) handlePrivateMessage(msg *data.PrivateMessage, from string) {
	if msg.Destination == gossiper.Name {
		gossiper.privateStack.AddMessage(*msg)
		logger.LogPrivate(*msg)
		return
	}
	msg.HopLimit--
	if msg.HopLimit > 0 {
		gossiper.sendPrivateMessage(msg)
	}
}

func (gossiper *Gossiper) handleRumorMessage(msg *data.RumorMessage, from string) {
	isRouteRumor := msg.IsRouteRumor()
	if !isRouteRumor {
		logger.LogRumor(*msg, from)
		logger.LogPeers(gossiper.peers.String())
	}
	msgStatus := gossiper.rumorStack.CompareMessage(msg.Origin, msg.ID)

	if msgStatus == NEW_MESSAGE {
		newEntry := gossiper.router.SetEntry(msg.Origin, from)
		if isRouteRumor {
			logger.Log(fmt.Sprintf("Received ROUTE RUMOR - new:%v", newEntry))
			return
		}
		// If I get own messages that i didn´t
		// have, set internal rumorCounter
		if msg.Origin == gossiper.Name {
			gossiper.rumorCounter.SetValue(msg.ID)
		}
		logger.Log("Received new rumor message, appending...")
		gossiper.resetusedPeers()

		gossiper.rumorStack.AddMessage(*msg)
		gossiper.mongerMessage(msg, from)
	}

	gossiper.sendStatusMessage(from)
}

func (gossiper *Gossiper) handleStatusMessage(msg *data.StatusPacket, from string) {
	handler := gossiper.findMonguerHandler(from)
	logger.Log(fmt.Sprint("Handler found:", handler != nil))

	if len(msg.Want) < len(*gossiper.rumorStack.getRumorStack()) {

		// check messages that i have from other peers that aren´t in the status message
		for origin, messages := range *gossiper.rumorStack.getRumorStack() {
			firstMessageID := messages[0].ID
			found := false
			for _, status := range msg.Want {
				if status.Identifier == origin {
					found = true
					break
				}
			}
			if !found {
				if handler != nil {
					handler.SetSynking(true)
				}
				logger.Log(fmt.Sprintf("Peer needs to update Origin:%v - ID:%v", origin, firstMessageID))
				gossiper.sendRumrorMessage(from, origin, firstMessageID)
				return
			}
		}
	}
	logger.LogStatus(*msg, from)
	logger.LogPeers(gossiper.peers.String())
	inSync := true

	for _, status := range msg.Want {
		messageStatus := gossiper.rumorStack.CompareMessage(status.Identifier, uint32(status.NextID-1))

		switch messageStatus {
		case NEW_MESSAGE:
			// logger.Log("Gossiper needs to update")
			if handler != nil {
				handler.SetSynking(true)
			}
			gossiper.sendStatusMessage(from)
			break
		case IN_SYNC:
			// logger.Log("Gossiper and Peer have same messages")
		case OLD_MESSAGE:
			// logger.Log("Peer needs to update")
			if handler != nil {
				handler.SetSynking(true)
			}
			gossiper.sendRumrorMessage(from, status.Identifier, status.NextID)
			break
		}
		inSync = inSync && messageStatus == IN_SYNC
	}
	if inSync {
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
func (gossiper *Gossiper) resetusedPeers() {
	gossiper.mux.Lock()
	gossiper.usedPeers = make(map[string]bool)
	gossiper.mux.Unlock()
}
func (gossiper *Gossiper) getUsedPeers() map[string]bool {
	gossiper.mux.Lock()
	defer gossiper.mux.Unlock()
	return gossiper.usedPeers
}

func (gossiper *Gossiper) GetLatestMessages() *[]data.RumorMessage {
	return gossiper.rumorStack.GetLatestMessages()
}
func (gossiper *Gossiper) GetPrivateMessages() *map[string][]data.PrivateMessage {
	return gossiper.privateStack.getPrivateStack()
}

func (gossiper *Gossiper) GetPeers() *[]string {
	var peersArr = []string{}
	for _, peer := range gossiper.peers.GetAdresses() {
		peersArr = append(peersArr, peer.String())
	}
	return &peersArr
}

func (gossiper *Gossiper) GetRoutes() *router.RoutingTable {
	return gossiper.router.GetTable()
}

func (gossiper *Gossiper) sendStatusMessage(destination string) {
	var message = gossiper.rumorStack.getStatusMessage()
	packet := &data.GossipPacket{Status: message}
	go gossiper.peerConection.sendPacketToPeer(destination, packet)
}

func (gossiper *Gossiper) sendRumrorMessage(destinationAdress, origin string, id uint32) {
	if message := gossiper.rumorStack.GetRumorMessage(origin, id); message != nil {
		packet := &data.GossipPacket{Rumor: message}
		logger.Log(fmt.Sprintf("Sending RUMOR ID:%v", message.ID))
		go gossiper.peerConection.sendPacketToPeer(destinationAdress, packet)
	} else {
		logger.Log("Message to send not found")
	}
}

func (gossiper *Gossiper) sendRouteRumorMessage(destinationAdress string) {
	latestMsgID := gossiper.rumorCounter.GetValue() + 1
	routeRumorMessage := data.NewRumorMessage(gossiper.Name, latestMsgID, "")
	packet := &data.GossipPacket{Rumor: routeRumorMessage}
	logger.Log(fmt.Sprintf("Sending ROUTE RUMOR ID:%v", latestMsgID))
	go gossiper.peerConection.sendPacketToPeer(destinationAdress, packet)
}

func (gossiper *Gossiper) sendPrivateMessage(msg *data.PrivateMessage) {
	packet := &data.GossipPacket{Private: msg}
	if destinationAdress, ok := gossiper.router.GetDestination(msg.Destination); ok {
		logger.Log(fmt.Sprintf("Sending PRIVATE Dest:%v", msg.Destination))
		go gossiper.peerConection.sendPacketToPeer(destinationAdress.String(), packet)
	}

}

func (gossiper *Gossiper) mongerMessage(msg *data.RumorMessage, originPeer string) {
	gossiper.mux.Lock()
	processName := fmt.Sprint(len(gossiper.monguerPocesses))
	logger.Log("Starting monger process - " + processName)
	monguerProcess := NewMongerHandler(originPeer, processName, msg, gossiper.peerConection, gossiper.peers)
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
func (gossiper *Gossiper) getMongerProcesses() map[string]*MongerHandler {
	gossiper.mux.Lock()
	defer gossiper.mux.Unlock()
	return gossiper.monguerPocesses
}
