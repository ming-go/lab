package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	errors "golang.org/x/xerrors"
)

type ModelResult struct {
	Value float64
}

func Model() (*ModelResult, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:8787"))
	if err != nil {
		return nil, errors.Errorf("mongo.NewClient: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, errors.Errorf("client.Connect: %w", err)
	}

	var result ModelResult
	collection := client.Database("testing").Collection("testing")
	err = collection.FindOne(ctx, bson.M{"name": "pi"}).Decode(&result)
	if err != nil {
		return nil, errors.Errorf("collection.FindOne: %w", err)
	}

	return &result, nil
}
