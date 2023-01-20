package storage

import (
	"context"
	"github.com/romankravchuk/toronto-bookings/types"
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

func (s *MongoUserStore) Insert(ctx context.Context, u *types.User) error {
	res, err := s.db.Collection(s.coll).InsertOne(ctx, u)
	u.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return err
}

func (s *MongoUserStore) GetByID(ctx context.Context, id string) (*types.User, error) {
	var (
		objID, _ = primitive.ObjectIDFromHex(id)
		user     = &types.User{}
		res      = s.db.Collection(s.coll).FindOne(ctx, bson.M{"_id": objID})
		err      = res.Decode(user)
	)
	return user, err
}

func (s *MongoUserStore) GetByEmail(ctx context.Context, email string) (*types.User, error) {
	var (
		user = &types.User{}
		res  = s.db.Collection(s.coll).FindOne(ctx, bson.M{"email": email})
		err  = res.Decode(user)
	)
	return user, err
}

func (s *MongoUserStore) GetAll(ctx context.Context) ([]*types.User, error) {
	cursor, err := s.db.Collection(s.coll).Find(ctx, map[string]any{})
	if err != nil {
		return nil, err
	}
	var user []*types.User
	err = cursor.All(ctx, &user)
	return user, err
}
