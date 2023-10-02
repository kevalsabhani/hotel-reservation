package db

import (
	"context"

	"github.com/kevalsabhani/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userColl = "users"

type UserStore interface {
	InsertUser(context.Context, *types.User) (*types.User, error)
	GetUserByID(context.Context, bson.M) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	DeleteUser(context.Context, bson.M) error
	UpdateUser(context.Context, bson.M, bson.M) error
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(userColl),
	}
}

func (s *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	users := []*types.User{}
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *MongoUserStore) GetUserByID(ctx context.Context, filter bson.M) (*types.User, error) {
	user := &types.User{}
	if err := s.coll.FindOne(ctx, filter).Decode(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, filter bson.M) error {
	_, err := s.coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := s.coll.UpdateOne(ctx, filter, update)
	return err
}
