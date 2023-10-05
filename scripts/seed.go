package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kevalsabhani/hotel-reservation/db"
	"github.com/kevalsabhani/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx        context.Context
	client     *mongo.Client
	hotelStore db.HotelStore
	roomStore  db.RoomStore
)

func seedHotel(name, location string, rating float64) {
	hotel := &types.Hotel{
		Name:     name,
		Location: location,
		Rating:   rating,
		Rooms:    []primitive.ObjectID{},
	}
	rooms := []*types.Room{
		{
			Type:  types.SingleRoomType,
			Price: 20000.0,
		},
		{
			Type:  types.DoubleRoomType,
			Price: 30000.0,
		},
		{
			Type:  types.DeluxeRoomType,
			Price: 35000.0,
		},
	}
	insertedHotel, err := hotelStore.InsertHotel(ctx, hotel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insertedHotel)

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		insertedRoom, err := roomStore.InsertRoom(ctx, room)
		if err != nil {
			log.Fatal(err)
		}
		filter := bson.M{"_id": insertedHotel.ID}
		update := bson.M{"$push": bson.M{"rooms": room.ID}}
		if err = hotelStore.UpdateHotel(ctx, filter, update); err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertedRoom)
	}

}

func main() {
	seedHotel("Taj", "India", 4.7)
	seedHotel("Marriot", "India", 4.3)
	seedHotel("Hyaat", "India", 4.5)
}

func init() {
	ctx = context.Background()
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client)
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
}
