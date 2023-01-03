package main

import (
	"context"
	"fmt"
	"github.com/anthdm/weavebox"
	"github.com/romankravchuk/toronto-bookings/api"
	"github.com/romankravchuk/toronto-bookings/store"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func handleAPIError(ctx *weavebox.Context, err error) {
	fmt.Println(err)
}

func main() {
	app := weavebox.New()
	app.ErrorHandler = handleAPIError
	adminRoute := app.Box("/admin")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	reservationStore := store.NewMongoReservationStore(client.Database("toronto-bookings"))
	reservationHandler := api.NewReservationHandler(reservationStore)

	adminRoute.Get("/reservation/:id", reservationHandler.HandleGetReservationByID)
	adminRoute.Post("/reservation", reservationHandler.HandlePostReservation)

	log.Fatal(app.Serve(3000))
}
