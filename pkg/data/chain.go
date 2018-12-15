package data

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"

	"github.com/ageapps/Peerster/pkg/file"
)

// TxPublish struct
type TxPublish struct {
	File     file.File
	HopLimit uint32
}

// TransactionBundle struct
type TransactionBundle struct {
	Tx     *TxPublish
	Blob   *file.Blob
	Origin string
}

// BlockBundle struct
type BlockBundle struct {
	BlockPublish *BlockPublish
	Origin       string
}

// BlockPublish struct
type BlockPublish struct {
	Block    Block
	HopLimit uint32
}

// Block stuct
type Block struct {
	PrevHash     [32]byte
	Nonce        [32]byte
	Transactions []TxPublish
}

// NewBlock func
func NewBlock(prev [32]byte) *Block {
	return &Block{
		PrevHash:     prev,
		Transactions: []TxPublish{},
	}
}

// AppendTransaction func
func (block *Block) AppendTransaction(tx TxPublish) {
	block.Transactions = append(block.Transactions, tx)
}

// AppendTransaction func
func (block *Block) String() string {
	return hex.EncodeToString(block.Nonce[:])
}

// NewTXPublish func
func NewTXPublish(file file.File, hops uint32) *TxPublish {
	return &TxPublish{file, hops}
}

// NewBlockPublish func
func NewBlockPublish(block Block, hops uint32) *BlockPublish {
	return &BlockPublish{block, hops}
}

func (b *Block) Hash() (out [32]byte) {
	h := sha256.New()
	h.Write(b.PrevHash[:])
	h.Write(b.Nonce[:])
	binary.Write(h, binary.LittleEndian, uint32(len(b.Transactions)))
	for _, t := range b.Transactions {
		th := t.Hash()
		h.Write(th[:])
	}
	copy(out[:], h.Sum(nil))
	return
}

func (t *TxPublish) Hash() (out [32]byte) {
	h := sha256.New()
	binary.Write(h, binary.LittleEndian,
		uint32(len(t.File.Name)))
	h.Write([]byte(t.File.Name))
	h.Write(t.File.MetafileHash)
	copy(out[:], h.Sum(nil))
	return
}
