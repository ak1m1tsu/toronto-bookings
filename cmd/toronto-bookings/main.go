package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/romankravchuk/toronto-bookings/internal/config"
	"github.com/romankravchuk/toronto-bookings/internal/router"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	config, err := config.LoadConfig(os.Getenv("CONFIG_FILE_PATH"))
	if err != nil {
		log.Fatal(err)
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.MongoURI))
	db := client.Database("toronto-bookings")
	if err != nil {
		log.Fatal(err)
	}

	router := router.NewRouter(db, config)

	fmt.Println("app listening on localhost:" + config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}
