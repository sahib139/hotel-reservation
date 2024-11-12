package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rating   int                  `bson:"rating" json:"rating"`
	Room     []primitive.ObjectID `bson:"room" json:"room"`
}

type RoomType int

const (
	_ RoomType = iota
	SingleRoomType
	DoubleRoomType
	SeaTypeRoomType
	DeluxRoomType
)

type Room struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	// Type      RoomType           `bson:"type" json:"type"`
	Size      string             `bson:"size" json:"size"` // "small" or "medium" or "king"
	BasePrice float64            `bson:"basePrice" json:"basePrice"`
	Price     float64            `bson:"price" json:"price"`
	HotelID   primitive.ObjectID `bson:"hotelId" json:"hotelId"`
}
