package api

import (
	"fmt"
	"hotel-reservation/db/fixtures"
	"testing"
	"time"
)

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	user := fixtures.AddUser(tdb.Store, "james", "foo", false)
	hotel := fixtures.AddHotel(tdb.Store, "bar hotel", "a", 4, nil)
	room := fixtures.AddRoom(tdb.Store, "small", true, 4.4, hotel.ID)

	from := time.Now()
	till := from.AddDate(0, 0, 5)
	booking := fixtures.AddBooking(tdb.Store, user.ID, room.ID, from, till, 2)
	fmt.Println(booking)
}
