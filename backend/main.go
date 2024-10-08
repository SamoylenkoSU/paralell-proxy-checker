package main

import (
	"log"
	"net"
	"proxy-checker-server/api/handler"

	"google.golang.org/grpc"
)

const (
	port = ":8080"
)

func main() {
	log.Println("Starting server")

	listner, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Listening on %s", port)
	server := grpc.NewServer()

	handler.ServerRegister(server)

	if err := server.Serve(listner); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
