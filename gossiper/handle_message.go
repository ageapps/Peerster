package gossiper

import (
	"fmt"

	"github.com/ageapps/Peerster/data"
	"github.com/ageapps/Peerster/logger"
)

func (gossiper *Gossiper) handleSimpleMessage(msg *data.SimpleMessage, from string) {
	if msg.OriginalName == gossiper.Name {
		logger.Log("Received own message")
		return
	}
	newMsg := data.NewSimpleMessage(msg.OriginalName, msg.Contents, gossiper.Address.String())
	gossiper.peerConection.BroadcastPacket(gossiper.peers, &data.GossipPacket{Simple: newMsg}, msg.RelayPeerAddr)
}

func (gossiper *Gossiper) handlePrivateMessage(msg *data.PrivateMessage, from string) {
	if msg.Destination == gossiper.Name {
		gossiper.privateStack.AddMessage(*msg)
		logger.LogPrivate(*msg)
		return
	}
	msg.HopLimit--
	if msg.HopLimit > 0 {
		gossiper.sendPrivateMessage(msg)
	}
}

func (gossiper *Gossiper) handleRumorMessage(msg *data.RumorMessage, from string) {
	isRouteRumor := msg.IsRouteRumor()
	if !isRouteRumor {
		logger.LogRumor(*msg, from)
		logger.LogPeers(gossiper.peers.String())
	}
	msgStatus := gossiper.rumorStack.CompareMessage(msg.Origin, msg.ID)

	if msgStatus == NEW_MESSAGE {
		newEntry := gossiper.router.SetEntry(msg.Origin, from)
		if isRouteRumor {
			logger.Log(fmt.Sprintf("Received ROUTE RUMOR - new:%v", newEntry))
			gossiper.mongerMessage(msg, from)
			return
		}
		// If I get own messages that i didn´t
		// have, set internal rumorCounter
		if msg.Origin == gossiper.Name {
			gossiper.rumorCounter.SetValue(msg.ID)
		}
		// logger.Log("Received new rumor message, appending...")

		// Reset used peers for timers
		go gossiper.resetUsedPeers()

		// message is new
		// -> add it to stack
		gossiper.rumorStack.AddMessage(*msg)
		// -> acknowledge message
		gossiper.sendStatusMessage(from)
		// -> start monguering message
		gossiper.mongerMessage(msg, from)
	} else {
		// message received is not new
		// send my status msg
		gossiper.sendStatusMessage(from)
	}
}

func (gossiper *Gossiper) handleStatusMessage(msg *data.StatusPacket, from string) {
	handler := gossiper.findMonguerHandler(from)
	logger.Log(fmt.Sprint("Handler found:", handler != nil))

	if len(msg.Want) < len(*gossiper.rumorStack.getRumorStack()) {

		// check messages that i have from other peers that aren´t in the status message
		for origin, messages := range *gossiper.rumorStack.getRumorStack() {
			firstMessageID := messages[0].ID
			found := false
			for _, status := range msg.Want {
				if status.Identifier == origin {
					found = true
					break
				}
			}
			if !found {
				if handler != nil {
					handler.SetSynking(true)
				}
				logger.Log(fmt.Sprintf("Peer needs to update Origin:%v - ID:%v", origin, firstMessageID))
				gossiper.sendRumrorMessage(from, origin, firstMessageID)
				return
			}
		}
	}
	logger.LogStatus(*msg, from)
	logger.LogPeers(gossiper.peers.String())
	inSync := true

	for _, status := range msg.Want {
		messageStatus := gossiper.rumorStack.CompareMessage(status.Identifier, uint32(status.NextID-1))

		switch messageStatus {
		case NEW_MESSAGE:
			// logger.Log("Gossiper needs to update")
			if handler != nil {
				handler.SetSynking(true)
			}
			gossiper.sendStatusMessage(from)
			break
		case IN_SYNC:
			// logger.Log("Gossiper and Peer have same messages")
		case OLD_MESSAGE:
			// logger.Log("Peer needs to update")
			if handler != nil {
				handler.SetSynking(true)
			}
			gossiper.sendRumrorMessage(from, status.Identifier, status.NextID)
			break
		}
		inSync = inSync && messageStatus == IN_SYNC
	}
	if inSync {
		logger.LogInSync(from)
		if handler != nil {
			handler.SetSynking(false)
			// Flip coin
			logger.Log("IN SYNC, FLIPPING COIN")
			if !keepRumorering() {
				handler.Stop()
				// delete handler from slice
				go gossiper.deleteMongerProcess(handler.Name)
			} else {
				handler.Reset()
			}
		}
	}
}
