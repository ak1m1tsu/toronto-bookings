package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/romankravchuk/toronto-bookings/config"
	"github.com/romankravchuk/toronto-bookings/internal/router"
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

	router := router.NewRouter(db)

	fmt.Println("app listening on localhost:" + config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}
