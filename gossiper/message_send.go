package gossiper

import (
	"fmt"

	"github.com/ageapps/Peerster/pkg/data"
	"github.com/ageapps/Peerster/pkg/logger"
)

func (gossiper *Gossiper) sendStatusMessage(destination, nodeName string) {
	var message *data.StatusPacket
	if nodeName != "" {
		message = data.NewStatusPacket(nil, nodeName)
	} else {
		message = gossiper.rumorStack.getStatusMessage()
	}
	logger.Log(fmt.Sprint("Sending STATUS route: ", nodeName != ""))
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

func (gossiper *Gossiper) broadcastRouteRumorMessage(destinationAdress string) {
	latestMsgID := gossiper.rumorCounter.GetValue() + 1
	routeRumorMessage := data.NewRumorMessage(gossiper.Name, latestMsgID, "")
	packet := &data.GossipPacket{Rumor: routeRumorMessage}
	logger.Log(fmt.Sprintf("Sending ROUTE RUMOR ID:%v", latestMsgID))
	gossiper.peerConection.SendPacketToPeer(destinationAdress, packet)
}

func (gossiper *Gossiper) sendPrivateMessage(msg *data.PrivateMessage) {
	packet := &data.GossipPacket{Private: msg}
	if destinationAdress, ok := gossiper.router.GetDestination(msg.Destination); ok {
		logger.Logf("Sending PRIVATE Dest:%v", msg.Destination)
		gossiper.peerConection.SendPacketToPeer(destinationAdress.String(), packet)
	} else {
		logger.Logf("INVALID PRIVATE Dest:%v", msg.Destination)
	}
}

func (gossiper *Gossiper) sendDataRequest(msg *data.DataRequest) {
	packet := &data.GossipPacket{DataRequest: msg}
	if destinationAdress, ok := gossiper.router.GetDestination(msg.Destination); ok {
		logger.Log(fmt.Sprintf("Sending DATA REQUEST Dest:%v", msg.Destination))
		gossiper.peerConection.SendPacketToPeer(destinationAdress.String(), packet)
	} else {
		logger.Logf("INVALID DATA REQUEST Dest:%v", msg.Destination)
	}
}
func (gossiper *Gossiper) sendDataReply(msg *data.DataReply) {
	packet := &data.GossipPacket{DataReply: msg}
	if destinationAdress, ok := gossiper.router.GetDestination(msg.Destination); ok {
		logger.Log(fmt.Sprintf("Sending DATA REPLY Dest:%v", msg.Destination))
		gossiper.peerConection.SendPacketToPeer(destinationAdress.String(), packet)
	} else {
		logger.Logf("INVALID DATA REPLY Dest:%v", msg.Destination)
	}
}
func (gossiper *Gossiper) sendSearchReply(msg *data.SearchReply) {
	packet := &data.GossipPacket{SearchReply: msg}
	if destinationAdress, ok := gossiper.router.GetDestination(msg.Destination); ok {
		logger.Log(fmt.Sprintf("Sending SEARCH REPLY Dest:%v", msg.Destination))
		gossiper.peerConection.SendPacketToPeer(destinationAdress.String(), packet)
	} else {
		logger.Logf("INVALID SEARCH REPLY Dest:%v", msg.Destination)
	}
}
