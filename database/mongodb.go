package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	greenColor = "\033[32m"
	resetColor = "\033[0m"
)

func colorize(text string, colorCode string) string {
	return colorCode + text + resetColor
}

func ConnectMongo(uri string, tag string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	mongoTag := colorize("[MongoDB]", greenColor)
	fmt.Printf("%s Connected to MongoDB at %s\n", mongoTag, uri)

	return client, nil
}
