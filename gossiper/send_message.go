package gossiper

import (
	"fmt"

	"github.com/ageapps/Peerster/data"
	"github.com/ageapps/Peerster/logger"
)

func (gossiper *Gossiper) sendStatusMessage(destination string) {
	var message = gossiper.rumorStack.getStatusMessage()
	packet := &data.GossipPacket{Status: message}
	gossiper.peerConection.SendPacketToPeer(destination, packet)
}

func (gossiper *Gossiper) sendRumrorMessage(destinationAdress, origin string, id uint32) {
	if message := gossiper.rumorStack.GetRumorMessage(origin, id); message != nil {
		packet := &data.GossipPacket{Rumor: message}
		logger.Log(fmt.Sprintf("Sending RUMOR ID:%v", message.ID))
		gossiper.peerConection.SendPacketToPeer(destinationAdress, packet)
	} else {
		logger.Log("Message to send not found")
	}
}

func (gossiper *Gossiper) sendRouteRumorMessage(destinationAdress string) {
	latestMsgID := gossiper.rumorCounter.GetValue() + 1
	routeRumorMessage := data.NewRumorMessage(gossiper.Name, latestMsgID, "")
	packet := &data.GossipPacket{Rumor: routeRumorMessage}
	logger.Log(fmt.Sprintf("Sending ROUTE RUMOR ID:%v", latestMsgID))
	gossiper.peerConection.SendPacketToPeer(destinationAdress, packet)
}

func (gossiper *Gossiper) sendPrivateMessage(msg *data.PrivateMessage) {
	packet := &data.GossipPacket{Private: msg}
	if destinationAdress, ok := gossiper.router.GetDestination(msg.Destination); ok {
		logger.Log(fmt.Sprintf("Sending PRIVATE Dest:%v", msg.Destination))
		gossiper.peerConection.SendPacketToPeer(destinationAdress.String(), packet)
	}
}
