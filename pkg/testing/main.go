package main

import (
	"github.com/ageapps/Peerster/pkg/file"
	"github.com/ageapps/Peerster/pkg/logger"
)

func main() {
	logger.CreateLogger("Test", "0.0.0.0", true)
	file.NewFileAsync("test.png")
}
