package gossiper

import (
	"errors"
	"fmt"
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
	udpAddr, udpConn, err := startListening(address)
	if err != nil {
		return nil, err
	}
	return &ConnectionHandler{
		address: udpAddr,
		conn:    udpConn,
		Name:    name,
	}, nil
}

func (handler *ConnectionHandler) Close() {
	if err := handler.conn.Close(); err != nil {
		logger.Log(fmt.Sprintln("Error closing connection", err))
		// log.Fatal(err1)
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
func (handler *ConnectionHandler) broadcastPacket(peers *utils.PeerAddresses, packet *data.GossipPacket, incommingPeer string) {
	logger.Log("Broadcasting packet " + packet.GetPacketType())
	for _, peer := range peers.Addresses {
		if incommingPeer == peer.String() {
			continue
		}
		handler.sendPacketToPeer(peer.String(), packet)
	}
}

func (handler *ConnectionHandler) sendPacketToPeer(address string, packet *data.GossipPacket) error {
	if handler.conn == nil {
		return errors.New("No connection")
	}
	udpaddr, err1 := net.ResolveUDPAddr("udp4", address)
	if err1 != nil {
		logger.Log("Error Resolving address")
		return err1
	}
	logger.Log("Sending packet " + packet.GetPacketType() + " to <" + address + ">")
	packetBytes, err2 := protobuf.Encode(packet)
	if err2 != nil {
		logger.Log("Warning Encoding")
	}
	if _, err3 := handler.conn.WriteToUDP(packetBytes, udpaddr); err3 != nil {
		logger.Log("Error Sending Packet")
		return err3
	}
	return nil
}

func (handler *ConnectionHandler) readPacket(packet *data.GossipPacket) (string, error) {

	if handler.conn == nil {
		return "", errors.New("No connection")
	}
	buffer := make([]byte, 1024)
	_, address, err1 := handler.conn.ReadFromUDP(buffer)
	if err1 != nil {
		logger.Log("Error Reading packet")
		return "", err1
	}
	err2 := protobuf.Decode(buffer, packet)
	if err2 != nil {
		logger.Log("Warning Decoding")
	}
	return address.String(), nil
}
func (handler *ConnectionHandler) readMessage(msg *data.Message) (string, error) {
	buffer := make([]byte, 1024)
	_, address, err := handler.conn.ReadFromUDP(buffer)
	protobuf.Decode(buffer, msg)
	return address.String(), err
}
