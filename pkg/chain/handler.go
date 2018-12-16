package chain

import (
	"sync"
	"time"

	"github.com/ageapps/Peerster/pkg/data"
	"github.com/ageapps/Peerster/pkg/file"
	"github.com/ageapps/Peerster/pkg/handler"
	"github.com/ageapps/Peerster/pkg/logger"
	"github.com/ageapps/Peerster/pkg/utils"
)

// ChainHandler is a handler that will be in
// charge of requesting data from other peers
// FileName            string
// MetaHash            data.HashValue
// stopped              bool
// connection          *ConnectionHandler
// router              *router.Router
// mux                 sync.Mutex
// quitChannel         chan bool
// resetChannel        chan bool
//
type ChainHandler struct {
	blockchain     *BlockChain
	gossiperAddres string
	stopped        bool
	connection     *handler.ConnectionHandler
	Peers          *utils.PeerAddresses
	fileStore      *file.Store
	mux            sync.Mutex
	timer          *time.Timer
	quitChannel    chan bool
	BundleChannel  chan *data.TransactionBundle
	BlockChannel   chan *data.BlockBundle
}

// NewChainHandler function
func NewChainHandler(address string, peerConection *handler.ConnectionHandler, store *file.Store, peers *utils.PeerAddresses) *ChainHandler {
	return &ChainHandler{
		blockchain:    NewBlockChain(),
		stopped:       false,
		connection:    peerConection,
		Peers:         peers,
		timer:         &time.Timer{},
		quitChannel:   make(chan bool),
		BundleChannel: make(chan *data.TransactionBundle),
		BlockChannel:  make(chan *data.BlockBundle),
		fileStore:     store,
	}
}

func (handler *ChainHandler) resetTimer() {
	//logger.Log("Launching new timer")
	if handler.getTimer().C != nil {
		handler.getTimer().Stop()
	}
	handler.mux.Lock()
	handler.timer = time.NewTimer(5 * time.Second)
	handler.mux.Unlock()
}
func (handler *ChainHandler) getTimer() *time.Timer {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	return handler.timer
}

// GetStore func
func (handler *ChainHandler) GetStore() *file.Store {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	return handler.fileStore
}

// Start handler
func (handler *ChainHandler) Start(onStopHandler func()) {
	go handler.resetTimer()

	go func() {
		for {
			select {
			case bundle := <-handler.BundleChannel:
				transaction := (*bundle).Tx
				file := transaction.File
				// fmt.Println(file)
				fileHash := file.GetMetaHash()
				logger.Logf("Received transaction - %v", fileHash.String())
				fileExists := !bundle.LocalFile && handler.fileStore.FileExists(fileHash.String())
				if !fileExists && !handler.blockchain.IsTransactionSaved(transaction) {
					transaction.HopLimit--
					handler.GetStore().IndexFile(file)
					handler.addTransaction(transaction)
					if transaction.HopLimit > 0 {
						handler.publishTX(file, transaction.HopLimit, bundle.Origin)
					}
				} else {
					logger.Logf("Transaction for %v already indexed", file.Name)
				}
			case minedBlock := <-handler.blockchain.MinedBlocks:
				timer2 := time.NewTimer(time.Duration(2 * handler.blockchain.getBlockTime()))
				go func() {
					<-timer2.C
					handler.publishBlock(minedBlock, uint32(20), handler.gossiperAddres)
				}()

			case blockBundle := <-handler.BlockChannel:
				blockMsg := blockBundle.BlockPublish
				block := blockMsg.Block
				logger.Logf("Received block - %v", block.String())
				blockMsg.HopLimit--
				if handler.blockchain.CanAddBlock(&block) {
					handler.addBlock(&block)
					handler.indexTransactionsInBlock(&blockMsg.Block)
					if blockMsg.HopLimit > 0 {
						handler.publishBlock(&blockMsg.Block, blockMsg.HopLimit, blockBundle.Origin)
					}
				}
			case <-handler.quitChannel:
				logger.Log("Finishing Blockchain handler")
				if handler.timer.C != nil {
					handler.timer.Stop()
				}
				handler.blockchain.Stop()
				close(handler.BundleChannel)
				close(handler.BlockChannel)
				close(handler.quitChannel)
				onStopHandler()
				return
			}
		}
	}()
}

func (handler *ChainHandler) StartBlockchain() {
	handler.blockchain.Start(func() {
		logger.Log("Stopped Blockchain succesfully")
	})
}

func (handler *ChainHandler) indexTransactionsInBlock(block *data.Block) {
	for _, tx := range block.Transactions {
		handler.GetStore().IndexFile(tx.File)
	}
}

func (handler *ChainHandler) addTransaction(tx *data.TxPublish) {
	handler.mux.Lock()
	handler.blockchain.TxChannel <- tx
	handler.mux.Unlock()
}
func (handler *ChainHandler) addBlock(bl *data.Block) {
	handler.mux.Lock()
	handler.blockchain.BlockChannel <- bl
	handler.mux.Unlock()
}

// Stop func
func (handler *ChainHandler) Stop() {
	logger.Log("Stopping data handler")
	if !handler.stopped {
		handler.stopped = true
		return
	}
	logger.Log("Data Handler already stopped....")
}

func (handler *ChainHandler) publishTX(file file.File, hops uint32, origin string) {
	// fmt.Println(file)
	msg := data.NewTXPublish(file, hops)
	// fmt.Println(msg.File)
	packet := &data.GossipPacket{TxPublish: msg}
	handler.connection.BroadcastPacket(handler.Peers, packet, origin)
}

func (handler *ChainHandler) publishBlock(block *data.Block, hops uint32, origin string) {
	msg := data.NewBlockPublish(*block, hops)
	packet := &data.GossipPacket{BlockPublish: msg}
	// logger.Logf("%v", handler.Peers.GetAdresses())
	handler.connection.BroadcastPacket(handler.Peers, packet, origin)
}
