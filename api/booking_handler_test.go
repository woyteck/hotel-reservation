package api

import (
	"encoding/json"
	"fmt"
	"hotel-reservation/db/fixtures"
	"hotel-reservation/types"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestUserGetBooking(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		nonAuthUser = fixtures.AddUser(tdb.Store, "Jimmy", "Watercooler", false)
		user        = fixtures.AddUser(tdb.Store, "james", "foo", false)
		hotel       = fixtures.AddHotel(tdb.Store, "bar hotel", "a", 4, nil)
		room        = fixtures.AddRoom(tdb.Store, "small", true, 4.4, hotel.ID)
		from        = time.Now()
		till        = from.AddDate(0, 0, 5)
		app         = fiber.New(fiber.Config{
			ErrorHandler: ErrorHandler,
		})
		route          = app.Group("/", JWTAuthentication(tdb.User))
		bookingHandler = NewBookingHandler(tdb.Store)
		booking        = fixtures.AddBooking(tdb.Store, user.ID, room.ID, from, till, 2)
	)

	route.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 code got %d", resp.StatusCode)
	}

	var bookingResponse *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResponse); err != nil {
		t.Fatal(err)
	}

	if bookingResponse.ID != booking.ID {
		t.Fatalf("expected %s got %s", booking.ID, bookingResponse.ID)
	}

	if bookingResponse.UserID != booking.UserID {
		t.Fatalf("expected %s got %s", booking.UserID, bookingResponse.UserID)
	}

	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected non 200 response got %d", resp.StatusCode)
	}
}

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		adminUser      = fixtures.AddUser(tdb.Store, "admin", "admin", true)
		user           = fixtures.AddUser(tdb.Store, "james", "foo", false)
		hotel          = fixtures.AddHotel(tdb.Store, "bar hotel", "a", 4, nil)
		room           = fixtures.AddRoom(tdb.Store, "small", true, 4.4, hotel.ID)
		from           = time.Now()
		till           = from.AddDate(0, 0, 5)
		app            = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		admin          = app.Group("/", JWTAuthentication(tdb.User), AdminAuth)
		bookingHandler = NewBookingHandler(tdb.Store)
		booking        = fixtures.AddBooking(tdb.Store, user.ID, room.ID, from, till, 2)
	)

	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response got %d", resp.StatusCode)
	}

	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}

	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking, got %d", len(bookings))
	}

	have := bookings[0]
	if have.ID != booking.ID {
		t.Fatalf("expected %s got %s", booking.ID, have.ID)
	}

	if have.UserID != booking.UserID {
		t.Fatalf("expected %s got %s", booking.UserID, have.UserID)
	}

	//test non-admin cannot access the bookings
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status unauthorized response got %d", resp.StatusCode)
	}
}
