package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sahib139/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.DbStore
}

func NewHotelHandler(store *db.DbStore) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

type HotelQueryParams struct {
	Room   bool `json:"room"`
	Rating int  `json:"rating"`
}

func (h *HotelHandler) HandlerGetHotels(c *fiber.Ctx) error {
	var query HotelQueryParams
	if err := c.QueryParser(&query); err != nil {
		return err
	}

	hotels, err := h.store.HotelStore.GetHotel(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}

func (h *HotelHandler) HandlerGetHotel(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)
	hotel, err := h.store.HotelStore.GetHotelByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(hotel)
}

func (h *HotelHandler) HandlerGetRoom(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"hotelId": oid}
	room, err := h.store.RoomStore.GetRoom(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(room)
}
