package mongo

import (
	"context"
	"errors"

	"github.com/albkvv/student-job-finder-back/internal/domain/entities"
	"github.com/albkvv/student-job-finder-back/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepo struct {
	coll *mongo.Collection
}

func NewMongoUserRepo(coll *mongo.Collection) repositories.UserRepository {
	return &MongoUserRepo{
		coll: coll,
	}
}

func (r *MongoUserRepo) Create(ctx context.Context, user *entities.User) error {
	_, err := r.coll.InsertOne(ctx, user)
	return err
}

func (r *MongoUserRepo) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	err := r.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *MongoUserRepo) FindByPhone(ctx context.Context, phone string) (*entities.User, error) {
	var user entities.User
	err := r.coll.FindOne(ctx, bson.M{"phone": phone}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *MongoUserRepo) FindByID(ctx context.Context, id string) (*entities.User, error) {
	var user entities.User
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *MongoUserRepo) Update(ctx context.Context, user *entities.User) error {
	update := bson.M{"$set": user}
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": user.ID}, update)
	return err
}
