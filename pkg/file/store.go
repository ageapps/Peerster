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

// BlobExists in store
func (store *Store) BlobExists(hash string) bool {
	store.mux.Lock()
	defer store.mux.Unlock()
	_, ok := store.blobs[hash]
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
	logger.Logf("Indexing new file : %v", newFile.Name)
	hash := newFile.GetMetaHash()
	_, exists = store.files[hash.String()]
	if exists {
		logger.Logf("Already exists")
		return true
	}
	store.mux.Lock()
	store.files[hash.String()] = newFile
	store.mux.Unlock()
	return false
}

// IndexBlob into store
func (store *Store) IndexBlob(newBlob Blob) bool {
	logger.Logf("Indexing new blob : %v", newBlob.fileName)
	hash := newBlob.GetMetaHash()
	_, exists := store.blobs[hash.String()]
	if exists {
		logger.Logf("Already exists")
		return true
	}
	store.mux.Lock()
	store.blobs[hash.String()] = newBlob
	store.mux.Unlock()
	return false
}

// GetBlobFromFile into store
func (store *Store) GetBlobFromFile(file File) (blob Blob, found bool) {
	hash := file.GetMetaHash()
	store.mux.Lock()
	defer store.mux.Unlock()
	blob, found = store.blobs[hash.String()]
	return
}
