package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sahib139/hotel-reservation/api/middleware"
	"github.com/sahib139/hotel-reservation/db"
	"github.com/sahib139/hotel-reservation/types"
)

type AuthHandler struct {
	store *db.DbStore
}

func NewAuthHandler(store *db.DbStore) *AuthHandler {
	return &AuthHandler{store: store}
}

type AuthUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

func (h *AuthHandler) HandleAuthentication(c *fiber.Ctx) error {
	var auth AuthUser

	if err := c.BodyParser(&auth); err != nil {
		return err
	}

	user, err := h.store.UserStore.GetUserByEmail(c.Context(), auth.Email)
	if err != nil {
		return err
	}

	if !types.IsValidPassword(user.EncryptedPassword, auth.Password) {
		return fmt.Errorf("unauthorized")
	}

	token, err := middleware.GenerateToken(user)
	if err != nil {
		return err
	}
	return c.JSON(&AuthResponse{User: user, Token: token})
}
