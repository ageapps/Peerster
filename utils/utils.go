package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
)

// Status struct
type Status struct {
	IsMongering   bool
	StatusChannel chan PeerAddress
	mux           sync.Mutex
}

// PeerAddress struct
type PeerAddress struct {
	IP   net.IP
	Port int64
}

// PeerAddreses struct
type PeerAddresses struct {
	Addresses []PeerAddress
}

// GetPeerAddress returns a PeerAdress
// from an string
func GetPeerAddress(value string) (PeerAddress, error) {
	var address PeerAddress
	return address, address.Set(value)
}

func (address *PeerAddress) String() string {
	return fmt.Sprint(address.IP.String(), ":", address.Port)
}

// Set PeerAddress from string
func (address *PeerAddress) Set(value string) error {
	ipPortStr := strings.Split(value, ":")
	if parsedIP := net.ParseIP(ipPortStr[0]); parsedIP == nil {
		return errors.New("IP was not parsed correctly")
	} else if parsedPort, err := strconv.ParseInt(ipPortStr[1], 10, 0); err != nil {
		return err
	} else {
		ad := PeerAddress{IP: parsedIP, Port: parsedPort}
		*address = ad
	}
	return nil
}

func (peers *PeerAddresses) String() string {
	var s []string
	for _, peer := range peers.Addresses {
		s = append(s, peer.String())
	}
	return strings.Join(s, ",")
}

// Set PeerAddreses from string
func (peers *PeerAddresses) Set(value string) error {
	adresses := strings.Split(value, ",")
	for _, item := range adresses {
		var adress PeerAddress
		if err := adress.Set(item); err != nil {
			return err
		} else if !strings.Contains(peers.String(), item) {
			peers.Addresses = append(peers.Addresses, adress)
		}
	}
	return nil
}

// GetRandomPeer func
func (peers *PeerAddresses) GetRandomPeer(usedPeers map[string]bool) int {
	peerNr := len(peers.Addresses)
	if len(usedPeers) >= peerNr {
		return -1
	}
	for {
		index := rand.Int() % peerNr
		peerAddress := peers.Addresses[index].String()
		if _, ok := usedPeers[peerAddress]; !ok {
			return index
		}
	}
}
