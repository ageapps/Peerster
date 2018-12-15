package chain

import (
	"encoding/hex"

	"github.com/ageapps/Peerster/pkg/data"
)

// Chain

type Chain struct {
	Blocks []*data.Block
}

func (chain *Chain) addBlock(block *data.Block) {
	chain.Blocks = append(chain.Blocks, block)
}

func (chain *Chain) isConsecutive(newBlock *data.Block) bool {
	if len(chain.Blocks) <= 0 {
		return false
	}
	hash := chain.Blocks[len(chain.Blocks)-1].Nonce
	if hex.EncodeToString(newBlock.PrevHash[:]) == hex.EncodeToString(hash[:]) {
		return true
	}
	return false
}
