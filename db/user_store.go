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
	GetUserByID(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(context.Context, string, *types.UpdateUserParams) error
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

func (s *MongoUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	user := &types.User{}
	if err := s.coll.FindOne(ctx, bson.M{"_id": oId}).Decode(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": oId})
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, id string, values *types.UpdateUserParams) error {
	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = s.coll.UpdateOne(ctx, bson.M{"_id": oId}, bson.M{"$set": values.ToBSON()})
	if err != nil {
		return err
	}
	return nil
}
