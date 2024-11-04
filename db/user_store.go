package db

import (
	"context"

	"github.com/sahib139/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	userCollections = "users"
)

type UserStore interface {
	GetUserByID(ctx context.Context, id string) (*types.User, error)
}

type mongoDbStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func (m *mongoDbStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": oid}

	var user types.User
	err = m.coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func NewMongoDBStore(client *mongo.Client) *mongoDbStore {
	return &mongoDbStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(userCollections),
	}
}
