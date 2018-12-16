package chain

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ageapps/Peerster/pkg/data"
	"github.com/ageapps/Peerster/pkg/utils"
)

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
	// if !yes {
	// 	logger.Logf("STOPPED MINIG")
	// }
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

func (bc *BlockChain) setCurrentNonce(nonce [32]byte) {
	bc.mux.Lock()
	bc.currentBlock.Nonce = nonce
	bc.mux.Unlock()
}
func (bc *BlockChain) setCurrentBlock(block data.Block) {
	bc.mux.Lock()
	bc.currentBlock = block
	bc.mux.Unlock()
}
func (bc *BlockChain) restoreCanonicalChain(newChain Chain) {
	bc.mux.Lock()
	bc.CanonicalChain = newChain
	bc.mux.Unlock()
}

func (bc *BlockChain) getCurrentBlock() data.Block {
	bc.mux.Lock()
	defer bc.mux.Unlock()
	return bc.currentBlock
}
func (bc *BlockChain) getCanonicalChain() Chain {
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

func (bc *BlockChain) getSideChains() []*Chain {
	bc.mux.Lock()
	defer bc.mux.Unlock()
	return bc.SideChains
}

func (bc *BlockChain) resetSideChains() {
	bc.mux.Lock()
	bc.SideChains = []*Chain{}
	bc.mux.Unlock()
}

func (bc *BlockChain) addSideChain(chain *Chain) {
	bc.mux.Lock()
	bc.SideChains = append(bc.SideChains, chain)
	bc.mux.Unlock()
}

func (bc *BlockChain) addToTransactionPool(tx *data.TxPublish) utils.HashValue {
	bc.mux.Lock()
	bc.tansactionPool[tx.String()] = tx
	bc.mux.Unlock()
	return tx.File.GetMetaHash()
}

func (bc *BlockChain) isBlockInCanonicalChain(newBlock *data.Block) bool {
	for _, block := range bc.getCanonicalChain().Blocks {
		if block.String() == newBlock.String() {
			return true
		}
	}
	return false
}

func (bc *BlockChain) deleteSideChain(index int) {
	bc.SideChains = append(bc.SideChains[:index], bc.SideChains[index+1:]...)
}

func (bc *BlockChain) isTransactionInCanonicalChain(newTransaction *data.TxPublish) bool {
	for _, block := range bc.getCanonicalChain().Blocks {
		for _, tx := range block.Transactions {
			if tx.String() == newTransaction.String() {
				return true
			}
		}
	}
	return false
}

func (bc *BlockChain) logChain() {
	str := " "
	cChain := bc.getCanonicalChain()
	for index := cChain.size() - 1; index >= 0; index-- {
		block := cChain.Blocks[index]
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
		str += " "
	}
	fmt.Println("CHAIN" + str)
}

func (bc *BlockChain) getMinedChannel() chan *data.Block {
	bc.mux.Lock()
	defer bc.mux.Unlock()
	return bc.MinedBlocks
}
func (bc *BlockChain) getBlockChannel() chan *data.Block {
	bc.mux.Lock()
	defer bc.mux.Unlock()
	return bc.BlockChannel
}
func (bc *BlockChain) getTxChannel() chan *data.TxPublish {
	bc.mux.Lock()
	defer bc.mux.Unlock()
	return bc.TxChannel
}
func (bc *BlockChain) isStopped() bool {
	bc.mux.Lock()
	defer bc.mux.Unlock()
	return bc.stopped
}
func (bc *BlockChain) setStopped(yes bool) {
	bc.mux.Lock()
	bc.stopped = yes
	bc.mux.Unlock()
}

func (bc *BlockChain) sendToBlockChannel(bl *data.Block) {
	if !bc.isStopped() {
		bc.getBlockChannel() <- (bl)
	}
}
func (bc *BlockChain) sendToMinedChannel(bl *data.Block) {
	if !bc.isStopped() {
		bc.getMinedChannel() <- bl
	}
}
func (bc *BlockChain) getBlockTime() uint64 {
	bc.mux.Lock()
	defer bc.mux.Unlock()
	return bc.BlockTime
}
func (bc *BlockChain) setBlockTime(t uint64) {
	bc.mux.Lock()
	bc.BlockTime = t
	bc.mux.Unlock()
}

func getTimestamp() int64 {
	return time.Now().UnixNano()
}
