package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID     primitive.ObjectID `bson:"userId,omitempty" json:"userId,omitempty"`
	RoomID     primitive.ObjectID `bson:"roomId,omitempty" json:"roomId,omitempty"`
	NumPersons int                `bson:"numPersons,omitempty" json:"numPersons,omitempty"`
	CheckIn    time.Time          `bson:"checkIn,omitempty" json:"checkIn,omitempty"`
	CheckOut   time.Time          `bson:"checkOut,omitempty" json:"checkOut,omitempty"`
}
