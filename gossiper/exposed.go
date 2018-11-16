package gossiper

import (
	"github.com/ageapps/Peerster/data"
	"github.com/ageapps/Peerster/logger"
	"github.com/ageapps/Peerster/router"
	"github.com/ageapps/Peerster/utils"
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
	gossiper.peers = newPeers
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

// GetPeers returns an array of address strings
func (gossiper *Gossiper) GetPeers() *[]string {
	var peersArr = []string{}
	for _, peer := range gossiper.peers.GetAdresses() {
		peersArr = append(peersArr, peer.String())
	}
	return &peersArr
}

// GetRoutes returns the routing table
func (gossiper *Gossiper) GetRoutes() *router.RoutingTable {
	return gossiper.router.GetTable()
}
