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

	authMw := api.AdminAuthMiddleware{}
	adminRoute := app.Box("/admin")
	adminRoute.Use(authMw.Authenticate)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	db := client.Database("toronto-bookings")
	if err != nil {
		log.Fatal(err)
	}
	reservationStore := store.NewMongoReservationStore(db)
	reservationHandler := api.NewReservationHandler(reservationStore)

	userStore := store.NewMongoUserStore(db)
	authHandler := api.NewAuthenticationHandler(userStore)

	// handle auth
	authRoute := app.Box("/account")
	authRoute.Post("/sign-in", authHandler.HandleSignIn)
	authRoute.Post("/sign-up", authHandler.HandleSignUp)

	// handle admin/reservation
	adminReservationRoute := adminRoute.Box("/reservation")
	adminReservationRoute.Get("/:id", reservationHandler.HandleGetReservationByID)
	adminReservationRoute.Get("/", reservationHandler.HandleGetAllReservations)
	adminReservationRoute.Post("/", reservationHandler.HandlePostReservation)

	log.Fatal(app.Serve(3000))
}	
