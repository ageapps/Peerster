package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/rs/cors"
)

// our main function
func main() {

	var UIPort = flag.String("port", "8080", "Port for the UI client")
	flag.Parse()

	router := NewRouter()
	log.Println("Listining on port: " + *UIPort)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Fatal(http.ListenAndServe(":"+*UIPort, handler))
}
