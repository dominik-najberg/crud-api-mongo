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

	// CREATE
	log.Println("creating a blog entry")
	req := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "Dominik Najberg",
			Title:    "This is my first blog entry",
			Content:  "To some, a thing is a solitude for absorbing. The doer experiences uniqueness which is not pictorial.",
		},
	}

	resp, err := c.CreateBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("error on Create Blog: %v", err)
	}

	log.Printf("blog item created: %s", resp.Blog)

	// READ
	log.Printf("reading blog item: %v", resp.Blog.Id)
	readReq := &blogpb.ReadBlogRequest{
		BlogId: resp.Blog.Id,
	}

	blogRes, err := c.ReadBlog(context.Background(), readReq)
	if err != nil {
		log.Fatalf("error on Read blog: %v", err)
	}

	log.Printf("blog item retrieved: %v", blogRes.Blog)
}
