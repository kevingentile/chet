package infrastructure

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//TODO client disconnect
type MongoStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewDefaultMongoStore(client *mongo.Client, database string, collection string) (*MongoStore, error) {
	return &MongoStore{
		client:     client,
		collection: client.Database(database).Collection(collection)}, nil
}

func (m *MongoStore) Update(id uuid.UUID, view interface{}) error {
	_, err := m.collection.ReplaceOne(context.TODO(), bson.M{"id": id}, view)
	return err
}

func (m *MongoStore) Create(view interface{}) error {
	_, err := m.collection.InsertOne(context.TODO(), view)
	return err
}
