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

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	params := &types.CreateUserParams{}
	if err := c.BodyParser(params); err != nil {
		return err
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	user, err := types.CreateUserFromParams(params)
	if err != nil {
		return err
	}
	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(insertedUser)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(users)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"error": "ID is invalid"})
	}
	filter := bson.M{"_id": oId}
	user, err := h.userStore.GetUserByID(c.Context(), filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(map[string]string{"error": "Not found"})
		}
		return err
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	oId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": oId}
	if err := h.userStore.DeleteUser(c.Context(), filter); err != nil {
		return nil
	}
	return c.Status(fiber.StatusOK).JSON(map[string]string{"deleted": userID})
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	params := &types.UpdateUserParams{}
	if err := c.BodyParser(params); err != nil {
		return err
	}
	oId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"error": "ID is invalid"})
	}
	filter := bson.M{"_id": oId}
	update := bson.M{"$set": params.ToBSON()}
	if err := h.userStore.UpdateUser(c.Context(), filter, update); err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(map[string]string{"update": userID})
}
