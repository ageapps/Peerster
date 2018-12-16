package gossiper

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"

	"github.com/ageapps/Peerster/pkg/file"

	"github.com/ageapps/Peerster/pkg/chain"

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
	fileProcesses   map[string]*handler.FileHandler
	searchProcesses map[string]*handler.SearchHandler
	fileStore       *file.Store
	rumorCounter    *utils.Counter // [name] address
	privateCounter  *utils.Counter // [name] address
	mux             sync.Mutex
	usedPeers       map[string]bool
	running         bool
	chainHandler    *chain.ChainHandler
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
		peers:           utils.EmptyAdresses(),
		rumorStack:      RumorStack{Messages: make(map[string][]data.RumorMessage)},
		privateStack:    PrivateStack{Messages: make(map[string][]data.PrivateMessage)},
		router:          router.NewRouter(),
		monguerPocesses: make(map[string]*handler.MongerHandler),
		fileProcesses:   make(map[string]*handler.FileHandler),
		searchProcesses: make(map[string]*handler.SearchHandler),
		fileStore:       file.NewStore(),
		rumorCounter:    utils.NewCounter(uint32(0)),
		privateCounter:  utils.NewCounter(uint32(0)),
		usedPeers:       make(map[string]bool),
		running:         true,
	}

	gossiper.chainHandler = chain.NewChainHandler(gossiper.Address.String(), gossiper.peerConection, gossiper.fileStore, gossiper.peers)

	gossiper.chainHandler.Start(func() {
		logger.Log("Chain handler stopped succesfully")
	})
	// var wg sync.WaitGroup

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()

	// }()

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

func (gossiper *Gossiper) StartBlockChain() {
	gossiper.chainHandler.StartBlockchain()
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

	logger.Logf("Message received from client \nindex: %v\nprivate: %v\nrequest: %v\nsearch: %v", msg.FileToIndex(), msg.IsDirectMessage(), msg.HasRequest(), msg.IsSearchMessage())

	if msg.FileToIndex() {
		blob, localFile := SaveLocalFile(msg.FileName)
		gossiper.IndexAndPublishBundle(localFile, blob, uint32(10))
		return
	}

	switch {
	case gossiper.simpleMode:
		logger.LogClient((*msg).Text)

		newMsg := data.NewSimpleMessage(gossiper.Name, msg.Text, gossiper.Address.String())
		gossiper.peerConection.BroadcastPacket(gossiper.peers, &data.GossipPacket{Simple: newMsg}, gossiper.Address.String())

	case msg.IsDirectMessage():
		gossiper.handleClientDirectMessage(msg)

	case msg.IsSearchMessage():
		// Message has keyboards to search
		gossiper.launchSearchProcess(msg.Keywords, msg.Budget, gossiper.Name)
		// Assign budget
	default:
		logger.LogClient((*msg).Text)
		// Reset used peers for timers
		go gossiper.resetUsedPeers()
		id := gossiper.rumorCounter.Increment()
		rumorMessage := data.NewRumorMessage(gossiper.Name, id, msg.Text)
		gossiper.rumorStack.AddMessage(*rumorMessage)
		gossiper.mongerMessage(rumorMessage, "", false)
	}
}

func (gossiper *Gossiper) handleClientDirectMessage(msg *data.Message) {
	// Message has request hash
	if msg.HasRequest() {
		hash, err := utils.GetHash(msg.RequestHash)
		if err != nil {
			log.Fatal(err)
		}
		gossiper.launchDataProcess(msg.FileName, msg.Destination, hash)
	} else {
		// Message is a private message
		logger.LogClient((*msg).Text)
		// Message is private
		id := gossiper.privateCounter.Increment()
		privateMessage := data.NewPrivateMessage(gossiper.Name, id, msg.Destination, msg.Text, uint32(10))
		gossiper.privateStack.AddMessage(*privateMessage)
		gossiper.sendPrivateMessage(privateMessage)
	}
}

