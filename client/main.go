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

func sendMessage(msg, dest, file, index, requestHash string) error {
	fmt.Println("Sending <" + msg + "> to address " + serverAdress.String())
	tmsg := &data.Message{
		Text:          msg,
		Destination:   dest,
		FileName:      file,
		IndexFilePath: index,
		RequestHash:   requestHash,
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
	var dest = flag.String("Dest", "", "Destination for the private message")
	var msg = flag.String("msg", "", "Message to be sent")
	var file = flag.String("file", "", "Name of file requested")
	var index = flag.String("index", "", "File to be indexed")
	var requestHash = flag.String("request", "", "HashValue to be requested")
	flag.Parse()
	serverAdress.Port = int64(*UIPort)
	if e := sendMessage(*msg, *dest, *file, *index, *requestHash); e != nil {
		log.Fatal(e)
	}
}
