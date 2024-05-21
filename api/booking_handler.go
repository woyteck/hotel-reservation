package api

import (
	"hotel-reservation/db"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	booking, err := h.store.Booking.GetBookingByID(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}

	user, err := getAuthUser(c)
	if err != nil {
		return err
	}

	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResponse{
			Type: "error",
			Msg:  "unauthorized",
		})
	}

	if err := h.store.Booking.UpdateBooking(c.Context(), c.Params("id"), bson.M{"isCancelled": true}); err != nil {
		return err
	}

	return c.JSON(genericResponse{
		Type: "msg",
		Msg:  "updated",
	})
}

// this needs to be admin authorized
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), nil)
	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

// this needs to be user authorized
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	booking, err := h.store.Booking.GetBookingByID(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}

	user, err := getAuthUser(c)
	if err != nil {
		return err
	}

	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResponse{
			Type: "error",
			Msg:  "unauthorized",
		})
	}

	return c.JSON(booking)
}
