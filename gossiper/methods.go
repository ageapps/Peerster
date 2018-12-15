package gossiper

import (
	"log"

	"github.com/ageapps/Peerster/pkg/chain"
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

// AddPeers peers
func (gossiper *Gossiper) AddPeers(newPeers *utils.PeerAddresses) {
	gossiper.mux.Lock()
	gossiper.peers.AppendPeers(newPeers)
	gossiper.mux.Unlock()
}

// AddAndNotifyPeer func
func (gossiper *Gossiper) AddAndNotifyPeer(newPeer string) {
	err := gossiper.peers.Set(newPeer)
	if err != nil {
		log.Fatal(err)
	}
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

// GetFileStore map <metahash>:name
func (gossiper *Gossiper) GetFileStore() *file.Store {
	gossiper.mux.Lock()
	defer gossiper.mux.Unlock()
	return gossiper.fileStore
}

// GetFiles map <metahash>:file
func (gossiper *Gossiper) GetFiles() map[string]file.File {
	gossiper.mux.Lock()
	defer gossiper.mux.Unlock()
	return gossiper.fileStore.GetFiles()
}

// IndexAndPublishBundle func
func (gossiper *Gossiper) IndexAndPublishBundle(file *file.File, blob *file.Blob, hops uint32) {
	existed := gossiper.GetFileStore().IndexBlob(*blob)
	if !existed {
		gossiper.GetFileStore().IndexFile(*file)
		msg := data.NewTXPublish(*file, hops)
		// gossiper.chainHandler.UpdatePeers(gossiper.GetPeers())
		gossiper.chainHandler.BundleChannel <- &data.TransactionBundle{Tx: msg, Origin: gossiper.Address.String()}
	}
}

// GetChainHandler func
func (gossiper *Gossiper) GetChainHandler() *chain.ChainHandler {
	gossiper.mux.Lock()
	defer gossiper.mux.Unlock()
	return gossiper.chainHandler
}

// GetRoutes returns the routing table
func (gossiper *Gossiper) GetRoutes() *router.RoutingTable {
	gossiper.mux.Lock()
	defer gossiper.mux.Unlock()

	return gossiper.router.GetTable()
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

// SaveLocalFile func
func SaveLocalFile(name string) (*file.Blob, *file.File) {
	blob, err := file.NewBlobFromLocalSync(name)
	if err != nil {
		log.Fatal(err)
	}
	return blob, &file.File{
		Name:         blob.GetName(),
		Size:         blob.GetBlobSize(),
		MetafileHash: blob.GetMetaHash(),
	}
}
