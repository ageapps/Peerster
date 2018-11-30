package gossiper

import (
	"fmt"

	"github.com/ageapps/Peerster/pkg/data"
	"github.com/ageapps/Peerster/pkg/logger"
)

func (gossiper *Gossiper) handleSimpleMessage(msg *data.SimpleMessage, address string) {
	if msg.OriginalName == gossiper.Name {
		logger.Log("Received own message")
		return
	}
	newMsg := data.NewSimpleMessage(msg.OriginalName, msg.Contents, gossiper.Address.String())
	gossiper.peerConection.BroadcastPacket(gossiper.peers, &data.GossipPacket{Simple: newMsg}, msg.RelayPeerAddr)
}

func (gossiper *Gossiper) handlePrivateMessage(msg *data.PrivateMessage, address string) {
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

func (gossiper *Gossiper) handleRumorMessage(msg *data.RumorMessage, address string) {
	msgStatus := gossiper.rumorStack.CompareMessage(msg.Origin, msg.ID)
	isRouteRumor := msg.IsRouteRumor()
	routeNode := ""
	if isRouteRumor {
		logger.Log(fmt.Sprintf("Received ROUTE RUMOR"))
		routeNode = msg.Origin
	} else {
		logger.LogRumor(*msg, address)
	}
	logger.LogPeers(gossiper.peers.String())

	gossiper.router.AddIfNotExists(msg.Origin, address)

	if msgStatus == NEW_MESSAGE {
		gossiper.router.SetEntry(msg.Origin, address)

		// If I get own messages that i didn´t
		// have, set internal rumorCounter
		if msg.Origin == gossiper.Name && gossiper.rumorCounter.GetValue() > msg.ID {
			gossiper.rumorCounter.SetValue(msg.ID)
			return
		}

		if !isRouteRumor {
			// Reset used peers for timers
			go gossiper.resetUsedPeers()

			// message is new
			// -> add it to stack
			gossiper.rumorStack.AddMessage(*msg)
		}
		// -> acknowledge message
		gossiper.sendStatusMessage(address, routeNode)
		// -> start monguering message
		gossiper.mongerMessage(msg, address, isRouteRumor)
	} else if !isRouteRumor {
		// message received is not new
		// send my status msg
		gossiper.sendStatusMessage(address, "")
	}
}

func (gossiper *Gossiper) handleStatusMessage(msg *data.StatusPacket, address string) {

	isRouteStatus := msg.IsRouteStatus()
	handler := gossiper.findMonguerHandler(address, isRouteStatus)
	logger.Log(fmt.Sprint("Handler found:", handler != nil))
	logger.Log(fmt.Sprint("STATUS received Route: ", isRouteStatus))

	if isRouteStatus {
		if msg.Route != gossiper.Name {
			gossiper.router.AddIfNotExists(msg.Route, address)
		}
		if handler != nil {
			handler.Stop()
		}
		return
	}

	if handler != nil {
		handler.SetSynking(true)
	}
	if len(msg.Want) < len(*gossiper.rumorStack.getRumorStack()) {
		// check messages that i have from other peers that aren´t in the status message
		missingMessage := gossiper.rumorStack.getFirstMissingMessage(&msg.Want)
		if missingMessage != nil {
			gossiper.sendRumrorMessage(address, missingMessage.Origin, missingMessage.ID)
		}
		return
	}
	logger.LogStatus(*msg, address)
	logger.LogPeers(gossiper.peers.String())
	inSync := true

	for _, status := range msg.Want {
		messageStatus := gossiper.rumorStack.CompareMessage(status.Identifier, uint32(status.NextID-1))

		switch messageStatus {
		case NEW_MESSAGE:
			// logger.Log("Gossiper needs to update")
			gossiper.sendStatusMessage(address, "")
			break
		case IN_SYNC:
			// logger.Log("Gossiper and Peer have same messages")
		case OLD_MESSAGE:
			// logger.Log("Peer needs to update")
			gossiper.sendRumrorMessage(address, status.Identifier, status.NextID)
			break
		}
		inSync = inSync && messageStatus == IN_SYNC
	}
	if inSync {
		logger.LogInSync(address)
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
