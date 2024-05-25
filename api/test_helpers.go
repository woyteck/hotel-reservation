package api

import (
	"context"
	"hotel-reservation/db"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	client *mongo.Client
	*db.Store
}

func setup(t *testing.T) *testdb {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGO_TEST_DB_URL")))
	if err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)

	return &testdb{
		client: client,
		Store: &db.Store{
			Hotel:   hotelStore,
			User:    db.NewMongoUserStore(client),
			Room:    db.NewMongoRoomStore(client, hotelStore),
			Booking: db.NewMongoBookingStore(client),
		},
	}
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.client.Database(os.Getenv("MONGO_TEST_DB_NAME")).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}
