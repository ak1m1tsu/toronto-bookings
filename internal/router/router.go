package router

import (
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/romankravchuk/toronto-bookings/internal/router/handlers"
	"github.com/romankravchuk/toronto-bookings/internal/service"
	"github.com/romankravchuk/toronto-bookings/internal/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

type Router struct {
	chi.Mux
	db *mongo.Database
}

func NewRouter(db *mongo.Database) *Router {
	router := chi.NewRouter()
	router.Use(
		chimiddleware.RequestID,
		chimiddleware.Logger,
		handlers.Recover,
	)

	userStore := storage.NewMongoUserStore(db)
	userService := service.NewUserService(userStore)
	authMw := handlers.NewAuthMiddleware(userStore)

	authHandler := handlers.NewAuthenticationHandler(userService)
	router.Route("/account", func(r chi.Router) {
		r.Post("/sign-in", authHandler.HandleSignIn)
		r.Post("/sign-up", authHandler.HandleSignUp)
		r.Post("/logout", authHandler.HandleLogout)
		r.Post("/refresh-token", authHandler.HandleRefreshToken)
	})

	reservationStore := storage.NewMongoReservationStore(db)
	reservationService := service.NewReservationService(reservationStore)
	reservationHandler := handlers.NewReservationHandler(reservationService)
	reservationMw := handlers.NewReservationMiddleware(reservationStore)
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

	return &Router{*router, db}
}
