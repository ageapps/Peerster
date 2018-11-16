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

// PeerAddresses struct
type PeerAddresses struct {
	Addresses []PeerAddress
	mux       sync.Mutex
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
		return errors.New(value + " IP was not parsed correctly")
	} else if parsedPort, err := strconv.ParseInt(ipPortStr[1], 10, 0); err != nil {
		return err
	} else {
		ad := PeerAddress{IP: parsedIP, Port: parsedPort}
		*address = ad
	}
	return nil
}

func (peers *PeerAddresses) String() string {
	peers.mux.Lock()
	defer peers.mux.Unlock()

	var s []string
	for _, peer := range peers.Addresses {
		s = append(s, peer.String())
	}
	return strings.Join(s, ",")
}

func (peers *PeerAddresses) GetAdresses() []PeerAddress {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	return peers.Addresses
}

func (peers *PeerAddresses) appendPeers(address PeerAddress) {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	peers.Addresses = append(peers.Addresses, address)
}

// Set PeerAddreses from string
func (peers *PeerAddresses) Set(value string) error {

	addresses := strings.Split(value, ",")
	for _, item := range addresses {
		var address PeerAddress
		if err := address.Set(item); err != nil {
			return err
		} else if !strings.Contains(peers.String(), item) {
			peers.appendPeers(address)
		}
	}
	return nil
}

// GetRandomPeer func
func (peers *PeerAddresses) GetRandomPeer(usedPeers map[string]bool) *PeerAddress {
	peers.mux.Lock()
	defer peers.mux.Unlock()

	peerNr := len(peers.Addresses)
	if len(usedPeers) >= peerNr {
		return nil
	}
	for {
		index := rand.Int() % peerNr
		peerAddress := peers.Addresses[index].String()
		if _, ok := usedPeers[peerAddress]; !ok {
			return &peers.Addresses[index]
		}
	}
}
