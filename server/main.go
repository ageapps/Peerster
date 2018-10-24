package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
)

// our main function
func main() {
	router := NewRouter()
	port := "8080"
	log.Println("Listining on port: " + port)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Fatal(http.ListenAndServe(":"+port, handler))
}
