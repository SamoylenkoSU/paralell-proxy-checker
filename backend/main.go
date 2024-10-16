package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"proxy-checker-server/api/handler"

	"google.golang.org/grpc"
)

const (
	port = ":8080"
)

func main() {
	go func() {

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
		})

		log.Printf("Listening http on %s", "8081")
		http.ListenAndServe(":8081", nil)
	}()

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
