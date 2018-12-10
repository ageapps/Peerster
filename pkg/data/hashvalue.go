package data

import (
	"crypto/sha256"
	"encoding/hex"
)

// Chunk struct
type Chunk struct {
	Data []byte
	Hash HashValue
}

// Valid check
func (chunk *Chunk) Valid() bool {
	hashArr := sha256.Sum256(chunk.Data)
	var hash HashValue = hashArr[:]
	return hash.String() == chunk.Hash.String()
}

// HashValue is a file containing the SHA-256 hashes of each chunk
type HashValue []byte

// String method
func (hash *HashValue) String() string {
	return hex.EncodeToString(*hash)
}

// Set HashValue from string
func (hash *HashValue) Set(value string) error {
	newHash, err := hex.DecodeString(value)
	*hash = newHash
	return err
}

// Equals from string
func (hash *HashValue) Equals(value string) bool {
	return hash.String() == value
}

// GetHash returns a HashValue
// from an string
func GetHash(value string) (HashValue, error) {
	var hash HashValue
	return hash, hash.Set(value)
}
