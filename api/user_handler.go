package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sahib139/hotel-reservation/types"
)

func HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "James",
		LastName:  "At the waterCooler",
	}
	// return c.JSON(map[string]string{"user": "Sahib Singh"})
	return c.JSON(u)
}

func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"user": "Sahib Singh"})
}