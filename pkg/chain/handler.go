package chain

import (
	"fmt"
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
	peers          *utils.PeerAddresses
	store          *file.Store
	mux            sync.Mutex
	timer          *time.Timer
	quitChannel    chan bool
	BundleChannel  chan *data.Bundle
	BlockChannel   chan *data.BlockPublish
}

// NewChainHandler function
func NewChainHandler(address string, peerConection *handler.ConnectionHandler, peers *utils.PeerAddresses) *ChainHandler {
	return &ChainHandler{
		blockchain:    NewBlockChain(),
		stopped:       false,
		connection:    peerConection,
		peers:         peers,
		timer:         &time.Timer{},
		BundleChannel: make(chan *data.Bundle),
		BlockChannel:  make(chan *data.BlockPublish),
		store:         file.NewStore(),
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

// GetFileStore func
func (handler *ChainHandler) GetFileStore() *file.Store {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	return handler.store
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
				if bundle.Blob != nil {
					handler.getStore().IndexBlob(bundle.Blob)
				}
				fmt.Println(file)
				fileHash := file.GetMetaHash()
				logger.Logf("Received transaction - %v", fileHash.String())

				if !handler.store.FileExists(fileHash.String()) && !handler.blockchain.isFileInTransactionPool(&file) {
					transaction.HopLimit--
					handler.getStore().IndexFile(&transaction.File)
					handler.addTransaction(transaction)
					if transaction.HopLimit > 0 {
						handler.publishTX(file, transaction.HopLimit)
					}
				} else {
					logger.Logf("Transaction for %v already indexed", transaction.File.Name)
				}
			case minedBlock := <-handler.blockchain.MinedBlocks:
				handler.publishBlock(minedBlock, uint32(20))

			case blockMsg := <-handler.BlockChannel:
				logger.Logf("Received block - %v", blockMsg.Block.String())
				blockMsg.HopLimit--
				if handler.blockchain.addBlock(&blockMsg.Block) {
					handler.indexTransactionsInBlock(&blockMsg.Block)
					if blockMsg.HopLimit > 0 {
						handler.publishBlock(&blockMsg.Block, blockMsg.HopLimit)
					}
				}
			case <-handler.quitChannel:
				logger.Log("Finishing Blockchain handler")
				if handler.timer.C != nil {
					handler.timer.Stop()
				}
				onStopHandler()
				return
			}
		}
	}()

}

func (handler *ChainHandler) indexTransactionsInBlock(block *data.Block) {
	for _, tx := range block.Transactions {
		handler.getStore().IndexFile(&tx.File)
	}
}
func (handler *ChainHandler) getStore() *file.Store {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	return handler.store
}
func (handler *ChainHandler) addTransaction(tx *data.TxPublish) {
	handler.mux.Lock()
	handler.blockchain.addTransaction(tx)
	handler.mux.Unlock()
}

func (handler *ChainHandler) UpdatePeers(peers *utils.PeerAddresses) {
	handler.peers = peers
}

// Stop func
func (handler *ChainHandler) Stop() {
	logger.Log("Stopping data handler")
	if !handler.stopped {
		handler.stopped = true
		close(handler.BundleChannel)
		close(handler.BlockChannel)
		return
	}
	logger.Log("Data Handler already stopped....")
}

func (handler *ChainHandler) publishTX(file file.File, hops uint32) {
	fmt.Println(file)
	msg := data.NewTXPublish(file, hops)
	fmt.Println(msg.File)
	packet := &data.GossipPacket{TxPublish: msg}
	handler.connection.BroadcastPacket(handler.peers, packet, handler.gossiperAddres)
}

func (handler *ChainHandler) publishBlock(block *data.Block, hops uint32) {
	msg := data.NewBlockPublish(*block, hops)
	packet := &data.GossipPacket{BlockPublish: msg}
	logger.Logf("%v", handler.peers.GetAdresses())
	handler.connection.BroadcastPacket(handler.peers, packet, handler.gossiperAddres)
}
