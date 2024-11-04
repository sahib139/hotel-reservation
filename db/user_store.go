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
	GetUsers(ctx context.Context) ([]*types.User, error)
	InsertUser(ctx context.Context, user *types.User) (*types.User, error)
}

type mongoDbStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func (m *mongoDbStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	insertedUser, err := m.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = insertedUser.InsertedID.(primitive.ObjectID)
	return user, nil
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

func (m *mongoDbStore) GetUsers(ctx context.Context) ([]*types.User, error) {

	filter := bson.M{}
	cur, err := m.coll.Find(ctx, filter)

	if err != nil {
		return nil, err
	}
	// defer cur.Close(ctx)

	var users []*types.User
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}
	// for cur.Next(ctx) {
	// 	var user types.User
	// 	if err := cur.Decode(&user); err != nil {
	// 		return nil, err
	// 	}
	// 	users = append(users, &user)
	// }

	// if err := cur.Err(); err != nil {
	// 	return nil, err
	// }

	return users, nil
}

func NewMongoDBStore(client *mongo.Client) *mongoDbStore {
	return &mongoDbStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(userCollections),
	}
}
