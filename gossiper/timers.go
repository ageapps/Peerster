package gossiper

import (
	"fmt"
	"time"

	"github.com/ageapps/Peerster/pkg/logger"
)

// StartEntropyTimer function
func (gossiper *Gossiper) StartEntropyTimer() {
	// logger.Log("Starting Entropy timer")
	for gossiper.IsRunning() {
		usedPeers := gossiper.GetUsedPeers()

		if len(usedPeers) >= len(gossiper.GetPeers().GetAdresses()) {
			// logger.Log("Entropy Timer - All peers where notified")
		}
		if newpeer := gossiper.GetPeers().GetRandomPeer(usedPeers); newpeer != nil {
			logger.Log(fmt.Sprintf("Entropy Timer - MESSAGE to %v", newpeer.String()))
			gossiper.mux.Lock()
			gossiper.usedPeers[newpeer.String()] = true
			gossiper.mux.Unlock()
			gossiper.sendStatusMessage(newpeer.String(), "")
		}
		time.Sleep(1 * time.Second)
	}
}

// StartRouteTimer function
func (gossiper *Gossiper) StartRouteTimer(rtimer int) {
	// logger.Log("Starting Route timer")
	for gossiper.IsRunning() {
		usedPeers := make(map[string]bool)
		// if len(usedPeers) >= len(gossiper.peers.GetAdresses()) {
		// 	// logger.Log("Entropy Timer - All peers where notified")
		// }
		if newpeer := gossiper.GetPeers().GetRandomPeer(usedPeers); newpeer != nil {
			logger.Log("Route Timer - MESSAGE")
			// gossiper.mux.Lock()
			// gossiper.usedPeers[newpeer.String()] = true
			// gossiper.mux.Unlock()
			gossiper.sendRouteRumorMessage(newpeer.String())
		}
		time.Sleep(time.Duration(rtimer) * time.Second)
	}
}
