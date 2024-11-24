package main

import (
	"log"

	"github.com/bad33ndj3/commander/example/golang/server"
)

func main() {
	sv := server.New("localhost", 8080)
	err := sv.Start()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