func (gossiper *Gossiper) handlePeerPacket(packet *data.GossipPacket, originAddress string) {
	// if originAddress != gossiper.Address.String() {
	// 	gossiper.AddPeer(originAddress)
	// }

	err := gossiper.GetPeers().Set(originAddress)
	if err != nil {
		log.Fatal(err)
	}
	packetType := packet.GetPacketType()
	logger.Log("Received packet peer: " + packetType)
	switch packetType {
	case data.PACKET_STATUS:
		gossiper.handleStatusMessage(packet.Status, originAddress)
	case data.PACKET_RUMOR:
		gossiper.handleRumorMessage(packet.Rumor, originAddress)
	case data.PACKET_PRIVATE:
		gossiper.handlePeerPrivateMessage(packet.Private, originAddress)
	case data.PACKET_DATA_REPLY:
		gossiper.handleDataReply(packet.DataReply, originAddress)
	case data.PACKET_DATA_REQUEST:
		gossiper.handleDataRequest(packet.DataRequest, originAddress)
	case data.PACKET_SEARCH_REQUEST:
		gossiper.handleSearchRequest(packet.SearchRequest, originAddress)
	case data.PACKET_SEARCH_REPLY:
		gossiper.handleSearchReply(packet.SearchReply, originAddress)
	case data.PACKET_TX_PUBLISH:
		gossiper.handleTXMessage(packet.TxPublish, originAddress)
	case data.PACKET_BLOCK_PUBLISH:
		gossiper.handleBlockMessage(packet.BlockPublish, originAddress)
	case data.PACKET_SIMPLE:
		msg := *packet.Simple
		logger.LogSimple(msg.OriginalName, msg.RelayPeerAddr, msg.Contents)
		logger.LogPeers(gossiper.peers.String())
		gossiper.handleSimpleMessage(packet.Simple, originAddress)
	default:
		logger.Log("Message not recognized")
		// log.Fatal(errors.New("Message not recognized"))
	}
}

func (gossiper *Gossiper) mongerMessage(msg *data.RumorMessage, originPeer string, routerMonguering bool) {
	gossiper.mux.Lock()
	name := fmt.Sprint(len(gossiper.monguerPocesses), "/", routerMonguering)
	monguerProcess := handler.NewMongerHandler(originPeer, name, routerMonguering, msg, gossiper.peerConection, gossiper.peers)
	gossiper.mux.Unlock()

	gossiper.registerProcess(monguerProcess, PROCESS_MONGUER)
	monguerProcess.Start(func() {
		gossiper.unregisterProcess(monguerProcess.Name, PROCESS_MONGUER)
	})
}

func (gossiper *Gossiper) launchSearchProcess(keywords []string, budget uint64, sender string) {
	name := utils.MakeHashString(strings.Join(keywords[:], ","))
	fromClient := sender == gossiper.Name
	searchProcess := handler.NewSearchHandler(name, budget, fromClient, sender, keywords, gossiper.peerConection, gossiper.router)

	gossiper.registerProcess(searchProcess, PROCESS_SEARCH)
	searchProcess.Start(func(fileFound *data.FileResult) { // onFileReceived
		if !gossiper.GetFileStore().FileExists(fileFound.MetafileHash.String()) {
			logger.LogFoundFile(fileFound.FileName, fileFound.Destination, fileFound.MetafileHash.String(), fileFound.ChunkMap)
			gossiper.launchDataProcess(fileFound.FileName, fileFound.Destination, fileFound.MetafileHash)
		}
	}, func() { // onStopHandler
		logger.LogSearchFinished()
		gossiper.unregisterProcess(searchProcess.Name, PROCESS_SEARCH)
	})
}

func (gossiper *Gossiper) launchDataProcess(filename, destination string, metahash utils.HashValue) {
	dataProcess := handler.NewDataHandler(buildDataProcessName(destination, filename), filename, gossiper.Name, destination, metahash, gossiper.peerConection, gossiper.router)
	gossiper.registerProcess(dataProcess, PROCESS_DATA)
	dataProcess.Start(func() {
		gossiper.IndexAndPublishBundle(dataProcess.GetFile(), dataProcess.GetBlob(), uint32(10))
		gossiper.unregisterProcess(dataProcess.GetCurrentPeer(), PROCESS_DATA)
	})
}

func buildDataProcessName(peer, file string) string {
	return fmt.Sprint(peer, "/", file)
}

func keepRumorering() bool {
	// flipCoin
	coin := rand.Int() % 3
	return coin != 0
}
