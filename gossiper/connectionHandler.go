package gossiper

import (
	"log"
	"net"

	"github.com/ageapps/Peerster/data"
	"github.com/ageapps/Peerster/logger"
	"github.com/ageapps/Peerster/utils"
	"github.com/dedis/protobuf"
)

var (
	buffer = make([]byte, 1024)
)

type ConnectionHandler struct {
	address *net.UDPAddr
	conn    *net.UDPConn
	Name    string
}

func NewConnectionHandler(address, name string) (*ConnectionHandler, error) {
	if udpAddr, udpConn, err := startListening(address); err != nil {
		return nil, err
	} else {
		return &ConnectionHandler{
			address: udpAddr,
			conn:    udpConn,
			Name:    name,
		}, nil
	}
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
func (handler *ConnectionHandler) BroadcastPacket(peers *utils.PeerAddresses, packet *data.GossipPacket, incommingPeer string) {
	logger.Log("Broadcasting packet")
	for _, peer := range peers.Addresses {
		if incommingPeer == peer.String() {
			continue
		}
		handler.SendPacketToPeer(peer.String(), packet)
	}
}

func (handler *ConnectionHandler) SendPacketToPeer(address string, packet *data.GossipPacket) {
	udpaddr, err1 := net.ResolveUDPAddr("udp4", address)
	// fmt.Println("Sending message to " + udpaddr.String())
	if err1 != nil {
		log.Fatal(err1)
	}
	logger.Log("Sending packet to <" + address + ">")
	packetBytes, err2 := protobuf.Encode(packet)
	handler.conn.WriteToUDP(packetBytes, udpaddr)
	if err2 != nil {
		log.Fatal(err2)
	}
}

func (handler *ConnectionHandler) ReadPacket(packet *data.GossipPacket) (string, error) {
	_, address, error := handler.conn.ReadFromUDP(buffer)
	protobuf.Decode(buffer, packet)
	return address.String(), error
}
func (handler *ConnectionHandler) ReadMessage(msg *data.Message) (string, error) {
	_, address, error := handler.conn.ReadFromUDP(buffer)
	protobuf.Decode(buffer, msg)
	return address.String(), error
}
