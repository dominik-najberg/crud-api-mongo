package main

import (
	"context"
	"github.com/dominik-najberg/crud-course/blog/bootstrap"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
)

var collection *mongo.Collection

type server struct{}

func main() {
	// in case we crash we know the filename and the line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("blog service launched...")

	log.Println("connecting to MongoDB")
	client, err := bootstrap.NewClient()
	if err != nil {
		log.Fatalf("error creating MongoDB connection: %v", err)
	}

	collection = client.Database("devdb").Collection("blog")

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	//blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		log.Println("starting the server")
		if err = s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	stopChannel := make(chan os.Signal, 1)
	signal.Notify(stopChannel, os.Interrupt)
	<-stopChannel

	log.Println("stopping the server")
	s.Stop()
	log.Println("closing the listener")
	_ = lis.Close()
	log.Println("closing the MongoDB connection")
	_ = client.Disconnect(context.Background())
	log.Println("program ended")
}
