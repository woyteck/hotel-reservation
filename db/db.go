package db

const (
	DBNAME     = "hotel-reservation"
	TestDBNAME = "hotel-reservation-test"
	DBURI      = "mongodb://root:example@localhost:27017"
)

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

type Map map[string]any

type Pagination struct {
	Limit int64
	Page  int64
}

type ResourceResponse struct {
	Results int `json:"results"`
	Data    any `json:"data"`
	Page    int `json:"page"`
}
