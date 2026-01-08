package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/lugondev/go-indexer-solana-starter/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

func NewMongoRepository(uri, dbName string) (*MongoRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("connect to mongodb: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping mongodb: %w", err)
	}

	database := client.Database(dbName)
	collection := database.Collection("events")

	return &MongoRepository{
		client:     client,
		database:   database,
		collection: collection,
	}, nil
}

func (r *MongoRepository) SaveEvent(ctx context.Context, event interface{}) error {
	_, err := r.collection.InsertOne(ctx, event)
	if err != nil {
		return fmt.Errorf("insert event: %w", err)
	}
	return nil
}

func (r *MongoRepository) GetEventsByTimeRange(ctx context.Context, from, to time.Time) ([]models.BaseEvent, error) {
	filter := bson.M{
		"block_time": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("find events: %w", err)
	}
	defer cursor.Close(ctx)

	var events []models.BaseEvent
	if err := cursor.All(ctx, &events); err != nil {
		return nil, fmt.Errorf("decode events: %w", err)
	}

	return events, nil
}

func (r *MongoRepository) GetEventsByType(ctx context.Context, eventType models.EventType, limit int) ([]interface{}, error) {
	filter := bson.M{"event_type": eventType}
	opts := options.Find().SetLimit(int64(limit)).SetSort(bson.D{{Key: "block_time", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("find events by type: %w", err)
	}
	defer cursor.Close(ctx)

	var events []interface{}
	if err := cursor.All(ctx, &events); err != nil {
		return nil, fmt.Errorf("decode events: %w", err)
	}

	return events, nil
}

func (r *MongoRepository) GetEventBySignature(ctx context.Context, signature string) (interface{}, error) {
	filter := bson.M{"signature": signature}

	var event interface{}
	if err := r.collection.FindOne(ctx, filter).Decode(&event); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("find event by signature: %w", err)
	}

	return event, nil
}

func (r *MongoRepository) Close(ctx context.Context) error {
	return r.client.Disconnect(ctx)
}

func (r *MongoRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "signature", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "event_type", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "block_time", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "slot", Value: -1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("create indexes: %w", err)
	}

	return nil
}
