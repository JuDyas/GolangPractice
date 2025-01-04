package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// ConnectDatabase - setup connect to database
func ConnectDatabase(uri string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		//TODO: add zap logger
		log.Fatal(err)
		return nil
	}

	if err := client.Ping(ctx, nil); err != nil {
		//TODO: add zap logger
		log.Fatal(err)
		return nil
	}

	return client
}
