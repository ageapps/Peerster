package chain

import (
	"github.com/ageapps/Peerster/pkg/data"
)

// Chain struct
type Chain struct {
	Blocks []*data.Block
}

// NewEmptyChain func
func NewEmptyChain() Chain {
	return Chain{Blocks: []*data.Block{}}
}
func (chain *Chain) appendBlock(block *data.Block) {
	chain.Blocks = append(chain.Blocks, block)
}
func (chain *Chain) size() int {
	return len(chain.Blocks)
}

func (chain *Chain) isNextBlockInChain(newBlock *data.Block) bool {
	if len(chain.Blocks) <= 0 {
		return true
	}
	lastBlock := chain.Blocks[len(chain.Blocks)-1]
	return lastBlock.IsNextBlock(newBlock)
}

func (chain *Chain) getSubchain(start, end int) *Chain {
	return &Chain{Blocks: chain.Blocks[start:end]}
}
