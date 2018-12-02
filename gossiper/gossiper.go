package gossiper

import (
	"fmt"
	"log"
	"math/rand"
	"sync"

	"github.com/ageapps/Peerster/pkg/data"
	"github.com/ageapps/Peerster/pkg/handler"
	"github.com/ageapps/Peerster/pkg/logger"
	"github.com/ageapps/Peerster/pkg/router"
	"github.com/ageapps/Peerster/pkg/utils"
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
	rumorCounter    *utils.Counter // [name] address
	privateCounter  *utils.Counter // [name] address
	mux             sync.Mutex
	usedPeers       map[string]bool
}

// NewGossiper return new instance
func NewGossiper(addressStr, name string, simple bool, rtimer int) (*Gossiper, error) {
	address, err1 := utils.GetPeerAddress(addressStr)
	connection, err2 := handler.NewConnectionHandler(addressStr, name, true)
	switch {
	case err1 != nil:
		return nil, err1
	case err2 != nil:
		return nil, err2
	}
	logger.Log(fmt.Sprint("Running simple mode: ", simple))

	gossiper := &Gossiper{
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
	}

	if !simple {
		if rtimer > 0 {
			go gossiper.StartRouteTimer(rtimer)
		}
		go gossiper.StartEntropyTimer()
	}

	return gossiper, nil
}

// ListenToPeers function
// Start listening to Packets from peers
func (gossiper *Gossiper) ListenToPeers() error {
	if gossiper.peerConection == nil {
		return fmt.Errorf("gossiper not connected to peers")
	}
	for pkt := range gossiper.peerConection.MessageQueue {
		if pkt.Address == "" {
			return fmt.Errorf("message received is not valid")
		}
		gossiper.handlePeerPacket(&pkt.Packet, pkt.Address)
	}
	return nil
}

// ListenToClients function
// start to listen for client messages
// in desired port
func (gossiper *Gossiper) ListenToClients(port int) {
	address := gossiper.Address.IP
	clientAddress := fmt.Sprintf("%v:%v", address.String(), port)
	connection, err := handler.NewConnectionHandler(clientAddress, gossiper.Name+"-Client", false)
	if err != nil {
		log.Fatal(err)
	}
	gossiper.clientConection = connection
	for pkt := range gossiper.clientConection.MessageQueue {
		if pkt.Address == "" {
			log.Fatal("message received is not valid")
		}
		gossiper.HandleClientMessage(&pkt.Message)
	}
}

// HandleClientMessage handles client messages
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
	err := gossiper.GetPeers().Set(originAddress)
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

func (gossiper *Gossiper) mongerMessage(msg *data.RumorMessage, originPeer string, routerMonguering bool) {
	gossiper.mux.Lock()
	processName := fmt.Sprint(len(gossiper.monguerPocesses), "/", routerMonguering)
	logger.Log("Starting monger process - " + processName)
	monguerProcess := handler.NewMongerHandler(originPeer, processName, routerMonguering, msg, gossiper.peerConection, gossiper.peers)
	gossiper.monguerPocesses[processName] = monguerProcess
	monguerProcess.Start()

	go func(){
		for {
			select {
			case <-monguerProcess.StopChannel:
				gossiper.deleteMongerProcess(monguerProcess.Name)
				return
			}
		}
	}()

	gossiper.mux.Unlock()
}

func (gossiper *Gossiper) findMonguerHandler(originAddress string, routeMonguer bool) *handler.MongerHandler {
	processes := gossiper.getMongerProcesses()
	for _, process := range processes {
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
