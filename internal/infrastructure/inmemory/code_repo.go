package inmemory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/albkvv/student-job-finder-back/internal/domain/entities"
	"github.com/albkvv/student-job-finder-back/internal/domain/repositories"
)

type InMemoryCodeRepo struct {
	mu    sync.RWMutex
	codes map[string]*entities.VerificationCode
}

func NewInMemoryCodeRepo() repositories.VerificationCodeRepository {
	return &InMemoryCodeRepo{
		codes: make(map[string]*entities.VerificationCode),
	}
}

func (r *InMemoryCodeRepo) SetCode(ctx context.Context, identifier, code, codeType string, expiresAt int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := fmt.Sprintf("%s:%s", identifier, codeType)
	r.codes[key] = &entities.VerificationCode{
		Identifier: identifier,
		Code:       code,
		Type:       codeType,
		Attempts:   0,
		ExpiresAt:  expiresAt,
	}
	return nil
}

func (r *InMemoryCodeRepo) GetCode(ctx context.Context, identifier, codeType string) (*entities.VerificationCode, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := fmt.Sprintf("%s:%s", identifier, codeType)
	vc, ok := r.codes[key]
	if !ok {
		return nil, nil
	}
	if time.Now().Unix() > vc.ExpiresAt {
		delete(r.codes, key)
		return nil, nil
	}
	return vc, nil
}

func (r *InMemoryCodeRepo) IncrementAttempts(ctx context.Context, identifier, codeType string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := fmt.Sprintf("%s:%s", identifier, codeType)
	if vc, ok := r.codes[key]; ok {
		vc.Attempts++
	}
	return nil
}

func (r *InMemoryCodeRepo) DeleteCode(ctx context.Context, identifier, codeType string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := fmt.Sprintf("%s:%s", identifier, codeType)
	delete(r.codes, key)
	return nil
}
