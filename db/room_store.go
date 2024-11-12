package db

import (
	"context"

	"github.com/sahib139/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	GetRoomByID(ctx context.Context, id string) (*types.Room, error)
	GetRoom(ctx context.Context, filter bson.M) ([]*types.Room, error)
	InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error)
	DeleteRoom(ctx context.Context, id string) error
	UpdateRoom(ctx context.Context, filter bson.M, update bson.M) error
	Drop(ctx context.Context) error
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	hotelstore HotelStore
}

func (m *MongoRoomStore) GetRoomByID(ctx context.Context, id string) (*types.Room, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": oid}
	var room types.Room
	err = m.coll.FindOne(ctx, filter).Decode(&room)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (m *MongoRoomStore) GetRoom(ctx context.Context, filter bson.M) ([]*types.Room, error) {
	cur, err := m.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var rooms []*types.Room
	if err := cur.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}

func (m *MongoRoomStore) DeleteRoom(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = m.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoRoomStore) UpdateRoom(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := m.coll.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	result, err := m.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = result.InsertedID.(primitive.ObjectID)
	m.hotelstore.UpdateHotel(ctx, bson.M{"_id": room.HotelID}, bson.M{"$push": bson.M{"room": room.ID}})
	return room, nil
}

func (m *MongoRoomStore) Drop(ctx context.Context) error {
	return m.coll.Drop(ctx)
}

func NewMongoRoomStore(DBNAME string, client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(DBNAME).Collection("rooms"),
		hotelstore: hotelStore,
	}
}
