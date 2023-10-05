package v1

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/kevalsabhani/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomHandler struct {
	rs db.RoomStore
}

func NewRoomHandler(rs db.RoomStore) *RoomHandler {
	return &RoomHandler{rs: rs}
}

func (h *RoomHandler) HandlePostRoom(c *fiber.Ctx) error {
	return nil
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.rs.GetRooms(c.Context())
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(map[string]string{"error": "Not found"})
		}
		return err
	}
	return c.Status(fiber.StatusOK).JSON(rooms)
}

func (h *RoomHandler) HandleGetRoom(c *fiber.Ctx) error {
	id := c.Params("id")
	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"error": "ID is invalid"})
	}
	filter := bson.M{"_id": oId}
	room, err := h.rs.GetRoomByID(c.Context(), filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(map[string]string{"error": "Not found"})
		}
		return err
	}
	return c.Status(fiber.StatusOK).JSON(room)
}

func (h *RoomHandler) HandlePutRoom(c *fiber.Ctx) error {
	return nil
}

func (h *RoomHandler) HandleDeleteRoom(c *fiber.Ctx) error {
	return nil
}
