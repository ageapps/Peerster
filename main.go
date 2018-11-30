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

	var peers = utils.PeerAddresses{}
	var gossipAddr = utils.PeerAddress{IP: net.ParseIP("127.0.0.1"), Port: 5000}
	var UIPort = flag.Int("UIPort", 10000, "Define the port to which the client will connect")
	var rtimer = flag.Int("rtimer", 0, "Route rumors sending period in seconds, 0 to disable")
	var name = flag.String("name", "node", "Define the name of the gossiper")
	flag.Var(&peers, "peers", "Define the addreses of the rest of the peers to connect to separeted by a colon")
	flag.Var(&gossipAddr, "gossipAddr", "Define the ip and port to connect and send gossip messages")
	var simple = flag.Bool("simple", false, "True if using Simple messaging")

	flag.Parse()

	// if len(peers.GetAdresses()) == 0 {
	// 	flag.PrintDefaults()
	// 	os.Exit(1)
	// }

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

	if !*simple {
		if *rtimer > 0 {
			gossiper.StartRouteTimer(*rtimer)
		}
		gossiper.StartEntropyTimer()
	}
	gossiper.ListenToPeers()
}
