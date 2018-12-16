package chain

import (
	"encoding/hex"
	"fmt"
	"log"
	"sync"

	"github.com/ageapps/Peerster/pkg/logger"

	"github.com/ageapps/Peerster/pkg/data"
)

const (
	// NumberOfZeros in the nonce
	NumberOfZeros       = 2
	BLOCK_OLD           = "BLOCK_OLD"
	BLOCK_CURRENT       = "BLOCK_CURRENT"
	BLOCK_UNKOWN_PARENT = "BLOCK_NEW"
)

// BlockChain struct
type BlockChain struct {
	minig          bool
	CanonicalChain Chain
	SideChains     []*Chain
	currentBlock   data.Block
	prevHash       [32]byte
	tansactionPool map[string]*data.TxPublish
	blockPool      map[string]*data.Block
	mux            sync.Mutex
	quitChannel    chan bool
	stopped        bool
	MinedBlocks    chan *data.Block
	TxChannel      chan *data.TxPublish
	BlockChannel   chan *data.Block
}

// NewBlockChain func
func NewBlockChain() *BlockChain {
	return &BlockChain{
		CanonicalChain: NewEmptyChain(),
		minig:          false,
		prevHash:       [32]byte{},
		tansactionPool: make(map[string]*data.TxPublish),
		blockPool:      make(map[string]*data.Block),
		quitChannel:    make(chan bool),
		stopped:        false,
		MinedBlocks:    make(chan *data.Block, 5),
		TxChannel:      make(chan *data.TxPublish, 5),
		BlockChannel:   make(chan *data.Block, 5),
	}

}

func (bc *BlockChain) Start(onStopHandler func()) {
	for {
		select {
		case transaction := <-bc.getTxChannel():
			bc.addTransaction(transaction)
		case block := <-bc.getBlockChannel():
			bc.addBlock(block, false)
		case <-bc.quitChannel:
			bc.setMining(false)
			logger.Log("Finishing Blockchain")
			bc.CloseChannels()
			onStopHandler()
			return
		}
	}
}

func (bc *BlockChain) CanAddBlock(bl *data.Block) bool {
	return !(bc.getBlockType(bl) == BLOCK_OLD)
}

func (bc *BlockChain) mine() {
	bc.setMining(true)
	currentBlock := bc.getCurrentBlock()
	prev := bc.getPrevHash()
	logger.Logf("Expecting - %v", hex.EncodeToString(prev[:]))
	logger.Logf("Mining block with parent - %v", currentBlock.PrintPrev())

	for bc.isMining() {
		bc.resetCurrentBlock()
		nonce := currentBlock.Hash()
		currentBlock.Nonce = nonce
		if checkZeros(nonce) {
			logger.LogFoundBlock(currentBlock.String())
			bc.sendToBlockChannel(&currentBlock)
			bc.sendToMinedChannel(&currentBlock)
			break
		}
		// fmt.Printf("|")
	}
	logger.Logf("FINISHED MINIG")
	bc.setMining(false)
}

func (bc *BlockChain) CloseChannels() {
	if bc.isStopped() {
		close(bc.getMinedChannel())
		close(bc.getBlockChannel())
		close(bc.getTxChannel())
	}
}

// checkZeros
func checkZeros(nonce [32]byte) bool {
	prefix := nonce[0:NumberOfZeros]
	flag := true
	for _, num := range prefix {
		flag = flag && int(num) == 0
	}

	return flag
}

func (bc *BlockChain) IsTransactionSaved(tx *data.TxPublish) bool {
	_, ok := bc.getTransactionPool()[tx.String()]
	if ok {
		return true
	}
	return bc.isTransactionInCanonicalChain(tx)
}

func (bc *BlockChain) addTransaction(transaction *data.TxPublish) {
	// store it in mempool
	hash := bc.addToTransactionPool(transaction)
	logger.Logf("Storing - %v in TXpool", hash.String())
	if !bc.isMining() {
		bc.buildBlockAndMine()
	}
}

