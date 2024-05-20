package main

import (
	"context"
	"flag"
	"hotel-reservation/api"
	"hotel-reservation/db"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listernAddr := flag.String("listenAddr", ":5000", "The listener address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	var (
		userHandler  = api.NewUserHandler(db.NewMongoUserStore(client, db.DBNAME))
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		hotelHandler = api.NewHotelHandler(hotelStore, roomStore)
		app          = fiber.New(config)
		apiv1        = app.Group("/api/v1")
	)

	//user
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/users/:id", userHandler.HandleGetUser)
	apiv1.Post("/users", userHandler.HandlePostUser)
	apiv1.Put("/users/:id", userHandler.HandlePutUser)
	apiv1.Delete("/users/:id", userHandler.HandleDeleteUser)

	//hotel
	apiv1.Get("/hotels", hotelHandler.HandleGetHotels)

	log.Fatal(app.Listen(*listernAddr))
}
