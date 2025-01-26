package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sahib139/hotel-reservation/db"
	"github.com/sahib139/hotel-reservation/types"
)

func JWTAuthentication(store *db.DbStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			return fmt.Errorf("unauthorized client")
		}
		val, err := parseJWT(token[0])

		if err != nil {
			return err
		}

		user_id := val["id"].(string)

		user, err := store.UserStore.GetUserByID(c.Context(), user_id)

		if err != nil {
			return fmt.Errorf("unauthorized client")
		}

		c.Context().SetUserValue("user", user)

		// fmt.Printf("user_id : %v\n", user_id)

		// pre_body := map[string]interface{}{}
		// if err := c.BodyParser(pre_body); err != nil {
		// 	return err
		// }
		// pre_body["user_id"] = user_id
		// body, err := json.Marshal(pre_body)
		// if err != nil {
		// 	return err
		// }
		// c.Request().SetBody(body)
		return c.Next()
	}
}

func parseJWT(JWTtoken string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(JWTtoken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Printf("Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("unauthorized")
	}
}

func GenerateToken(user *types.User) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"nbf":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 4).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))

	return tokenString, err
}
