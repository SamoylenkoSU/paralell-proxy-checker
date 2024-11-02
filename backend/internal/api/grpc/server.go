package grpc

import (
	"log"
	"net"
	pb "proxy-checker-server/generated/grpc/proxy-checker.api"

	"google.golang.org/grpc"
)

func StartServer(port string) {
	log.Println("Starting grpc server")

	listner, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Listening on %s", port)
	server := grpc.NewServer()

	pb.RegisterProxyCheckerServer(server, &ApiServer{})

	if err := server.Serve(listner); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
