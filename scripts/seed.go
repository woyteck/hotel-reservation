package main

import (
	"context"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	hotelStore db.HotelStore
	roomStore  db.RoomStore
	userStore  db.UserStore
	ctx        = context.Background()
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
}

func main() {
	seedHotel("Bellucia", "France", 3)
	seedHotel("The Cozy Hotel", "The Netherlands", 4)
	seedHotel("Don't die in your sleep", "London", 1)
	seedUser("James", "Foo", "james@foo.com", "supersecurepassword", false)
	seedUser("admin", "admin", "admin@admin.com", "admin", true)
}

func seedUser(firstName, lastName, email, password string, isAdmin bool) {
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

	_, err = userStore.InsertUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
}

func seedHotel(name string, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}
	rooms := []types.Room{
		{
			Size:  "small",
			Price: 99.9,
		},
		{
			Size:  "normal",
			Price: 122.9,
		},
		{
			Size:  "kingsize",
			Price: 222.9,
		},
	}

	insertedHotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}
}
