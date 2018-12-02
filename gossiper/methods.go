package gossiper

import (
	"fmt"

	"github.com/ageapps/Peerster/pkg/data"
	"github.com/ageapps/Peerster/pkg/logger"
	"github.com/ageapps/Peerster/pkg/router"
	"github.com/ageapps/Peerster/pkg/handler"
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
	gossiper.peerConection.Close()
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

// GetRoutes returns the routing table
func (gossiper *Gossiper) GetRoutes() *router.RoutingTable {
	return gossiper.router.GetTable()
}

func (gossiper *Gossiper) deleteMongerProcess(name string) {
	gossiper.mux.Lock()
	logger.Log(fmt.Sprintf("Deleting monguer - %v", name))
	gossiper.monguerPocesses[name]=nil // free for garbage collection
	delete(gossiper.monguerPocesses, name)
	gossiper.mux.Unlock()
}
func (gossiper *Gossiper) getMongerProcesses() map[string]*handler.MongerHandler {
	gossiper.mux.Lock()
	defer gossiper.mux.Unlock()
	return gossiper.monguerPocesses
}
func (gossiper *Gossiper) addMongerProcess(process *handler.MongerHandler) {
	gossiper.mux.Lock()
	gossiper.monguerPocesses[process.Name] = process
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
