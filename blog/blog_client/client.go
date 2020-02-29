package main

import (
	"context"
	"github.com/dominik-najberg/crud-course/blog/blogpb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	log.Println("starting client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)
	log.Printf("client created: %#v\n", c)

	doUnary(c)
}

func doUnary(c blogpb.BlogServiceClient) {
	log.Println("doing Unary RPC")
	req := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "Dominik Najberg",
			Title:    "This is my first blog entry",
			Content:  "To some, a thing is a solitude for absorbing. The doer experiences uniqueness which is not pictorial.",
		},
	}

	resp, err := c.CreateBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("error on Greet RPC: %v", err)
	}

	log.Printf("server response: %s", resp.Blog)
}
