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

	userHandler := v1.NewUserHandler(db.NewMongoUserStore(client))
	hotelHandler := v1.NewHotelHandler(db.NewMongoHotelStore(client), db.NewMongoRoomStore(client))

	app := fiber.New(
		fiber.Config{
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				return c.JSON(map[string]string{"error": err.Error()})
			},
		})
	api := app.Group("/api")

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

	app.Listen(*listenAddr)
}
