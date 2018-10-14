package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/ageapps/Peerster/gossiper"
	"github.com/ageapps/Peerster/utils"
)

func main() {

	// Setup flags with this sintax
	// gossiper -UIPort=10000 -gossipAddr=127.0.0.1:5000 -name=nodeA -peers=127.0.0.1:5001,10.1.1.7:5002 -simple
	var peers = utils.PeerAddreses{}
	var gossipAddr = utils.GossipAddress{IP: net.ParseIP("127.0.0.1"), Port: 5000}
	var UIPort = flag.Int("UIPort", 10000, "Define the port to which the client will connect")
	var name = flag.String("name", "node", "Define the name of the gossiper")
	flag.Var(&peers, "peers", "Define the addreses of the rest of the peers to connect to separeted by a colon")
	flag.Var(&gossipAddr, "gossipAddr", "Define the ip and port to connect and send gossip messages")
	var simple = flag.Bool("simple", true, "True if using Simple messaging")

	if *simple {

	}

	flag.Parse()

	if len(peers.Addresses) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// fmt.Println(gossipAddr.String())
	// fmt.Println(peers.String())

	var gossiper, err = gossiper.NewGossiper(gossipAddr.String(), *name)
	if err != nil {
		log.Fatal(err)
	}
	gossiper.Peers = &peers
	go gossiper.ListenToClients(*UIPort)
	go gossiper.ReceivePeerMessages()
	for {

	}
}
