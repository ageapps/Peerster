package main

import (
	"fmt"

	"github.com/ageapps/Peerster/pkg/logger"

	"github.com/ageapps/Peerster/pkg/file"

	"github.com/ageapps/Peerster/pkg/data"
)

func testHashValue() {
	var pepe data.HashValue
	pepe.Set("0c515910c21c81b00d899705c2da2afc70db2d0c5b29d4293f5e698fd5afa5c0")
	fmt.Println(pepe.String())
}
func testFiles() {
	logger.CreateLogger("file", "0.0.0.0", true)
	f, err := file.NewFileFromLocalSync("test.png")
	err = f.Reconstruct()
	if err != nil {
		fmt.Println(err)
	}
}

func main() {

}
