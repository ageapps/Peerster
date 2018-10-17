package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/ageapps/Peerster/data"
	"github.com/ageapps/Peerster/utils"
	"github.com/dedis/protobuf"
)

var (
	protocol     = "udp"
	serverAdress = utils.PeerAddress{IP: net.ParseIP("127.0.0.1")}
)

func sendMessage(msg string) error {
	fmt.Println("Sending <" + msg + "> to address " + serverAdress.String())
	tmsg := &data.Message{
		Text: msg,
	}
	buf, err1 := protobuf.Encode(tmsg)
	conn, err2 := net.Dial(protocol, serverAdress.String())
	defer conn.Close()
	switch {
	case err1 != nil:
		return err1
	case err2 != nil:
		return err2
	}
	conn.Write(buf)
	return nil
}

func main() {

	// Setup flags with this sintax
	// ./client -UIPort=10000 -msg=Hello
	var UIPort = flag.Int("UIPort", 10000, "Port for the UI client")
	var msg = flag.String("msg", "", "Message to be sent")
	flag.Parse()
	serverAdress.Port = int64(*UIPort)
	if e := sendMessage(*msg); e != nil {
		log.Fatal(e)
	}

}
