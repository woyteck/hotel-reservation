package main

import (
	"flag"
	"hotel-reservation/api"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	listernAddr := flag.String("listenAddr", ":5000", "The listener address of the API server")
	flag.Parse()

	app := fiber.New()

	apiv1 := app.Group("/api/v1")
	apiv1.Get("/users", api.HandleGetUsers)
	apiv1.Get("/user/:id", api.HandleGetUser)

	log.Fatal(app.Listen(*listernAddr))
}
