package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/romankravchuk/toronto-bookings/api"
	"github.com/romankravchuk/toronto-bookings/config"
	"github.com/romankravchuk/toronto-bookings/service"
	"github.com/romankravchuk/toronto-bookings/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.MongoURI))
	db := client.Database("toronto-bookings")
	if err != nil {
		log.Fatal(err)
	}

	r := router(db)

	fmt.Println("app listening on localhost:" + config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, r))
}

func router(db *mongo.Database) http.Handler {
	router := chi.NewRouter()
	router.Use(
		middleware.RequestID,
		middleware.Logger,
		api.Recover,
	)

	userStore := storage.NewMongoUserStore(db)
	userService := service.NewUserService(userStore)
	authMw := api.NewAuthMiddleware(userStore)

	// handle account auth
	authHandler := api.NewAuthenticationHandler(userService)
	router.Route("/account", func(r chi.Router) {
		r.Post("/sign-in", authHandler.HandleSignIn)
		r.Post("/sign-up", authHandler.HandleSignUp)
		r.Post("/logout", authHandler.HandleLogout)
		r.Post("/refresh-token", authHandler.HandleRefreshToken)
	})

	// handle admin
	reservationStore := storage.NewMongoReservationStore(db)
	reservationHandler := api.NewReservationHandler(reservationStore)
	reservationMw := api.NewReservationMiddleware(reservationStore)
	router.Route("/admin", func(r chi.Router) {
		r.Use(authMw.JWTRequired)
		r.Route("/reservation", func(r chi.Router) {
			r.Route("/{id}", func(r chi.Router) {
				r.Use(reservationMw.Context)
				r.Get("/", reservationHandler.HandleGetReservationById)
			})
			r.Get("/", reservationHandler.HandleGetReservations)
			r.Post("/", reservationHandler.HandlePostReservation)
		})
	})

	return router
}
