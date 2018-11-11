package gossiper
import (
	"time"

	"github.com/ageapps/Peerster/logger"
)
// StartEntropyTimer function
func (gossiper *Gossiper) StartEntropyTimer() {
	go func() {
		logger.Log("Starting Entropy timer")
		for {
			if len(gossiper.getUsedPeers()) >= len(gossiper.peers.GetAdresses()) {
				// logger.Log("Entropy Timer - All peers where notified")
			}
			if newpeer := gossiper.peers.GetRandomPeer(gossiper.usedPeers); newpeer != nil {
				logger.Log("Entropy Timer - MESSAGE")
				gossiper.usedPeers[newpeer.String()] = true
				gossiper.sendStatusMessage(newpeer.String())
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

// StartEntropyTimer function
func (gossiper *Gossiper) StartRouteTimer(rtimer int) {
	var usedPeers = make(map[string]bool)

	go func() {
		logger.Log("Starting Route timer")
		for {
			if len(usedPeers) >= len(gossiper.peers.GetAdresses()) {
				usedPeers = make(map[string]bool)
				// logger.Log("Entropy Timer - All peers where notified")
			}
			if newpeer := gossiper.peers.GetRandomPeer(usedPeers); newpeer != nil {
				logger.Log("Route Timer - MESSAGE")
				gossiper.usedPeers[newpeer.String()] = true
				gossiper.sendRouteRumorMessage(newpeer.String())
			}
			time.Sleep(time.Duration(rtimer) * time.Second)
		}
	}()
}

