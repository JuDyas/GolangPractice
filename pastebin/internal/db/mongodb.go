package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectDatabase(uri string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Println("failed to connect to MongoDB: %w", err)
		return nil
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Println("failed to ping MongoDB: %w", err)
		return nil
	}

	fmt.Println("Connected to MongoDB!")
	return client
}
