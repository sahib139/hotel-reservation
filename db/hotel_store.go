package db

import (
	"context"

	"github.com/sahib139/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelStore interface {
	GetHotelByID(ctx context.Context, id string) (*types.Hotel, error)
	GetHotel(ctx context.Context) ([]*types.Hotel, error)
	InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error)
	DeleteHotel(ctx context.Context, id string) error
	UpdateHotel(ctx context.Context, filter bson.M, update bson.M) error
	Drop(ctx context.Context) error
}

type mongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func (m *mongoHotelStore) GetHotelByID(ctx context.Context, id string) (*types.Hotel, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": oid}
	var hotel types.Hotel
	err = m.coll.FindOne(ctx, filter).Decode(&hotel)
	if err != nil {
		return nil, err
	}
	return &hotel, nil
}

func (m *mongoHotelStore) GetHotel(ctx context.Context) ([]*types.Hotel, error) {
	filter := bson.M{"_id": ""}
	cur, err := m.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var hotels []*types.Hotel
	if err := cur.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}

func (m *mongoHotelStore) DeleteHotel(ctx context.Context, id string) error {
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

func (m *mongoHotelStore) UpdateHotel(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := m.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (m *mongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := m.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = res.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (m *mongoHotelStore) Drop(ctx context.Context) error {
	return m.coll.Drop(ctx)
}

func NewMongoHotelStore(DBNAME string, client *mongo.Client) *mongoHotelStore {
	return &mongoHotelStore{
		client: client,
		coll:   client.Database(DBNAME).Collection("hotels"),
	}
}
