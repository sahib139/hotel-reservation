package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/sahib139/hotel-reservation/db"
)

type UserHandler struct {
	UserStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		UserStore: userStore,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id  = c.Params("id")
		ctx = context.Background()
	)
	user, err := h.UserStore.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

// func HandleGetUsers(c *fiber.Ctx) error {
// 	u := types.User{
// 		FirstName: "James",
// 		LastName:  "At the waterCooler",
// 	}
// 	// return c.JSON(map[string]string{"user": "Sahib Singh"})
// 	return c.JSON(u)
// }

// func HandleGetUser(c *fiber.Ctx) error {
// 	return c.JSON(map[string]string{"user": "Sahib Singh"})
// }
