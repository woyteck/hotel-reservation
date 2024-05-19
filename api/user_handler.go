package api

import (
	"hotel-reservation/types"

	"github.com/gofiber/fiber/v2"
)

func HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "James",
		LastName:  "Hetfield",
	}
	return c.JSON(u)
}

func HandleGetUser(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "James",
		LastName:  "Hetfield",
	}
	return c.JSON(u)
}
