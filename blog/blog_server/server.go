package main

import (
	"context"
	"github.com/dominik-najberg/crud-course/blog/blogpb"
	"github.com/dominik-najberg/crud-course/blog/bootstrap"
	"github.com/dominik-najberg/crud-course/blog/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"
)

var collection *mongo.Collection

type server struct{}

func (s *server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	// 5e5aad36034c2fef6df9fedf

	log.Printf("ReadBlog request: %v", req)
	blogID := req.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		log.Fatalf("blogID conversion error: %v", err)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	blogItem := model.BlogItem{}

	if err := collection.FindOne(ctx, filter).Decode(&blogItem); err != nil {
		return nil, status.Errorf(codes.NotFound, "not found: %v", err)
	}

	return &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{
			Id:       blogItem.ID.Hex(),
			AuthorId: blogItem.AuthorId,
			Title:    blogItem.Title,
			Content:  blogItem.Content,
		},
	}, nil
}

func (s *server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	log.Printf("CreateBlog request: %v", req)
	blog := req.GetBlog()

	data := model.BlogItem{
		AuthorId: blog.GetAuthorId(),
		Content:  blog.GetContent(),
		Title:    blog.GetTitle(),
	}

	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while inserting into DB: %v", err)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, "cannot convert object ID")
	}

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		},
	}, nil
}

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
	blogpb.RegisterBlogServiceServer(s, &server{})

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
