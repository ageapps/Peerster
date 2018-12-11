package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/ageapps/Peerster/pkg/data"
	"github.com/ageapps/Peerster/pkg/utils"
	"github.com/dedis/protobuf"
)

var (
	protocol     = "udp"
	serverAdress = utils.PeerAddress{IP: net.ParseIP("127.0.0.1")}
)

func sendMessage(msg, dest, file, index, requestHash string, budget int, keywords []string) error {
	fmt.Println("Sending <" + msg + "> to address " + serverAdress.String())
	fmt.Println("Text: " + msg)
	fmt.Println("Destination: " + file)
	fmt.Println("FileName: " + index)
	fmt.Println("IndexFilePath: " + requestHash)
	fmt.Printf("Keywords: %v\n", keywords)
	fmt.Printf("Budget: %v\n", uint64(budget))
	tmsg := &data.Message{
		Text:          msg,
		Destination:   dest,
		FileName:      file,
		IndexFilePath: index,
		RequestHash:   requestHash,
		Keywords:      keywords,
		Budget:        uint64(budget),
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
	var keywordsString = flag.String("keywords", "", "Keyboards to be searched")
	var budget = flag.Int("budget", 0, "Budget to be assingned")
	flag.Parse()
	serverAdress.Port = int64(*UIPort)
	keywords := []string{}
	if *keywordsString != "" {
		keywords = strings.Split(*keywordsString, ",")
	}
	if e := sendMessage(*msg, *dest, *file, *index, *requestHash, *budget, keywords); e != nil {
		log.Fatal(e)
	}
}
