package storage

import (
	"context"

	"github.com/romankravchuk/toronto-bookings/internal/storage/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoReservationStore struct {
	db   *mongo.Database
	coll string
}

func NewMongoReservationStore(db *mongo.Database) *MongoReservationStore {
	return &MongoReservationStore{
		db:   db,
		coll: "reservations",
	}
}

func (s *MongoReservationStore) Insert(ctx context.Context, r *models.Reservation) error {
	res, err := s.db.Collection(s.coll).InsertOne(ctx, r)
	r.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return err
}

func (s *MongoReservationStore) GetByID(ctx context.Context, id string) (*models.Reservation, error) {
	var (
		objID, _    = primitive.ObjectIDFromHex(id)
		reservation = &models.Reservation{}
		res         = s.db.Collection(s.coll).FindOne(ctx, bson.M{"_id": objID})
		err         = res.Decode(reservation)
	)
	return reservation, err
}

func (s *MongoReservationStore) GetAll(ctx context.Context) ([]*models.Reservation, error) {
	cursor, err := s.db.Collection(s.coll).Find(ctx, map[string]any{})
	if err != nil {
		return nil, err
	}
	var reservations []*models.Reservation
	err = cursor.All(ctx, &reservations)
	return reservations, err
}
