package main

import (
	"github.com/dominik-najberg/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

func main() {
	log.Println("server running...")

	creds, sslErr := credentials.NewServerTLSFromFile("ssl/server.crt", "ssl/server.pem")
	if sslErr != nil {
		log.Fatalf("failed while loading certificates: %v", sslErr)
	}

	s := grpc.NewServer(grpc.Creds(creds))
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	greetpb.RegisterGreetServiceServer(s, &server{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
