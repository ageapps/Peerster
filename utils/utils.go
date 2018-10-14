package utils

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// GossipAddress struct
type GossipAddress struct {
	IP   net.IP
	Port int64
}

// PeerAddreses struct
type PeerAddreses struct {
	Addresses []GossipAddress
}

// SimpleMessage struct
type SimpleMessage struct {
	OriginalName  string
	RelayPeerAddr string
	Contents      string
}

// GossipPacket struct
type GossipPacket struct {
	Simple *SimpleMessage
}

// Message to send
type Message struct {
	Text string
}

// NewSimpleMessage create
func NewSimpleMessage(ogname, msg, relay string) *SimpleMessage {
	return &SimpleMessage{
		OriginalName:  ogname,
		RelayPeerAddr: relay,
		Contents:      msg,
	}
}

func (address *GossipAddress) String() string {
	return fmt.Sprint(address.IP.String(), ":", address.Port)
}

// Set GossipAddress from string
func (address *GossipAddress) Set(value string) error {
	ipPortStr := strings.Split(value, ":")
	if parsedIP := net.ParseIP(ipPortStr[0]); parsedIP == nil {
		return errors.New("IP was not parsed correctly")
	} else if parsedPort, err := strconv.ParseInt(ipPortStr[1], 10, 0); err != nil {
		return err
	} else {
		ad := GossipAddress{IP: parsedIP, Port: parsedPort}
		*address = ad
	}
	return nil
}

func (peers *PeerAddreses) String() string {
	var s []string
	for _, peer := range peers.Addresses {
		s = append(s, peer.String())
	}
	return strings.Join(s, ",")
}

// Set PeerAddreses from string
func (peers *PeerAddreses) Set(value string) error {
	adresses := strings.Split(value, ",")
	for _, item := range adresses {
		var adress GossipAddress
		if err := adress.Set(item); err != nil {
			return err
		} else if !strings.Contains(peers.String(), item) {
			peers.Addresses = append(peers.Addresses, adress)
		}
	}
	return nil
}
