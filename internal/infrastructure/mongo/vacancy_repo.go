package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/albkvv/student-job-finder-back/internal/domain/entities"
	"github.com/albkvv/student-job-finder-back/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoVacancyRepo struct {
	coll *mongo.Collection
}

func NewMongoVacancyRepo(coll *mongo.Collection) repositories.VacancyRepository {
	return &MongoVacancyRepo{
		coll: coll,
	}
}

func (r *MongoVacancyRepo) Create(ctx context.Context, vacancy *entities.Vacancy) error {
	// Генерируем ObjectID для новой вакансии
	objectID := primitive.NewObjectID()
	vacancy.ID = objectID.Hex()
	vacancy.CreatedAt = time.Now()
	vacancy.UpdatedAt = time.Now()
	
	if vacancy.Status == "" {
		vacancy.Status = entities.VacancyStatusActive
	}
	
	_, err := r.coll.InsertOne(ctx, vacancy)
	return err
}

func (r *MongoVacancyRepo) FindByID(ctx context.Context, id string) (*entities.Vacancy, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid vacancy ID format")
	}

	var vacancy entities.Vacancy
	err = r.coll.FindOne(ctx, bson.M{"_id": objectID}).Decode(&vacancy)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &vacancy, nil
}

func (r *MongoVacancyRepo) FindAll(ctx context.Context, status string) ([]*entities.Vacancy, error) {
	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}

	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var vacancies []*entities.Vacancy
	if err := cursor.All(ctx, &vacancies); err != nil {
		return nil, err
	}

	return vacancies, nil
}

func (r *MongoVacancyRepo) Update(ctx context.Context, vacancy *entities.Vacancy) error {
	objectID, err := primitive.ObjectIDFromHex(vacancy.ID)
	if err != nil {
		return errors.New("invalid vacancy ID format")
	}

	vacancy.UpdatedAt = time.Now()
	
	update := bson.M{
		"$set": bson.M{
			"title":            vacancy.Title,
			"type":             vacancy.Type,
			"format":           vacancy.Format,
			"location":         vacancy.Location,
			"salary_type":      vacancy.SalaryType,
			"salary_from":      vacancy.SalaryFrom,
			"salary_to":        vacancy.SalaryTo,
			"salary_fixed":     vacancy.SalaryFixed,
			"skills":           vacancy.Skills,
			"description":      vacancy.Description,
			"responsibilities": vacancy.Responsibilities,
			"requirements":     vacancy.Requirements,
			"benefits":         vacancy.Benefits,
			"status":           vacancy.Status,
			"deadline":         vacancy.Deadline,
			"updated_at":       vacancy.UpdatedAt,
		},
	}

	result, err := r.coll.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}
	
	if result.MatchedCount == 0 {
		return errors.New("vacancy not found")
	}

	return nil
}

func (r *MongoVacancyRepo) UpdateStatus(ctx context.Context, id string, status string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid vacancy ID format")
	}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := r.coll.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}
	
	if result.MatchedCount == 0 {
		return errors.New("vacancy not found")
	}

	return nil
}

func (r *MongoVacancyRepo) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid vacancy ID format")
	}

	result, err := r.coll.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	
	if result.DeletedCount == 0 {
		return errors.New("vacancy not found")
	}

	return nil
}

func (r *MongoVacancyRepo) IncrementViews(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid vacancy ID format")
	}

	update := bson.M{
		"$inc": bson.M{
			"views_count": 1,
		},
	}

	_, err = r.coll.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (r *MongoVacancyRepo) IncrementResponses(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid vacancy ID format")
	}

	update := bson.M{
		"$inc": bson.M{
			"responses_count": 1,
		},
	}

	_, err = r.coll.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}
