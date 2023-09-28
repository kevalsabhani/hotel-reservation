package v1

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/kevalsabhani/hotel-reservation/db"
	"github.com/kevalsabhani/hotel-reservation/types"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	return c.JSON(&types.User{ID: "", FirstName: "keval", LastName: "sabhani"})
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id  = c.Params("id")
		ctx = context.Background()
	)
	user, err := h.userStore.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}
