package storage

import (
	"context"

	"github.com/romankravchuk/toronto-bookings/internal/storage/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserStore struct {
	db   *mongo.Database
	coll string
}

func NewMongoUserStore(db *mongo.Database) *MongoUserStore {
	return &MongoUserStore{
		db:   db,
		coll: "users",
	}
}

func (s *MongoUserStore) Insert(ctx context.Context, u *models.User) error {
	res, err := s.db.Collection(s.coll).InsertOne(ctx, u)
	u.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return err
}

func (s *MongoUserStore) GetByID(ctx context.Context, id string) (*models.User, error) {
	var (
		objID, _ = primitive.ObjectIDFromHex(id)
		user     = &models.User{}
		res      = s.db.Collection(s.coll).FindOne(ctx, bson.M{"_id": objID})
		err      = res.Decode(user)
	)
	return user, err
}

func (s *MongoUserStore) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var (
		user = &models.User{}
		res  = s.db.Collection(s.coll).FindOne(ctx, bson.M{"email": email})
		err  = res.Decode(user)
	)
	return user, err
}

func (s *MongoUserStore) GetAll(ctx context.Context) ([]*models.User, error) {
	cursor, err := s.db.Collection(s.coll).Find(ctx, map[string]any{})
	if err != nil {
		return nil, err
	}
	var user []*models.User
	err = cursor.All(ctx, &user)
	return user, err
}
