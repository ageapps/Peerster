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

func (gossiper *Gossiper) GetLatestMessages() *[]data.RumorMessage {
	return gossiper.rumorStack.GetLatestMessages()
}
func (gossiper *Gossiper) GetPrivateMessages() *map[string][]data.PrivateMessage {
	return gossiper.privateStack.getPrivateStack()
}

func (gossiper *Gossiper) GetPeers() *[]string {
	var peersArr = []string{}
	for _, peer := range gossiper.peers.GetAdresses() {
		peersArr = append(peersArr, peer.String())
	}
	return &peersArr
}

func (gossiper *Gossiper) GetRoutes() *router.RoutingTable {
	return gossiper.router.GetTable()
}
