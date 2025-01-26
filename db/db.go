package db

import "go.mongodb.org/mongo-driver/mongo"

const (
	DBUrl  = "mongodb://localhost:27017"
	DBNAME = "hotel-reservation"
	TestDb = "hotel-reservation-test"
)

type DbStore struct {
	Client     *mongo.Client
	UserStore  UserStore
	HotelStore HotelStore
	RoomStore  RoomStore
	BookStore  BookStore
}

func NewDbStore(client *mongo.Client) *DbStore {
	hotelStore := NewMongoHotelStore(DBNAME, client)
	return &DbStore{
		Client:     client,
		UserStore:  NewMongoUserStore(DBNAME, client),
		HotelStore: hotelStore,
		RoomStore:  NewMongoRoomStore(DBNAME, client, hotelStore),
		BookStore:  NewMongoBookStore(client),
	}
}
