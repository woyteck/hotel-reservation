package main

import (
	"context"
	"fmt"
	"hotel-reservation/api"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	hotelStore   db.HotelStore
	roomStore    db.RoomStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          = context.Background()
)

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
	bookingStore = db.NewMongoBookingStore(client)
}

func main() {
	james := seedUser("James", "Foo", "james@foo.com", "supersecurepassword", false)
	seedUser("admin", "admin", "admin@admin.com", "admin", true)

	seedHotel("Bellucia", "France", 3)
	seedHotel("The Cozy Hotel", "The Netherlands", 4)
	hotel := seedHotel("Don't die in your sleep", "London", 1)

	seedRoom("small", true, 98.99, hotel.ID)
	seedRoom("medium", true, 18.99, hotel.ID)
	room := seedRoom("large", false, 28.99, hotel.ID)

	seedBooking(james.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 2), 2)
}

func seedUser(firstName, lastName, email, password string, isAdmin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err)
	}

	user.IsAdmin = isAdmin

	insertedUser, err := userStore.InsertUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))

	return insertedUser
}

func seedHotel(name string, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertedHotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func seedRoom(size string, seaside bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: seaside,
		Price:   price,
		HotelID: hotelID,
	}

	insertedRoom, err := roomStore.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func seedBooking(userID, roomID primitive.ObjectID, from, till time.Time, numPersons int) *types.Booking {
	booking := &types.Booking{
		UserID:     userID,
		RoomID:     roomID,
		FromDate:   from,
		TillDate:   till,
		NumPersons: numPersons,
	}

	insertedBooking, err := bookingStore.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("booking:", insertedBooking.ID)

	return insertedBooking
}

// func seedHotel(name string, location string, rating int) {
// 	hotel := types.Hotel{
// 		Name:     name,
// 		Location: location,
// 		Rooms:    []primitive.ObjectID{},
// 		Rating:   rating,
// 	}
// 	rooms := []types.Room{
// 		{
// 			Size:  "small",
// 			Price: 99.9,
// 		},
// 		{
// 			Size:  "normal",
// 			Price: 122.9,
// 		},
// 		{
// 			Size:  "kingsize",
// 			Price: 222.9,
// 		},
// 	}

// 	insertedHotel, err := hotelStore.Insert(ctx, &hotel)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for _, room := range rooms {
// 		room.HotelID = insertedHotel.ID
// 		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("room -> %s\n", insertedRoom.ID)
// 	}
// }
