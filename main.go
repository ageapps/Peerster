package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/ageapps/Peerster/gossiper"
	"github.com/ageapps/Peerster/logger"
	"github.com/ageapps/Peerster/utils"
)

func main() {

	// Setup flags with this sintax
	// gossiper -UIPort=10000 -gossipAddr=127.0.0.1:5000 -name=nodeA -peers=127.0.0.1:5001,127.0.0.1:5002 -simple
	// TO TEST
	// go run main.go -UIPort=10000 -gossipAddr=127.0.0.1:5000 -name=nodeA -peers=127.0.0.1:5001
	// go run main.go -UIPort=10001 -gossipAddr=127.0.0.1:5001 -name=nodeB -peers=127.0.0.1:5002
	// go run main.go -UIPort=10002 -gossipAddr=127.0.0.1:5002 -name=nodeC -peers=127.0.0.1:5000 
	var peers = utils.PeerAddresses{}
	var gossipAddr = utils.PeerAddress{IP: net.ParseIP("127.0.0.1"), Port: 5000}
	var UIPort = flag.Int("UIPort", 10000, "Define the port to which the client will connect")
	var name = flag.String("name", "node", "Define the name of the gossiper")
	flag.Var(&peers, "peers", "Define the addreses of the rest of the peers to connect to separeted by a colon")
	flag.Var(&gossipAddr, "gossipAddr", "Define the ip and port to connect and send gossip messages")
	var simple = flag.Bool("simple", false, "True if using Simple messaging")

	flag.Parse()

	if len(peers.GetAdresses()) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// fmt.Println(UIPort)
	// fmt.Println(name)
	// fmt.Println(simple)
	// fmt.Println(peers.String())

	logger.CreateLogger(*name, gossipAddr.String(), true)

	var gossiper, err = gossiper.NewGossiper(gossipAddr.String(), *name, *simple)
	if err != nil {
		log.Fatal(err)
	}
	gossiper.SetPeers(&peers)
	go gossiper.ListenToClients(*UIPort)
	go gossiper.ListenToPeers()
	if !*simple {
		gossiper.StartEntropyTimer()
	}
	for {

	}
}
