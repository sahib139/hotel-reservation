package main

import (
	"context"
	"fmt"

	"github.com/sahib139/hotel-reservation/db"
	"github.com/sahib139/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	hotelStore *db.MongoHotelStore
	roomStore  *db.MongoRoomStore
	userStore  *db.MongoUserStore
)

func seedUser(fName string, lName string, email string) {
	usr := &types.CreateUserParams{
		FirstName: fName,
		LastName:  lName,
		Email:     email,
		Password:  "password",
	}
	user, err := types.NewUserFromParams(*usr)
	if err != nil {
		panic(err)
	}
	insertedUser, err := userStore.InsertUser(context.Background(), user)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Inserted User: %+v\n", insertedUser)
}

func seedHotel(hotelName string, location string, rating int) {
	hotel := &types.Hotel{
		Name:     hotelName,
		Location: location,
		Rating:   rating,
		Room:     []primitive.ObjectID{},
	}
	insertedHostel, err := hotelStore.InsertHotel(context.Background(), hotel)

	if err != nil {
		panic(err)
	}

	rooms := []types.Room{
		{
			Size:      "small",
			BasePrice: 100.0,
		},
		{
			Size:      "medium",
			BasePrice: 200.0,
		},
		{
			Size:      "king",
			BasePrice: 300.0,
		},
	}

	for _, room := range rooms {
		room.HotelID = insertedHostel.ID
		room, err := roomStore.InsertRoom(context.Background(), &room)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Inserted Room: %+v\n", room)
	}
}

func main() {
	fmt.Println("Performing seeding...!")
	defer fmt.Println("Seeding Done!")

	// Create Hotels
	seedHotel("spark", "Mumbai", 4)
	seedHotel("current", "Delhi", 3)
	seedHotel("Arrow", "Gurugram", 5)
	seedHotel("Sunami", "Pune", 1)
	seedHotel("Technolo", "Surat", 2)

	// create Users
	seedUser("sahib", "singh", "sahib@singh.com")
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBUrl))
	if err != nil {
		panic(err)
	}

	hotelStore = db.NewMongoHotelStore(db.DBNAME, client)
	roomStore = db.NewMongoRoomStore(db.DBNAME, client, hotelStore)
	userStore = db.NewMongoUserStore(db.DBNAME, client)

	hotelStore.Drop(context.Background())
	roomStore.Drop(context.Background())
	userStore.Drop(context.Background())
}
