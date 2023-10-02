package db

import (
	"context"

	"github.com/kevalsabhani/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const hotelColl = "hotels"

type HotelStore interface {
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	GetHotels(context.Context) ([]*types.Hotel, error)
	GetHotelByID(context.Context, bson.M) (*types.Hotel, error)
	UpdateHotel(context.Context, bson.M, bson.M) error
	DeleteHotel(context.Context, bson.M) error
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(hotelColl),
	}
}

func (s *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = res.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (s *MongoHotelStore) GetHotels(ctx context.Context) ([]*types.Hotel, error) {
	var hotels []*types.Hotel
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err := cur.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}

func (s *MongoHotelStore) GetHotelByID(ctx context.Context, filter bson.M) (*types.Hotel, error) {
	hotel := &types.Hotel{}
	if err := s.coll.FindOne(ctx, filter).Decode(hotel); err != nil {
		return nil, err
	}
	return hotel, nil
}

func (s *MongoHotelStore) UpdateHotel(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := s.coll.UpdateOne(ctx, filter, update)
	return err
}

func (s *MongoHotelStore) DeleteHotel(ctx context.Context, filter bson.M) error {
	_, err := s.coll.DeleteOne(ctx, filter)
	return err
}
