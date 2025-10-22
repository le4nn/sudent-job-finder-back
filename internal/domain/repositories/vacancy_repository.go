package repositories

import (
	"context"
	"github.com/albkvv/student-job-finder-back/internal/domain/entities"
)

type VacancyRepository interface {
	Create(ctx context.Context, vacancy *entities.Vacancy) error
	FindByID(ctx context.Context, id string) (*entities.Vacancy, error)
	FindAll(ctx context.Context, status string) ([]*entities.Vacancy, error)
	Update(ctx context.Context, vacancy *entities.Vacancy) error
	UpdateStatus(ctx context.Context, id string, status string) error
	Delete(ctx context.Context, id string) error
	IncrementViews(ctx context.Context, id string) error
	IncrementResponses(ctx context.Context, id string) error
}
