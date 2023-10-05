package db

import (
	"context"

	"github.com/kevalsabhani/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const roomColl = "rooms"

type RoomStore interface {
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	GetRooms(context.Context) ([]*types.Room, error)
	GetRoomsByHotelID(context.Context, bson.M) ([]*types.Room, error)
	GetRoomByID(context.Context, bson.M) (*types.Room, error)
	// UpdateRoom(context.Context, bson.M, bson.M) error
	// DeleteRoom(context.Context, bson.M) error
	// DeleteRoomsByHotelID(context.Context, bson.M) error
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoRoomStore(client *mongo.Client) *MongoRoomStore {
	return &MongoRoomStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(roomColl),
	}
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = res.InsertedID.(primitive.ObjectID)
	return room, nil
}

func (s *MongoRoomStore) GetRooms(ctx context.Context) ([]*types.Room, error) {
	var rooms []*types.Room
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err := cur.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}

func (s *MongoRoomStore) GetRoomsByHotelID(ctx context.Context, filter bson.M) ([]*types.Room, error) {
	var rooms []*types.Room
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err := cur.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}

func (s *MongoRoomStore) GetRoomByID(ctx context.Context, filter bson.M) (*types.Room, error) {
	room := &types.Room{}
	if err := s.coll.FindOne(ctx, filter).Decode(room); err != nil {
		return nil, err
	}
	return room, nil
}

func (s *MongoRoomStore) UpdateRoom(ctx context.Context, filter bson.M, update bson.M) error {
	return nil
}

func (s *MongoRoomStore) DeleteRoom(ctx context.Context, filter bson.M) error {
	return nil
}

func (s *MongoRoomStore) DeleteRoomsByHotelID(ctx context.Context, filter bson.M) error {
	return nil
}
