package gossiper

import (
	"log"
	"net"

	"github.com/ageapps/Peerster/data"
	"github.com/ageapps/Peerster/logger"
	"github.com/ageapps/Peerster/utils"
	"github.com/dedis/protobuf"
)
// ConnectionHandler handles 
// all connections coding and decoding all packets
type ConnectionHandler struct {
	address *net.UDPAddr
	conn    *net.UDPConn
	Name    string
}

// NewConnectionHandler function
func NewConnectionHandler(address, name string) (*ConnectionHandler, error) {
	udpAddr, udpConn, err := startListening(address); 
	if err != nil {
		return nil, err
	}
	return &ConnectionHandler{
		address: udpAddr,
		conn:    udpConn,
		Name:    name,
	}, nil
}

func startListening(address string) (*net.UDPAddr, *net.UDPConn, error) {
	logger.Log("Starting to linten in address: " + address)
	if udpAddr, err1 := net.ResolveUDPAddr("udp4", address); err1 != nil {
		return nil, nil, err1
	} else if udpConn, err2 := net.ListenUDP("udp4", udpAddr); err2 != nil {
		return nil, nil, err2
	} else {
		return udpAddr, udpConn, nil
	}
}

// BroadcastPacket function
func (handler *ConnectionHandler) broadcastPacket(peers *utils.PeerAddresses, packet *data.GossipPacket, incommingPeer string) {
	logger.Log("Broadcasting packet " + packet.GetPacketType())
	for _, peer := range peers.Addresses {
		if incommingPeer == peer.String() {
			continue
		}
		handler.sendPacketToPeer(peer.String(), packet)
	}
}

func (handler *ConnectionHandler) sendPacketToPeer(address string, packet *data.GossipPacket) {
	udpaddr, err1 := net.ResolveUDPAddr("udp4", address)
	// fmt.Println("Sending message to " + udpaddr.String())
	if err1 != nil {
		log.Fatal(err1)
	}
	logger.Log("Sending packet " + packet.GetPacketType() + " to <" + address + ">")
	// fmt.Println(packet)
	packetBytes, err2 := protobuf.Encode(packet)
	if err2 != nil {
		logger.Log("Error Encoding")
		log.Fatal(err2)
	}
	if _, err3 := handler.conn.WriteToUDP(packetBytes, udpaddr); err3 != nil {
		logger.Log("Error Sending Packet")
		log.Fatal(err3)
	}
}

func (handler *ConnectionHandler) readPacket(packet *data.GossipPacket) string {
	buffer := make([]byte, 1024)
	_, address, err1 := handler.conn.ReadFromUDP(buffer)
	if err1 != nil {
		logger.Log("Error Reading packet")
		log.Fatal(err1)
	}
	err2 := protobuf.Decode(buffer, packet)
	if err2 != nil {
		// logger.Log("Error Decoding")
		// log.Fatal(err2)
	}

	return address.String()
}
func (handler *ConnectionHandler) readMessage(msg *data.Message) (string, error) {
	buffer := make([]byte, 1024)
	_, address, error := handler.conn.ReadFromUDP(buffer)
	protobuf.Decode(buffer, msg)
	return address.String(), error
}
