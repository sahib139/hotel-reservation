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

func main() {
	fmt.Println("Performing seeding...!")
	defer fmt.Println("Seeding Done!")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBUrl))
	if err != nil {
		panic(err)
	}

	hotelStore := db.NewMongoHotelStore(db.DBNAME, client)
	roomStore := db.NewMongoRoomStore(db.DBNAME, client, hotelStore)

	hotelStore.Drop(context.Background())
	roomStore.Drop(context.Background())

	hotel := &types.Hotel{
		Name:     "The Hotel",
		Location: "New York",
		Room:     []primitive.ObjectID{},
	}
	insertedHostel, err := hotelStore.InsertHotel(context.Background(), hotel)

	if err != nil {
		panic(err)
	}

	rooms := []types.Room{
		{
			Type:      types.SingleRoomType,
			BasePrice: 100.0,
		},
		{
			Type:      types.DoubleRoomType,
			BasePrice: 200.0,
		},
		{
			Type:      types.SeaTypeRoomType,
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
