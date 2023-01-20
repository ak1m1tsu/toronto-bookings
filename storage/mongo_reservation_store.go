package storage

import (
	"context"
	"github.com/romankravchuk/toronto-bookings/types"
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

func (s *MongoReservationStore) Insert(ctx context.Context, r *types.Reservation) error {
	res, err := s.db.Collection(s.coll).InsertOne(ctx, r)
	r.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return err
}

func (s *MongoReservationStore) GetByID(ctx context.Context, id string) (*types.Reservation, error) {
	var (
		objID, _    = primitive.ObjectIDFromHex(id)
		reservation = &types.Reservation{}
		res         = s.db.Collection(s.coll).FindOne(ctx, bson.M{"_id": objID})
		err         = res.Decode(reservation)
	)
	return reservation, err
}

func (s *MongoReservationStore) GetAll(ctx context.Context) ([]*types.Reservation, error) {
	cursor, err := s.db.Collection(s.coll).Find(ctx, map[string]any{})
	if err != nil {
		return nil, err
	}
	var reservations []*types.Reservation
	err = cursor.All(ctx, &reservations)
	return reservations, err
}
