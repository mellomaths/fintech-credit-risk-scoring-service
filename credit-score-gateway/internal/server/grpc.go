package server

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func StartGrpcServer(port string) {
	grpcListener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen for gRPC: %v", err)
	}
	grpcServer := grpc.NewServer()
	// Register gRPC services here
	log.Printf("gRPC server listening on %v", grpcListener.Addr())
	if err := grpcServer.Serve(grpcListener); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
