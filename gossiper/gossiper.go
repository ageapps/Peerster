package gossiper

import (
	"fmt"
	"log"
	"math/rand"
	"sync"

	"github.com/ageapps/Peerster/data"
	"github.com/ageapps/Peerster/handler"
	"github.com/ageapps/Peerster/logger"
	"github.com/ageapps/Peerster/router"
	"github.com/ageapps/Peerster/utils"
)

// Gossiper struct
type Gossiper struct {
	peerConection   *handler.ConnectionHandler
	clientConection *handler.ConnectionHandler
	Name            string
	Address         utils.PeerAddress
	simpleMode      bool
	peers           *utils.PeerAddresses
	rumorStack      RumorStack
	privateStack    PrivateStack
	router          *router.Router
	monguerPocesses map[string]*handler.MongerHandler
	rumorCounter    *utils.Counter // [name] adress
	privateCounter  *utils.Counter // [name] adress
	mux             sync.Mutex
	usedPeers       map[string]bool
}

// NewGossiper return new instance
func NewGossiper(addressStr, name string, simple bool) (*Gossiper, error) {
	address, err1 := utils.GetPeerAddress(addressStr)
	connection, err2 := handler.NewConnectionHandler(addressStr, name)
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
		monguerPocesses: make(map[string]*handler.MongerHandler),
		rumorCounter:    utils.NewCounter(uint32(0)),
		privateCounter:  utils.NewCounter(uint32(0)),
		usedPeers:       make(map[string]bool),
	}, nil
}

// ListenToPeers function
// Start listening to Packets from peers
func (gossiper *Gossiper) ListenToPeers() {
	for {
		packet := &data.GossipPacket{}
		address, err := gossiper.peerConection.ReadPacket(packet)
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
	if connection, err := handler.NewConnectionHandler(clientAddress, gossiper.Name+"-Client"); err != nil {
		log.Fatal(err)
	} else {
		logger.Log("Listening to client in " + clientAddress)
		for {
			msg := &data.Message{}
			_, err := connection.ReadMessage(msg)
			if err != nil {
				logger.Log("Error reading message")
				break
			}
			go gossiper.HandleClientMessage(msg)
		}
	}
}

func (gossiper *Gossiper) HandleClientMessage(msg *data.Message) {
	// logger.Log("Message received from client")
	logger.LogClient(*msg)

	if gossiper.simpleMode {
		newMsg := data.NewSimpleMessage(gossiper.Name, msg.Text, gossiper.Address.String())
		gossiper.peerConection.BroadcastPacket(gossiper.peers, &data.GossipPacket{Simple: newMsg}, gossiper.Address.String())
	} else if msg.IsPrivate() {
		// Message is private
		id := gossiper.privateCounter.Increment()
		privateMessage := data.NewPrivateMessage(gossiper.Name, id, msg.Destination, msg.Text, uint32(10))
		gossiper.privateStack.AddMessage(*privateMessage)
		gossiper.sendPrivateMessage(privateMessage)
	} else {
		// Reset used peers for timers
		go gossiper.resetUsedPeers()
		id := gossiper.rumorCounter.Increment()
		rumorMessage := data.NewRumorMessage(gossiper.Name, id, msg.Text)
		gossiper.rumorStack.AddMessage(*rumorMessage)
		gossiper.mongerMessage(rumorMessage, "", false)
	}

}

func (gossiper *Gossiper) handlePeerPacket(packet *data.GossipPacket, originAddress string) {
	err := gossiper.peers.Set(originAddress)
	if err != nil {
		log.Fatal(err)
	}
	packetType := packet.GetPacketType()
	logger.Log("Peer packet received: " + packetType)

	switch packetType {
	case data.PACKET_STATUS:
		gossiper.handleStatusMessage(packet.Status, originAddress)
	case data.PACKET_RUMOR:
		gossiper.handleRumorMessage(packet.Rumor, originAddress)
	case data.PACKET_PRIVATE:
		gossiper.handlePrivateMessage(packet.Private, originAddress)
	case data.PACKET_SIMPLE:
		logger.LogSimple(*packet.Simple)
		logger.LogPeers(gossiper.peers.String())
		gossiper.handleSimpleMessage(packet.Simple, originAddress)
	default:
		logger.Log("Message not recognized")
		// log.Fatal(errors.New("Message not recognized"))
	}
}

func (gossiper *Gossiper) resetUsedPeers() {
	gossiper.mux.Lock()
	gossiper.usedPeers = make(map[string]bool)
	gossiper.mux.Unlock()
}
func (gossiper *Gossiper) getUsedPeers() map[string]bool {
	gossiper.mux.Lock()
	defer gossiper.mux.Unlock()
	return gossiper.usedPeers
}

func (gossiper *Gossiper) mongerMessage(msg *data.RumorMessage, originPeer string, routerMonguering bool) {
	gossiper.mux.Lock()
	processName := fmt.Sprint(len(gossiper.monguerPocesses), "/", routerMonguering)
	logger.Log("Starting monger process - " + processName)
	monguerProcess := handler.NewMongerHandler(originPeer, processName, routerMonguering, msg, gossiper.peerConection, gossiper.peers)
	gossiper.monguerPocesses[processName] = monguerProcess
	monguerProcess.Start()
	gossiper.mux.Unlock()
}

func (gossiper *Gossiper) findMonguerHandler(originAddress string, routeMonguer bool) *handler.MongerHandler {
	processes := gossiper.getMongerProcesses()
	// fmt.Println(processes)
	for _, process := range processes {
		if !process.IsActive() {
			gossiper.deleteMongerProcess(process.Name)
			continue
		}
		if process.GetMonguerPeer() == originAddress && process.IsRouteMonguer() == routeMonguer {
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
	logger.Log(fmt.Sprintf("Deleting monguer - %v", name))
	delete(gossiper.monguerPocesses, name)
	gossiper.mux.Unlock()
}
func (gossiper *Gossiper) getMongerProcesses() map[string]*handler.MongerHandler {
	gossiper.mux.Lock()
	defer gossiper.mux.Unlock()
	return gossiper.monguerPocesses
}
func (gossiper *Gossiper) addMongerProcess(process *handler.MongerHandler) {
	gossiper.mux.Lock()
	gossiper.monguerPocesses[process.Name] = process
	gossiper.mux.Unlock()
}
