package repositories

import "github.com/albkvv/student-job-finder-back/internal/domain/entities"

type UserRepository interface {
    FindByPhone(phone string) (*entities.User, error)
    Create(user *entities.User) error
}

type VerificationCodeRepository interface {
    SetCode(phone, code string, expiresAt int64) error
    GetCode(phone string) (*entities.VerificationCode, error)
    DeleteCode(phone string) error
}

