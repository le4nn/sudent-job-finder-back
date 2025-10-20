package mongo

import (
	"context"
	"github.com/albkvv/student-job-finder-back/internal/domain/entities"
	"github.com/albkvv/student-job-finder-back/internal/domain/repositories"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepo struct {
	coll *mongo.Collection
}

func NewMongoUserRepo(coll *mongo.Collection) repositories.UserRepository {
	return &MongoUserRepo{coll: coll}
}

func (r *MongoUserRepo) FindByPhone(phone string) (*entities.User, error) {
	ctx := context.Background()
	var m struct {
		ID    primitive.ObjectID `bson:"_id"`
		Phone string            `bson:"phone"`
		Name  string            `bson:"name"`
		Role  string            `bson:"role"`
	}
	if err := r.coll.FindOne(ctx, bson.M{"phone": phone}).Decode(&m); err != nil {
		return nil, err
	}
	return &entities.User{
		ID:    m.ID.Hex(),
		Phone: m.Phone,
		Name:  m.Name,
		Role:  m.Role,
	}, nil
}

func (r *MongoUserRepo) Create(u *entities.User) error {
	ctx := context.Background()
	insert := bson.M{
		"phone": u.Phone,
		"name":  u.Name,
		"role":  u.Role,
	}
	res, err := r.coll.InsertOne(ctx, insert)
	if err != nil {
		return err
	}
	id, ok := res.InsertedID.(primitive.ObjectID)
	if ok {
		u.ID = id.Hex()
	}
	return nil
}
