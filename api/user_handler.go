package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sahib139/hotel-reservation/db"
	"github.com/sahib139/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
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
		id = c.Params("id")
	)
	user, err := h.UserStore.GetUserByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if errors := params.Validation(); len(errors) > 0 {
		return c.JSON(errors)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return nil
	}
	insertedUser, err := h.UserStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {

	users, err := h.UserStore.GetUsers(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(users)
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)
	err := h.UserStore.DeleteUser(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(map[string]string{"msg": "User Deleted!"})
}

func (h *UserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	var (
		id     = c.Params("id")
		params types.UpdateUserParams
	)

	if err := c.BodyParser(&params); err != nil {
		return err
	}
	updates := params.ToBson()
	err := h.UserStore.UpdateUser(c.Context(), bson.M{"_id": id}, updates)
	if err != nil {
		return err
	}
	return c.JSON(map[string]string{"msg": "User Updated!"})
}
