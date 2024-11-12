package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/sahib139/hotel-reservation/db"
	"github.com/sahib139/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewTestUserHandler(store *db.DbStore) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

func NewTestStore(client *mongo.Client, userStore db.UserStore, hotelStore db.HotelStore, roomStore db.RoomStore) *db.DbStore {
	return &db.DbStore{
		Client:     client,
		UserStore:  userStore,
		HotelStore: hotelStore,
		RoomStore:  roomStore,
	}
}

func CreateTestUserHandler() *UserHandler {
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBUrl))
	userStore := db.NewMongoUserStore(db.TestDb, client)
	hotelStore := db.NewMongoHotelStore(db.TestDb, client)
	roomStore := db.NewMongoRoomStore(db.TestDb, client, hotelStore)
	store := NewTestStore(client, userStore, hotelStore, roomStore)
	return NewTestUserHandler(store)
}

func TestUserPostRequest(t *testing.T) {
	testUserHandler := CreateTestUserHandler()
	defer testUserHandler.store.UserStore.Drop(context.Background())

	api := fiber.New()
	api.Post("/", testUserHandler.HandlePostUser)

	reqBody := types.CreateUserParams{
		FirstName: "Ashutosh",
		LastName:  "Sharma",
		Email:     "ashutosh@example.com",
		Password:  "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	res, _ := api.Test(req)

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		t.Error("Expected status code 201, got", res.StatusCode)
	}
	var user types.CreateUserParams
	json.NewDecoder(req.Body).Decode(&user)

	if user.FirstName != "Ashutosh" || user.LastName != "Sharma" || user.Email != "ashutosh@example.com" || user.Password != "password123" {
		t.Error("Expected user details, got", user)
	}

	insetUser, _ := types.NewUserFromParams(user)

	createdUser, _ := testUserHandler.store.UserStore.InsertUser(context.Background(), insetUser)

	if createdUser.ID.IsZero() {
		t.Error("Expected user ID, got", createdUser.ID)
	}

}
