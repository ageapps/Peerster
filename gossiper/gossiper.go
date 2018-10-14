package gossiper

import (
	"fmt"
	"log"
	"net"

	"github.com/ageapps/Peerster/utils"
	"github.com/dedis/protobuf"
)

// Gossiper struct
type Gossiper struct {
	address *net.UDPAddr
	conn    *net.UDPConn
	Name    string
	Peers   *utils.PeerAddreses
}

// NewGossiper return new instance
func NewGossiper(address, name string) (*Gossiper, error) {
	if udpAddr, udpConn, err := startListening(address); err != nil {
		return nil, err
	} else {
		return &Gossiper{
			address: udpAddr,
			conn:    udpConn,
			Name:    name,
		}, nil
	}
}

func startListening(address string) (*net.UDPAddr, *net.UDPConn, error) {
	fmt.Println("Starting to linten in address: " + address)
	if udpAddr, err1 := net.ResolveUDPAddr("udp4", address); err1 != nil {
		return nil, nil, err1
	} else if udpConn, err2 := net.ListenUDP("udp4", udpAddr); err2 != nil {
		return nil, nil, err2
	} else {
		return udpAddr, udpConn, nil
	}
}

// BroadcastMessage function
func (gossiper *Gossiper) BroadcastMessage(text string, incommingPeer string, origin string) {
	for _, peer := range gossiper.Peers.Addresses {
		if incommingPeer == peer.String() {
			continue
		}
		udpaddr, err1 := net.ResolveUDPAddr("udp4", peer.String())
		// fmt.Println("Sending message to " + udpaddr.String())
		if err1 != nil {
			log.Fatal(err1)
		}
		sender := origin
		if sender == "" {
			sender = gossiper.Name
		}
		msg := utils.NewSimpleMessage(sender, text, gossiper.address.String())
		packet := &utils.GossipPacket{Simple: msg}
		packetBytes, err2 := protobuf.Encode(packet)
		gossiper.conn.WriteToUDP(packetBytes, udpaddr)
		if err2 != nil {
			log.Fatal(err2)
		}
	}
}

// ReceivePeerMessages function
func (gossiper *Gossiper) ReceivePeerMessages() {
	buffer := make([]byte, 1024)
	packet := &utils.GossipPacket{}
	for {
		gossiper.conn.Read(buffer)
		protobuf.Decode(buffer, packet)
		go gossiper.handleMessage(packet.Simple)
	}
}

func (gossiper *Gossiper) handleMessage(msg *utils.SimpleMessage) {
	fmt.Printf("SIMPLE MESSAGE origin %v from %v contents %v \n", msg.OriginalName, msg.RelayPeerAddr, msg.Contents)
	err := gossiper.Peers.Set(msg.RelayPeerAddr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("PEERS " + gossiper.Peers.String())
	gossiper.BroadcastMessage(msg.Contents, msg.RelayPeerAddr, msg.OriginalName)
}

//ListenToClients function
func (gossiper *Gossiper) ListenToClients(port int) {
	address := gossiper.address.IP
	fullAddress := fmt.Sprintf("%v:%v", address.String(), port)

	if udpAddr, udpConn, err := startListening(fullAddress); err != nil {
		log.Fatal(err)
	} else {
		buffer := make([]byte, 1024)
		msg := &utils.Message{}
		fmt.Println("Listening to client in " + udpAddr.String())
		for {
			_, err := udpConn.Read(buffer)
			if err != nil {
				log.Fatal(err)
			}
			protobuf.Decode(buffer, msg)
			fmt.Printf("CLIENT MESSAGE %v\n", msg.Text)
			go gossiper.BroadcastMessage(msg.Text, "", "")
		}
	}
}
