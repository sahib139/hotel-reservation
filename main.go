package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sahib139/hotel-reservation/api"
	"github.com/sahib139/hotel-reservation/api/middleware"
	"github.com/sahib139/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var errorConfig = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		if err != nil {
			ctx.JSON(map[string]string{"err": fmt.Sprintf("%v", err)})
		}
		return nil
	},
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBUrl))
	if err != nil {
		log.Fatal(err)
	}

	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	var (
		app = fiber.New(errorConfig)

		auth  = app.Group("/api")
		apiv1 = app.Group("/api/v1", middleware.JWTAuthentication)

		// Handler's initialization
		store = db.NewDbStore(client)

		userHandler  = api.NewUserHandler(store)
		hotelHandler = api.NewHotelHandler(store)
		authHandler  = api.NewAuthHandler(store)
	)

	//auth
	auth.Post("/auth", authHandler.HandleAuthentication)

	// User endpoints
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/users/:id", userHandler.HandleGetUser)
	apiv1.Post("/users", userHandler.HandlePostUser)
	apiv1.Delete("/users/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/users/:id", userHandler.HandleUpdateUser)

	// Hotel endpoints
	apiv1.Get("/hotels", hotelHandler.HandlerGetHotels)
	apiv1.Get("/hotels/:id", hotelHandler.HandlerGetHotel)
	apiv1.Get("/hotels/:id/rooms", hotelHandler.HandlerGetRoom)

	app.Listen(*listenAddr)
}
