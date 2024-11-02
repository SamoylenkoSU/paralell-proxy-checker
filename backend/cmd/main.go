package main

import (
	gapi "proxy-checker-server/internal/api/grpc"
	rapi "proxy-checker-server/internal/api/rest"
)

func main() {
	go rapi.StartServer("8081")
	gapi.StartServer("8080")
}
