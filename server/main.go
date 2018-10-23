package main

import (
	"log"
	"net/http"
)

// our main function
func main() {
	router := NewRouter()
	port := "8080"
	log.Println("Listining on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
