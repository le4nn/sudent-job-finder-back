package repositories

import (
    "context"
    "github.com/albkvv/student-job-finder-back/internal/domain/entities"
)

type UserRepository interface {
    Create(ctx context.Context, user *entities.User) error
    FindByEmail(ctx context.Context, email string) (*entities.User, error)
    FindByPhone(ctx context.Context, phone string) (*entities.User, error)
    FindByID(ctx context.Context, id string) (*entities.User, error)
    Update(ctx context.Context, user *entities.User) error
}

type VerificationCodeRepository interface {
    SetCode(ctx context.Context, identifier, code, codeType string, expiresAt int64) error
    GetCode(ctx context.Context, identifier, codeType string) (*entities.VerificationCode, error)
    IncrementAttempts(ctx context.Context, identifier, codeType string) error
    DeleteCode(ctx context.Context, identifier, codeType string) error
}

