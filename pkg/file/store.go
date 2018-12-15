package file

import (
	"sync"

	"github.com/ageapps/Peerster/pkg/logger"
)

// Store of files
type Store struct {
	blobs map[string]Blob
	files map[string]File
	mux   sync.Mutex
}

// NewStore func
func NewStore() *Store {
	return &Store{
		files: make(map[string]File),
		blobs: make(map[string]Blob),
	}
}

// FileExists in store
func (store *Store) FileExists(hash string) bool {
	store.mux.Lock()
	defer store.mux.Unlock()
	_, ok := store.files[hash]
	return ok
}

// GetFiles in store
func (store *Store) GetFiles() map[string]File {
	store.mux.Lock()
	defer store.mux.Unlock()
	return store.files
}

// GetBlobs in store
func (store *Store) GetBlobs() map[string]Blob {
	store.mux.Lock()
	defer store.mux.Unlock()
	return store.blobs
}

// IndexFile into store
func (store *Store) IndexFile(newFile File) (exists bool) {
	store.mux.Lock()
	logger.Logf("Indexing new file : %v", newFile.Name)
	hash := newFile.GetMetaHash()
	_, exists = store.files[hash.String()]
	if exists {
		return exists
	}
	store.files[hash.String()] = newFile
	store.mux.Unlock()
	return exists
}

// IndexBlob into store
func (store *Store) IndexBlob(newBlob Blob) (exists bool) {
	store.mux.Lock()
	logger.Logf("Indexing new blob : %v", newBlob.GetName())
	hash := newBlob.GetMetaHash()
	_, exists = store.blobs[hash.String()]
	if exists {
		return exists
	}
	store.blobs[hash.String()] = newBlob
	store.mux.Unlock()
	return exists
}

// GetBlobFromFile into store
func (store *Store) GetBlobFromFile(file File) (blob Blob, found bool) {
	hash := file.GetMetaHash()
	store.mux.Lock()
	defer store.mux.Unlock()
	blob, found = store.blobs[hash.String()]
	return
}