func (bc *BlockChain) addBlock(newBlock *data.Block, forking bool) bool {
	if !checkZeros(newBlock.Nonce) {
		return false
	}
	blockType := bc.getBlockType(newBlock)
	logger.Logf("Adding Block of type %v", blockType)
	switch blockType {
	case BLOCK_CURRENT:
		// stop mining
		bc.setMining(false)
		bc.addToBlockChain(newBlock)
		// reference the prev hash to the new added block
		bc.setPrevHash(newBlock.Nonce)
		if !bc.checkTransactionsInCurrentBlock(newBlock) {
			logger.Log("Transactions NOT found in current block")
			// there is no transactions in the block that i am currently minig,
			// jclean transaction pool from new block and keep mining
			bc.cleanTransactionPoolByAddedBlock(newBlock)
			if !forking {
				bc.buildBlockAndMine()
			}
		} else {
			logger.Log("Transactions found in current block")
			fmt.Println(newBlock.Transactions)
			fmt.Println(bc.getCurrentBlock().Transactions)
			// there is transactions in the currentBlock,
			// add them again to the Txpool, clean transaction pool from new block,
			// reset current block and rebuild it to mine again
			bc.addBlockTransactionsToPool(bc.getCurrentBlock())
			bc.cleanTransactionPoolByAddedBlock(newBlock)
			bc.resetCurrentBlock()
			if !forking {
				bc.buildBlockAndMine()
			}
		}

		return true
	case BLOCK_UNKOWN_PARENT:
		if forking {
			log.Fatal("SOMETHING IS WRONG, UNKOWN PARENT FOUND WHILE FORKING")
		}
		// just save it
		sideChainIndex, parentIndex := bc.addToBlockPool(newBlock)
		if sideChainIndex >= 0 && parentIndex >= 0 {
			logger.Log("Side chains found")
			bc.setMining(false)
			bc.forkCanonicalChain(sideChainIndex, parentIndex)
			bc.buildBlockAndMine()
		} else {
			logger.Log("NO Side chains found")
		}

		return true
	case BLOCK_OLD:
		logger.Log("Block already in blockchain, dropping it...")
		// if it's an old newBlock, just discard it
		return false
	}
	return false
}

func (bc *BlockChain) getBlockType(newBlock *data.Block) string {
	// same prev hash than out prev block
	canonicalChain := bc.getCanonicalChain()
	if canonicalChain.isNextBlockInChain(newBlock) {
		return BLOCK_CURRENT
	}
	_, ok := bc.getBlockPool()[newBlock.String()]
	if ok {
		logger.Log("Block already in pool")
		return BLOCK_OLD
	}

	if existing := bc.isBlockInCanonicalChain(newBlock); existing {
		logger.Log("Block already in chain")
		return BLOCK_OLD
	}

	return BLOCK_UNKOWN_PARENT
}

func (bc *BlockChain) cleanTransactionPoolByAddedBlock(newBlock *data.Block) {
	// Look in tx pool
	for _, newTx := range newBlock.Transactions {
		newHash := newTx.File.GetMetaHash()
		bc.deleteFromTxPool(newHash.String())
	}
}

func (bc *BlockChain) addBlockTransactionsToPool(newBlock data.Block) {
	// Look in tx pool
	for _, newTx := range newBlock.Transactions {
		bc.addToTransactionPool(&newTx)
	}
}

func (bc *BlockChain) resetCurrentBlock() {
	bc.setCurrentBlock(*data.NewBlock(bc.getPrevHash()))
}

// check if new block contains transactions that i'm currently mining
func (bc *BlockChain) checkTransactionsInCurrentBlock(newBlock *data.Block) bool {
	cb := bc.getCurrentBlock()
	// Look in current block and non coincident tx, add to tx pool
	for _, newTx := range newBlock.Transactions {
		for _, tx := range cb.Transactions {
			if newTx.String() == tx.String() {
				return true
			}
		}
	}
	return false
}

func (bc *BlockChain) deleteFromTxPool(hash string) {
	bc.mux.Lock()
	delete(bc.tansactionPool, hash)
	bc.mux.Unlock()
}

func (bc *BlockChain) buildBlockAndMine() {
	canonicalChain := bc.getCanonicalChain()
	for _, block := range bc.getBlockPool() {
		if canonicalChain.isNextBlockInChain(block) {
			logger.Logf("Found stored block matching prev - %v", block.String())
			bc.sendToBlockChannel(block)
			return
		}
	}
	if len(bc.getTransactionPool()) > 0 && !bc.isMining() {
		currentBlock := *data.NewBlock(bc.getPrevHash())
		for _, tx := range bc.getTransactionPool() {
			currentBlock.AppendTransaction(*tx)
		}
		bc.setCurrentBlock(currentBlock)
		logger.Logf("Mining current block with transactions %v", len(currentBlock.Transactions))
		go func() {
			bc.mine()
		}()
	} else {
		logger.Logf("No transactions to mine...")
	}
}

