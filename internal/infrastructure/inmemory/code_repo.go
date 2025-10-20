package inmemory

import (
	"errors"
	"sync"

	"github.com/albkvv/student-job-finder-back/internal/domain/entities"
	"github.com/albkvv/student-job-finder-back/internal/domain/repositories"
)

type InMemoryCodeRepo struct {
	mu   sync.Mutex
	data map[string]*entities.VerificationCode
}

func NewInMemoryCodeRepo() repositories.VerificationCodeRepository {
	return &InMemoryCodeRepo{data: make(map[string]*entities.VerificationCode)}
}

func (r *InMemoryCodeRepo) SetCode(phone, code string, expiresAt int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[phone] = &entities.VerificationCode{Phone: phone, Code: code, ExpiresAt: expiresAt}
	return nil
}
func (r *InMemoryCodeRepo) GetCode(phone string) (*entities.VerificationCode, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	vc, ok := r.data[phone]
	if !ok {
		return nil, errors.New("not found")
	}
	return vc, nil
}
func (r *InMemoryCodeRepo) DeleteCode(phone string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.data, phone)
	return nil
}
