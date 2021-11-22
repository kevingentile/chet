package infrastructure

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//TODO client disconnect
type MongoStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewDefaultMongoStore(mongoUri string) (*MongoStore, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		return nil, err
	}
	return &MongoStore{
		client:     client,
		collection: client.Database("chet").Collection("events")}, nil
}

func (m *MongoStore) Save(message interface{}) error {
	_, err := m.collection.InsertOne(context.TODO(), message)
	return err
}
