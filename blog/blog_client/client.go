package main

import (
	"context"
	"github.com/dominik-najberg/crud-course/blog/blogpb"
	"google.golang.org/grpc"
	"io"
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
	createReq := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "Dominik Najberg",
			Title:    "This is my first blog entry",
			Content:  "To some, a thing is a solitude for absorbing. The doer experiences uniqueness which is not pictorial.",
		},
	}

	resp, err := c.CreateBlog(context.Background(), createReq)
	if err != nil {
		log.Fatalf("error on Create Blog: %v", err)
	}

	log.Printf("blog item created: %s", resp.Blog)

	// READ
	log.Printf("reading blog item: %v", resp.Blog.Id)
	readReq := &blogpb.ReadBlogRequest{
		BlogId: resp.Blog.GetId(),
	}

	blogRes, err := c.ReadBlog(context.Background(), readReq)
	if err != nil {
		log.Fatalf("error on Read blog: %v", err)
	}

	log.Printf("blog item retrieved: %v", blogRes.Blog)

	// UPDATE
	log.Printf("updating blog item: %v", resp.Blog.Id)

	updateReq := &blogpb.UpdateBlogRequest{
		Blog: &blogpb.Blog{
			Id:       resp.Blog.GetId(),
			AuthorId: createReq.GetBlog().GetAuthorId(),
			Title:    "(updated) " + createReq.Blog.GetTitle(),
			Content:  createReq.GetBlog().GetContent(),
		},
	}
	updateRes, err := c.UpdateBlog(context.Background(), updateReq)
	if err != nil {
		log.Fatalf("update error: %v", err)
	}
	log.Printf("blog item updated: %v", updateRes)

	//// DELETE
	//log.Printf("deleting blog item: %v", resp.Blog.Id)
	//
	//deleteReq := &blogpb.DeleteBlogRequest{
	//	BlogId: resp.Blog.Id,
	//}
	//
	//deleteRes, err := c.DeleteBlog(context.Background(), deleteReq)
	//if err != nil {
	//	log.Fatalf("error while deleting item: %v", err)
	//}
	//
	//log.Printf("blog item deleted: %v", deleteRes)

	// LIST
	log.Println("Listing blog items")

	listReq := &blogpb.ListBlogRequest{}
	stream, err := c.ListBlog(context.Background(), listReq)
	if err != nil {
		log.Fatalf("error on list request: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while receiving blog list: %v", err)
		}
		log.Printf("blog item: %v", res.Blog)
	}

	log.Println("list ended")
}
