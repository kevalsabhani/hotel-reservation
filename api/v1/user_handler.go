package v1

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/kevalsabhani/hotel-reservation/db"
	"github.com/kevalsabhani/hotel-reservation/types"
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
	user, err := h.userStore.GetUserByID(c.Context(), id)
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
	if err := h.userStore.DeleteUser(c.Context(), userID); err != nil {
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
	if err := h.userStore.UpdateUser(c.Context(), userID, params); err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(map[string]string{"update": userID})
}
