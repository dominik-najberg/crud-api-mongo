package bootstrap

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewClient() (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:Zaq12wsx@localhost:27017"))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}
