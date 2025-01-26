package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sahib139/hotel-reservation/db"
	"github.com/sahib139/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.DbStore
}

func NewBookingHandler(store *db.DbStore) *BookingHandler {
	return &BookingHandler{store: store}
}

func (h *BookingHandler) HandleBookings(c *fiber.Ctx) error {
	user := c.Context().Value("user").(*types.User)
	bookings, err := h.store.BookStore.GetBookings(c.Context(), bson.M{"userId": user.ID})
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}
