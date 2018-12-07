package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/ageapps/Peerster/pkg/data"
	"github.com/ageapps/Peerster/pkg/utils"
	"github.com/dedis/protobuf"
)

var (
	protocol     = "udp"
	serverAdress = utils.PeerAddress{IP: net.ParseIP("127.0.0.1")}
)

func sendMessage(msg, dest string) error {
	fmt.Println("Sending <" + msg + "> to address " + serverAdress.String())
	tmsg := &data.Message{
		Text:        msg,
		Destination: dest,
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
	// go run . -UIPort=10000 -msg=Hello
	var UIPort = flag.Int("UIPort", 10000, "Port for the UI client")
	var dest = flag.String("dest", "", "Destination for the private message")
	var msg = flag.String("msg", "", "file to be indexed")
	var file = flag.String("file", "", "Message to be sent")
	flag.Parse()
	serverAdress.Port = int64(*UIPort)
	if e := sendMessage(*msg, *dest); e != nil {
		log.Fatal(e)
	}
	_ = file

}
