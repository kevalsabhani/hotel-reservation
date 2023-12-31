package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	v1 "github.com/kevalsabhani/hotel-reservation/api/v1"
	"github.com/kevalsabhani/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "Get Server Port")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	var (
		userStore  = db.NewMongoUserStore(client)
		hotelStore = db.NewMongoHotelStore(client)
		roomStore  = db.NewMongoRoomStore(client)
		store      = &db.Store{
			User:  userStore,
			Hotel: hotelStore,
			Room:  roomStore,
		}
		userHandler  = v1.NewUserHandler(userStore)
		hotelHandler = v1.NewHotelHandler(store)
		roomHandler  = v1.NewRoomHandler(roomStore)
		app          = fiber.New(
			fiber.Config{
				ErrorHandler: func(c *fiber.Ctx, err error) error {
					return c.JSON(map[string]string{"error": err.Error()})
				},
			})
		api = app.Group("/api")
	)

	// /api/v1 route
	v1Route := api.Group("/v1")

	// user routes
	v1Route.Get("/users", userHandler.HandleGetUsers)
	v1Route.Post("/users", userHandler.HandlePostUser)
	v1Route.Get("/users/:id", userHandler.HandleGetUser)
	v1Route.Delete("/users/:id", userHandler.HandleDeleteUser)
	v1Route.Put("/users/:id", userHandler.HandlePutUser)

	// hotel routes
	v1Route.Get("/hotels", hotelHandler.HandleGetHotels)
	v1Route.Get("/hotels/:id", hotelHandler.HandleGetHotel)
	v1Route.Get("/hotels/:id/rooms", hotelHandler.HandleGetRoomsByHotelID)
	v1Route.Put("/hotels/:id", hotelHandler.HandleGetHotel)
	v1Route.Delete("/hotels/:id", hotelHandler.HandleGetHotel)

	// room routes
	v1Route.Get("/rooms", roomHandler.HandleGetRooms)
	v1Route.Get("/rooms/:id", roomHandler.HandleGetRoom)
	// v1Route.Put("/rooms/:id", roomHandler.HandleGetRoom)
	// v1Route.Delete("/rooms/:id", roomHandler.HandleGetRoom)

	app.Listen(*listenAddr)
}
