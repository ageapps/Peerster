package chain

import (
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/ageapps/Peerster/pkg/utils"

	"github.com/ageapps/Peerster/pkg/file"

	"github.com/ageapps/Peerster/pkg/logger"

	"github.com/ageapps/Peerster/pkg/data"
)

const (
	// NumberOfZeros in the nonce
	NumberOfZeros = 2
	BLOCK_OLD     = "BLOCK_OLD"
	BLOCK_CURRENT = "BLOCK_CURRENT"
	BLOCK_NEW     = "BLOCK_NEW"
)

// BlockChain struct
type BlockChain struct {
	minig          bool
	CanonicalChain []data.Block
	SideChains     [][]data.Block
	currentBlock   data.Block
	prevHash       [32]byte
	tansactionPool map[string]*data.TxPublish
	blockPool      map[string]*data.Block
	mux            sync.Mutex
	MinedBlocks    chan *data.Block
}

// NewBlockChain func
func NewBlockChain() *BlockChain {
	chain := []data.Block{}

	return &BlockChain{
		CanonicalChain: chain,
		minig:          false,
		prevHash:       [32]byte{},
		tansactionPool: make(map[string]*data.TxPublish),
		blockPool:      make(map[string]*data.Block),
		MinedBlocks:    make(chan *data.Block),
	}
}

