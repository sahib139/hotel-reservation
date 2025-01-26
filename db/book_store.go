package db

import (
	"context"

	"github.com/sahib139/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookStore interface {
	InsertBooking(context context.Context, book *types.Booking) (*types.Booking, error)
	GetBookings(context context.Context, filter bson.M) ([]*types.Booking, error)
	GetBookingById(context context.Context, id string) (*types.Booking, error)
}
type MongoBookStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoBookStore(client *mongo.Client) *MongoBookStore {
	return &MongoBookStore{
		client: nil,
		coll:   client.Database(DBNAME).Collection("bookings"),
	}
}

func (m *MongoBookStore) InsertBooking(ctx context.Context, book *types.Booking) (*types.Booking, error) {
	insertBooking, err := m.coll.InsertOne(ctx, book)
	if err != nil {
		return nil, err
	}
	book.ID = insertBooking.InsertedID.(primitive.ObjectID)
	return book, nil
}

func (m *MongoBookStore) GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error) {
	cur, err := m.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var bookings []*types.Booking
	if err := cur.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (m *MongoBookStore) GetBookingById(ctx context.Context, id string) (*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var booking types.Booking
	err = m.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}
