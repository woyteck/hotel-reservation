package api

import (
	"hotel-reservation/db"

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
		return ErrResourceNotFound("booking")
	}

	user, err := getAuthUser(c)
	if err != nil {
		return ErrUnauthorized()
	}

	if booking.UserID != user.ID {
		return ErrUnauthorized()
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
		return ErrResourceNotFound("booking")
	}

	return c.JSON(bookings)
}

// this needs to be user authorized
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	booking, err := h.store.Booking.GetBookingByID(c.Context(), c.Params("id"))
	if err != nil {
		return ErrResourceNotFound("booking")
	}

	user, err := getAuthUser(c)
	if err != nil {
		return ErrUnauthorized()
	}

	if booking.UserID != user.ID {
		return ErrUnauthorized()
	}

	return c.JSON(booking)
}
