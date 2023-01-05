package main

import (
	"context"
	"github.com/anthdm/weavebox"
	"github.com/romankravchuk/toronto-bookings/api"
	"github.com/romankravchuk/toronto-bookings/store"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {
	app := weavebox.New()
	app.ErrorHandler = api.HandleAPIError
	adminRoute := app.Box("/admin")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	reservationStore := store.NewMongoReservationStore(client.Database("toronto-bookings"))
	reservationHandler := api.NewReservationHandler(reservationStore)

	// handle admin/reservation
	reservationRoute := adminRoute.Box("/reservation")
	reservationRoute.Get("/:id", reservationHandler.HandleGetReservationByID)
	reservationRoute.Get("", reservationHandler.HandleGetAllReservations)
	reservationRoute.Post("", reservationHandler.HandlePostReservation)

	log.Fatal(app.Serve(3000))
}
