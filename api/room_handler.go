package api

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sahib139/hotel-reservation/db"
	"github.com/sahib139/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	store db.DbStore
}

type BookRoomParams struct {
	NumPersons int       `json:"number"`
	CheckIn    time.Time `json:"checkIn"`
	CheckOut   time.Time `json:"checkOut"`
}

func NewRoomHandler(store *db.DbStore) *RoomHandler {
	return &RoomHandler{store: *store}
}

func (h *RoomHandler) HandlerBookRoom(c *fiber.Ctx) error {
	roomID := c.Params("id")
	var bookingParams BookRoomParams
	if err := c.BodyParser(&bookingParams); err != nil {
		return err
	}
	roomOID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&genericResponse{
				Msg:    "unauthorized",
				Status: 500,
			},
		)
	}
	if err := bookingParams.Validation(); err != nil {
		return err
	}

	booking := &types.Booking{
		RoomID:     roomOID,
		UserID:     user.ID,
		NumPersons: bookingParams.NumPersons,
		CheckIn:    bookingParams.CheckIn,
		CheckOut:   bookingParams.CheckOut,
	}

	ok, err = h.checkIfBookingPossible(c, booking)
	if err != nil {
		return err
	}
	if !ok {
		return c.Status(fiber.StatusConflict).JSON(
			&genericResponse{
				Msg:    "booking not possible",
				Status: 409,
			},
		)
	}
	booking, err = h.store.BookStore.InsertBooking(c.Context(), booking)
	if err != nil {
		return err
	}
	return c.JSON(booking)
}

func (h *RoomHandler) GetRooms(c *fiber.Ctx) error {
	filter := bson.M{}
	rooms, err := h.store.RoomStore.GetRoom(c.Context(), filter)

	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (b *BookRoomParams) Validation() error {
	time := time.Now()
	if time.After(b.CheckIn) || time.After(b.CheckOut) {
		return fmt.Errorf("CheckIn and CheckOut cannot be set in past")
	}
	return nil
}

func (h *RoomHandler) checkIfBookingPossible(c *fiber.Ctx, booking *types.Booking) (bool, error) {
	filter := bson.M{
		"roomId": booking.RoomID,
		"checkIn": bson.M{
			"$gte": booking.CheckIn,
		},
		"checkOut": bson.M{
			"$lte": booking.CheckOut,
		},
	}
	bookings, err := h.store.BookStore.GetBookings(c.Context(), filter)

	if err != nil {
		return false, err
	}
	if len(bookings) > 0 {
		return false, nil
	}
	return true, nil
}