func (bc *BlockChain) addToBlockChain(block *data.Block) {
	bc.mux.Lock()
	logger.Logf("Adding Block to Canonical Chain - %v", block.String())
	logger.Logf("With prev - %v", block.PrintPrev())
	bc.CanonicalChain.appendBlock(block)
	bc.mux.Unlock()
	bc.logChain()
}

func (bc *BlockChain) addToBlockPool(block *data.Block) (sideChainIndex, parentIndex int) {
	logger.Logf("Adding Block to Blok Pool - %v", block.String())
	bc.mux.Lock()
	bc.blockPool[block.String()] = block
	bc.mux.Unlock()
	// everytime a block is added to the block pool
	// explore and build sidechains
	bc.buildSideChains(block)
	return bc.checkLongestChain()
}

// build all sidechains possible
func (bc *BlockChain) buildSideChains(block *data.Block) {
	bc.resetSideChains()
	logger.Log("Trying to build new sidechains...")
	for _, block := range bc.getBlockPool() {
		newChain := NewEmptyChain()
		lastBlock := block
		for lastBlock != nil {
			newChain.appendBlock(lastBlock)
			lastBlock = bc.findNextBlock(lastBlock)
		}
		if newChain.size() > 1 {
			logger.Logf("Found sidechain of size - %v", newChain.size())
			bc.addSideChain(&newChain)
		}
	}
}

func (bc *BlockChain) forkCanonicalChain(sideChainIndex, parentIndex int) {
	canonicalChain := bc.getCanonicalChain()
	headCanonicalChain := canonicalChain.getSubchain(0, parentIndex+1)
	removingChain := canonicalChain.getSubchain(parentIndex+1, canonicalChain.size())

	// add removed block's transactions to pool
	for _, block := range removingChain.Blocks {
		bc.addBlockTransactionsToPool(*block)
	}
	// restore canonical chain to head
	bc.restoreCanonicalChain(*headCanonicalChain)

	sideChain := bc.getSideChains()[sideChainIndex]
	// add sidechain blocks to blockchain
	for _, newBlock := range sideChain.Blocks {
		bc.addBlock(newBlock, true)
	}
}

func (bc *BlockChain) findNextBlock(block *data.Block) *data.Block {
	for _, newBlock := range bc.getBlockPool() {
		if hex.EncodeToString(newBlock.PrevHash[:]) == hex.EncodeToString(block.Nonce[:]) {
			return newBlock
		}
	}
	return nil
}

func (bc *BlockChain) checkLongestChain() (sidechain, parentBlock int) {
	sideChains := bc.getSideChains()
	canonicalChain := bc.getCanonicalChain()
	for sideChainIndex := 0; sideChainIndex < len(sideChains); sideChainIndex++ {
		sideChain := sideChains[sideChainIndex]
		chainHead := sideChain.Blocks[0]
		for blockIndex := 0; blockIndex < canonicalChain.size(); blockIndex++ {
			block := canonicalChain.Blocks[blockIndex]
			if block.IsNextBlock(chainHead) {
				logger.Log("Parent of SideChain found in CanonicalChain")
				// if a block in the canonical chain is the parent
				// of a head of a chain, lets check the sizes
				headChainSize := blockIndex + 1
				subChainSize := canonicalChain.size() - headChainSize
				// if sidechain is longer, there is a fork,
				// return the index of the restoring sidechein
				// and the heah in the blockchain
				if sideChain.size() > subChainSize {
					logger.Logf("SideChain of size %v found bigger than canonical %v", sideChain.size(), subChainSize)
					sidechain = sideChainIndex
					parentBlock = blockIndex
					return
				}
			}

		}
	}
	return -1, -1
}

// Stop func
func (bc *BlockChain) Stop() {
	logger.Log("Stopping BlockChain handler")
	if !bc.isStopped() {
		bc.setStopped(true)
		close(bc.quitChannel)
		return
	}
	logger.Log("BlockChain already stopped....")
}
