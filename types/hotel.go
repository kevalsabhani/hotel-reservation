package types

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Hotel struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
	Rating   float64              `bson:"rating" json:"rating"`
}

type CreateHotelParams struct {
	Name     string  `json:"name"`
	Location string  `json:"location"`
	Rating   float64 `json:"rating"`
}

type UpdateHotelParams struct {
	CreateHotelParams
}

func (p *CreateHotelParams) Validate() map[string]string {
	errors := make(map[string]string)
	if p.Name == "" {
		errors["name"] = "Name is required"
	}
	if p.Location == "" {
		errors["location"] = "Location is required"
	}
	if p.Rating == 0 {
		errors["rating"] = "Rating is required"
	}
	return errors
}

func CreateHotelFromParams(params *CreateHotelParams) *Hotel {
	return &Hotel{
		Name:     params.Name,
		Location: params.Location,
		Rooms:    []primitive.ObjectID{},
		Rating:   params.Rating,
	}
}

func (params *UpdateHotelParams) ToBSON() bson.M {
	m := bson.M{}
	if len(params.Name) > 0 {
		m["name"] = params.Name
	}
	if len(params.Location) > 0 {
		m["location"] = params.Location
	}
	if params.Rating > 0 {
		m["rating"] = params.Rating
	}
	return m
}
