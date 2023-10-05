package v1

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/kevalsabhani/hotel-reservation/db"
	"github.com/kevalsabhani/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandlePostHotel(c *fiber.Ctx) error {
	params := &types.CreateHotelParams{}
	if err := c.BodyParser(params); err != nil {
		return err
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	hotel := types.CreateHotelFromParams(params)
	hotel, err := h.store.Hotel.InsertHotel(c.Context(), hotel)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}
	return c.Status(fiber.StatusCreated).JSON(hotel)
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	hotels, err := h.store.Hotel.GetHotels(c.Context())
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(map[string]string{"error": "Not found"})
		}
		return err
	}
	return c.Status(fiber.StatusOK).JSON(hotels)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"error": "ID is invalid"})
	}
	filter := bson.M{"_id": oId}
	user, err := h.store.Hotel.GetHotelByID(c.Context(), filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(map[string]string{"error": "Not found"})
		}
		return err
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *HotelHandler) HandleGetRoomsByHotelID(c *fiber.Ctx) error {
	id := c.Params("id")
	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"error": "ID is invalid"})
	}
	filter := bson.M{"hotelId": oId}
	rooms, err := h.store.Room.GetRoomsByHotelID(c.Context(), filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(map[string]string{"error": "Not found"})
		}
		return err
	}
	return c.Status(fiber.StatusOK).JSON(rooms)
}

func (h *HotelHandler) HandlePutHotel(c *fiber.Ctx) error {
	userID := c.Params("id")
	params := &types.UpdateHotelParams{}
	if err := c.BodyParser(params); err != nil {
		return err
	}
	oId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"error": "ID is invalid"})
	}
	filter := bson.M{"_id": oId}
	update := bson.M{"$set": params.ToBSON()}
	if err := h.store.Hotel.UpdateHotel(c.Context(), filter, update); err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(map[string]string{"updated": userID})
}

func (h *HotelHandler) HandleDeleteHotel(c *fiber.Ctx) error {
	hotelID := c.Params("id")
	oId, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"error": "ID is invalid"})
	}
	filter := bson.M{"_id": oId}
	if err := h.store.Hotel.DeleteHotel(c.Context(), filter); err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(map[string]string{"deleted": hotelID})
}