func (bc *BlockChain) mine() {
	bc.setMining(true)
	currentBlock := bc.getCurrentBlock()
	logger.Logf("Mining block - %v", currentBlock.String())
	for bc.isMining() {
		nonce := currentBlock.Hash()
		currentBlock.Nonce = nonce
		if checkZeros(nonce) {
			logger.LogFoundBlock(hex.EncodeToString(nonce[:]))
			bc.setNonce(nonce)
			bc.setPrevHash(nonce)
			minedBlock := bc.getCurrentBlock()
			bc.MinedBlocks <- &minedBlock
			bc.addToBlockChain(minedBlock)
			break
		}
	}
	bc.setMining(false)
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

func (bc *BlockChain) isFileInTransactionPool(file *file.File) bool {
	fileMataHash := file.GetMetaHash()
	_, ok := bc.getTransactionPool()[fileMataHash.String()]
	return ok
}

func (bc *BlockChain) addTransaction(transaction *data.TxPublish) {
	// currently minig, store it in mempool
	hash := bc.addToTransactionPool(transaction)
	logger.Logf("Storing - %v in mempool", hash.String())
	if !bc.isMining() {
		bc.shouldMine()
	}
}

func (bc *BlockChain) addBlock(block *data.Block) bool {
	if !checkZeros(block.Nonce) {
		return false
	}
	blockType := bc.getBlockType(block)
	logger.Logf("Block received of type %v", blockType)
	switch blockType {
	case BLOCK_CURRENT:
		// stop mining
		bc.setMining(false)
		bc.addToBlockChain(*block)
		bc.setPrevHash(block.PrevHash)
		bc.checkTransactionPool(block)
		if bc.checkCurrentBlock(block) {
			bc.resetCurrentBlock(block)
			bc.shouldMine()
		} else {
			go bc.mine()
		}
		return true
	case BLOCK_NEW:
		if bc.checkCurrentBlock(block) {
			// new block contains tx that i'm mining
			// for now i still add it,
			// this could lead to duplicate transactions
			// but... since we're indexing files in a map,
			// trying to save a file twice is note a big deal
			bc.addToBlockPool(block)
		}
		// new block doesn´t contain tx that i'm mining
		// just save it
		bc.addToBlockPool(block)
		return true
	case BLOCK_OLD:
		// if it's an old block, just discard it
		return false
	}
	return false
}
func (bc *BlockChain) getBlockType(newBlock *data.Block) string {
	// same prev hash than out prev block
	hash := bc.getPrevHash()
	if hex.EncodeToString(newBlock.PrevHash[:]) == hex.EncodeToString(hash[:]) {
		return BLOCK_CURRENT
	}
	_, ok := bc.getBlockPool()[newBlock.String()]
	if ok {
		return BLOCK_OLD
	}

	return BLOCK_NEW
}

func (bc *BlockChain) checkTransactionPool(newBlock *data.Block) {
	// Look in tx pool
	for _, newTx := range newBlock.Transactions {
		newHash := newTx.File.GetMetaHash()
		bc.deleteFromTxPool(newHash.String())
	}
}

func (bc *BlockChain) resetCurrentBlock(newBlock *data.Block) {
	// Look in current block and non coincident tx, add to tx pool
	for _, newTx := range newBlock.Transactions {
		for _, tx := range bc.getCurrentBlock().Transactions {
			hash := tx.File.GetMetaHash()
			if newHash := newTx.File.GetMetaHash(); newHash.String() != hash.String() {
				bc.addToTransactionPool(&tx)
			}
		}
	}
}

// check if new block contains transactions that i'm currently mining
func (bc *BlockChain) checkCurrentBlock(newBlock *data.Block) bool {
	// Look in current block and non coincident tx, add to tx pool
	for _, newTx := range newBlock.Transactions {
		for _, tx := range bc.getCurrentBlock().Transactions {
			hash := tx.File.GetMetaHash()
			if newHash := newTx.File.GetMetaHash(); newHash.String() == hash.String() {
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

func (bc *BlockChain) shouldMine() {
	for _, block := range bc.getBlockPool() {
		if currentPrev := bc.getPrevHash(); hex.EncodeToString(block.PrevHash[:]) == hex.EncodeToString(currentPrev[:]) {
			logger.Logf("Found stored block matching prev - %v", block.String())
			bc.addBlock(block)
			return
		}
	}
	if len(bc.getTransactionPool()) > 0 && !bc.isMining() {
		currentBlock := *data.NewBlock(bc.getPrevHash())
		for _, tx := range bc.getTransactionPool() {
			currentBlock.AppendTransaction(*tx)
		}
		bc.setCurrentBlock(currentBlock)
		// restore mempool
		bc.resetTransactionPool()
		logger.Logf("Mining current block - %v", currentBlock.String())
		go bc.mine()
	}
}

func (bc *BlockChain) getBlockPool() map[string]*data.Block {
	bc.mux.Lock()
	defer bc.mux.Unlock()
	return bc.blockPool
}

func (bc *BlockChain) setPrevHash(newPrev [32]byte) {
	bc.mux.Lock()
	bc.prevHash = newPrev
	bc.mux.Unlock()
}

func (bc *BlockChain) setMining(yes bool) {
	bc.mux.Lock()
	bc.minig = yes
	bc.mux.Unlock()

}
func (bc *BlockChain) isMining() bool {
	bc.mux.Lock()
	defer bc.mux.Unlock()
	return bc.minig
}

func (bc *BlockChain) resetTransactionPool() {
	bc.mux.Lock()
	bc.tansactionPool = make(map[string]*data.TxPublish)
	bc.mux.Unlock()
}

func (bc *BlockChain) setNonce(nonce [32]byte) {
	bc.mux.Lock()
	bc.currentBlock.Nonce = nonce
	bc.mux.Unlock()
}
func (bc *BlockChain) setCurrentBlock(block data.Block) {
	bc.mux.Lock()
	bc.currentBlock = block
	bc.mux.Unlock()
}

func (bc *BlockChain) getCurrentBlock() data.Block {
	bc.mux.Lock()
	defer bc.mux.Unlock()
	return bc.currentBlock
}
func (bc *BlockChain) getChain() []data.Block {
	bc.mux.Lock()
	defer bc.mux.Unlock()
	return bc.CanonicalChain
}
func (bc *BlockChain) getPrevHash() [32]byte {
	bc.mux.Lock()
	defer bc.mux.Unlock()
	return bc.prevHash
}

func (bc *BlockChain) getTransactionPool() map[string]*data.TxPublish {
	bc.mux.Lock()
	defer bc.mux.Unlock()
	return bc.tansactionPool
}

func (bc *BlockChain) addToTransactionPool(tx *data.TxPublish) utils.HashValue {
	fileMataHash := tx.File.GetMetaHash()
	bc.mux.Lock()
	bc.tansactionPool[fileMataHash.String()] = tx
	bc.mux.Unlock()
	return fileMataHash
}

func (bc *BlockChain) addToBlockPool(block *data.Block) {
	bc.mux.Lock()
	bc.blockPool[block.String()] = block
	bc.mux.Unlock()
}

func (bc *BlockChain) addToBlockChain(block data.Block) {
	bc.mux.Lock()
	bc.CanonicalChain = append(bc.CanonicalChain, block)
	bc.mux.Unlock()
	bc.logChain()
}

func (bc *BlockChain) logChain() {
	str := " "
	for _, block := range bc.CanonicalChain {
		str += block.String()
		str += ":"
		str += hex.EncodeToString(block.PrevHash[:])
		str += ":"
		for index := 0; index < len(block.Transactions); index++ {
			str += block.Transactions[index].File.Name
			if index < len(block.Transactions)-1 {
				str += ","
			}
		}
	}
	fmt.Println("CHAIN" + str)
}
