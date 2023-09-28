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

const (
	dbUri    = "mongodb://localhost:27017"
	database = "hotel-reservation"
	userColl = "users"
)

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "Get Server Port")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))
	if err != nil {
		log.Fatal(err)
	}

	userHandler := v1.NewUserHandler(db.NewMongoUserStore(client))

	app := fiber.New(
		fiber.Config{
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				return c.JSON(map[string]string{"error": err.Error()})
			},
		})
	// status check
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(map[string]string{"status": "ok"})
	})

	// /api route
	api := app.Group("/api")

	// /api/v1 route
	v1Route := api.Group("/v1")
	v1Route.Get("/users", userHandler.HandleGetUsers)
	v1Route.Get("/users/:id", userHandler.HandleGetUser)

	app.Listen(*listenAddr)
}
