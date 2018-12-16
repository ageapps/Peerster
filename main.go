package main

import (
	"flag"
	"log"
	"net"

	"github.com/ageapps/Peerster/gossiper"
	"github.com/ageapps/Peerster/pkg/logger"
	"github.com/ageapps/Peerster/pkg/utils"
)

// Setup flags with this sintax
// gossiper -UIPort=10000 -gossipAddr=127.0.0.1:5000 -name=nodeA -peers=127.0.0.1:5001,127.0.0.1:5002 -simple
// TO TEST
// go run main.go -UIPort=10000 -gossipAddr=127.0.0.1:5000 -name=nodeA -peers=127.0.0.1:5001 -rtimer=3
// go run main.go -UIPort=10001 -gossipAddr=127.0.0.1:5001 -name=nodeB -peers=127.0.0.1:5002 -rtimer=3
// go run main.go -UIPort=10002 -gossipAddr=127.0.0.1:5002 -name=nodeC -peers=127.0.0.1:5000 -rtimer=3

func main() {

	var peers = utils.EmptyAdresses()
	var gossipAddr = utils.PeerAddress{IP: net.ParseIP("127.0.0.1"), Port: 5000}
	var UIPort = flag.Int("UIPort", 10000, "Define the port to which the client will connect")
	var rtimer = flag.Int("rtimer", 3, "Route rumors sending period in seconds, 0 to disable")
	var name = flag.String("name", "node", "Define the name of the gossiper")
	flag.Var(peers, "peers", "Define the addreses of the rest of the peers to connect to separeted by a colon")
	flag.Var(&gossipAddr, "gossipAddr", "Define the ip and port to connect and send gossip messages")
	var simple = flag.Bool("simple", false, "True if using Simple messaging")

	flag.Parse()

	logger.CreateLogger(*name, gossipAddr.String(), false)

	var gossiper, err = gossiper.NewGossiper(gossipAddr.String(), *name, *simple, *rtimer)
	if err != nil {
		log.Fatal(err)
	}
	go gossiper.AddPeers(peers)
	go gossiper.ListenToClients(*UIPort)
	go func() {
		if err := gossiper.ListenToPeers(); err != nil {
			log.Fatal(err)
		}
	}()
	gossiper.StartBlockChain()
}
