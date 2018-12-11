package gossiper

import (
	"fmt"
	"log"

	"github.com/ageapps/Peerster/pkg/file"

	"github.com/ageapps/Peerster/pkg/data"
	"github.com/ageapps/Peerster/pkg/logger"
	"github.com/ageapps/Peerster/pkg/router"
	"github.com/ageapps/Peerster/pkg/utils"
)

// Kill func
func (gossiper *Gossiper) Kill() {
	if gossiper == nil {
		return
	}
	logger.Log("Finishing Gossiper " + gossiper.Name)
	for _, process := range gossiper.getMongerProcesses() {
		process.Stop()
	}
	for _, process := range gossiper.getDataProcesses() {
		process.Stop()
	}
	gossiper.peerConection.Close()
	gossiper.Stop()
	gossiper = nil
}

// SetPeers peers
func (gossiper *Gossiper) SetPeers(newPeers *utils.PeerAddresses) {
	gossiper.mux.Lock()
	gossiper.peers = newPeers
	gossiper.mux.Unlock()
}

// AddPeer func
func (gossiper *Gossiper) AddPeer(newPeer string) {
	gossiper.peers.Set(newPeer)
	gossiper.sendStatusMessage(newPeer, "")
}

// GetLatestMessages returns last rumor messages
func (gossiper *Gossiper) GetLatestMessages() *[]data.RumorMessage {
	return gossiper.rumorStack.GetLatestMessages()
}

// GetPrivateMessages returns last private messages
func (gossiper *Gossiper) GetPrivateMessages() *map[string][]data.PrivateMessage {
	return gossiper.privateStack.getPrivateStack()
}

// GetPeerArray returns an array of address strings
func (gossiper *Gossiper) GetPeerArray() *[]string {
	var peersArr = []string{}
	for _, peer := range gossiper.GetPeers().GetAdresses() {
		peersArr = append(peersArr, peer.String())
	}
	return &peersArr
}

// GetPeers returns current peers
func (gossiper *Gossiper) GetPeers() *utils.PeerAddresses {
	gossiper.mux.Lock()
	defer gossiper.mux.Unlock()
	return gossiper.peers
}

// IsRunning gossiper
func (gossiper *Gossiper) IsRunning() bool {
	gossiper.mux.Lock()
	defer gossiper.mux.Unlock()
	return gossiper.running
}

// Stop gossiper
func (gossiper *Gossiper) Stop() {
	gossiper.mux.Lock()
	gossiper.running = false
	gossiper.mux.Unlock()
}

// GetFiles map <metahash>:name
func (gossiper *Gossiper) GetFiles() map[string]string {
	gossiper.mux.Lock()
	defer gossiper.mux.Unlock()
	files := make(map[string]string)
	fmt.Sprintln(len(gossiper.indexedFiles))
	for k, v := range gossiper.indexedFiles {
		files[k] = v.Name
	}
	return files
}

// GetIndexedFiles func
func (gossiper *Gossiper) GetIndexedFiles() map[string]*file.File {
	gossiper.mux.Lock()
	defer gossiper.mux.Unlock()
	return gossiper.indexedFiles
}

// IndexFile into gossiper
func (gossiper *Gossiper) IndexFile(name string) {
	file, err := file.NewFileFromLocalSync(name)
	if err != nil {
		log.Fatal(err)
	}
	gossiper.addFile(file)
}

// GetRoutes returns the routing table
func (gossiper *Gossiper) GetRoutes() *router.RoutingTable {
	gossiper.mux.Lock()
	defer gossiper.mux.Unlock()

	return gossiper.router.GetTable()
}

func (gossiper *Gossiper) addFile(newFile *file.File) {
	logger.Logf("Indexing new file : %v", newFile.Name)
	gossiper.mux.Lock()
	hash := newFile.GetMetaHash()
	gossiper.indexedFiles[hash.String()] = newFile
	gossiper.mux.Unlock()
}

func (gossiper *Gossiper) resetUsedPeers() {
	gossiper.mux.Lock()
	gossiper.usedPeers = make(map[string]bool)
	gossiper.mux.Unlock()
}

// GetUsedPeers funct
func (gossiper *Gossiper) GetUsedPeers() map[string]bool {
	gossiper.mux.Lock()
	defer gossiper.mux.Unlock()
	return gossiper.usedPeers
}
