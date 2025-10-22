package usecases

import (
	"context"
	"errors"

	"github.com/albkvv/student-job-finder-back/internal/domain/entities"
	"github.com/albkvv/student-job-finder-back/internal/domain/repositories"
)

type VacancyService struct {
	repo repositories.VacancyRepository
}

func NewVacancyService(repo repositories.VacancyRepository) *VacancyService {
	return &VacancyService{
		repo: repo,
	}
}

// CreateVacancy создает новую вакансию с валидацией
func (s *VacancyService) CreateVacancy(ctx context.Context, vacancy *entities.Vacancy) error {
	// Валидация обязательных полей
	if vacancy.Title == "" {
		return errors.New("title is required")
	}
	if vacancy.Type == "" {
		return errors.New("type is required")
	}
	if vacancy.Format == "" {
		return errors.New("format is required")
	}
	if vacancy.SalaryType == "" {
		return errors.New("salary_type is required")
	}

	// Валидация типа зарплаты
	switch vacancy.SalaryType {
	case entities.SalaryTypeRange:
		if vacancy.SalaryFrom == nil || vacancy.SalaryTo == nil {
			return errors.New("salary_from and salary_to are required when salary_type is 'range'")
		}
		if *vacancy.SalaryFrom > *vacancy.SalaryTo {
			return errors.New("salary_from cannot be greater than salary_to")
		}
	case entities.SalaryTypeFixed:
		if vacancy.SalaryFixed == nil {
			return errors.New("salary_fixed is required when salary_type is 'fixed'")
		}
	default:
		return errors.New("invalid salary_type, must be 'range' or 'fixed'")
	}

	// Валидация типа занятости
	if vacancy.Type != entities.VacancyTypeFull &&
		vacancy.Type != entities.VacancyTypePartial &&
		vacancy.Type != entities.VacancyTypeInternship {
		return errors.New("invalid type, must be 'Полная', 'Частичная', or 'Стажировка'")
	}

	// Валидация формата работы
	if vacancy.Format != entities.VacancyFormatOffice &&
		vacancy.Format != entities.VacancyFormatRemote &&
		vacancy.Format != entities.VacancyFormatHybrid {
		return errors.New("invalid format, must be 'Офис', 'Удалённо', or 'Гибрид'")
	}

	// Инициализация значений по умолчанию
	if vacancy.Status == "" {
		vacancy.Status = entities.VacancyStatusActive
	}
	if vacancy.Skills == nil {
		vacancy.Skills = []string{}
	}
	if vacancy.Responsibilities == nil {
		vacancy.Responsibilities = []string{}
	}
	if vacancy.Requirements == nil {
		vacancy.Requirements = []string{}
	}
	if vacancy.Benefits == nil {
		vacancy.Benefits = []string{}
	}

	return s.repo.Create(ctx, vacancy)
}

// GetVacancy получает вакансию по ID
func (s *VacancyService) GetVacancy(ctx context.Context, id string) (*entities.Vacancy, error) {
	vacancy, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if vacancy == nil {
		return nil, errors.New("vacancy not found")
	}
	return vacancy, nil
}

// GetAllVacancies получает все вакансии с фильтрацией по статусу
func (s *VacancyService) GetAllVacancies(ctx context.Context, status string) ([]*entities.Vacancy, error) {
	// Валидация статуса, если он указан
	if status != "" {
		if status != entities.VacancyStatusActive &&
			status != entities.VacancyStatusPaused &&
			status != entities.VacancyStatusClosed {
			return nil, errors.New("invalid status filter")
		}
	}

	return s.repo.FindAll(ctx, status)
}

// UpdateVacancy обновляет существующую вакансию
func (s *VacancyService) UpdateVacancy(ctx context.Context, vacancy *entities.Vacancy) error {
	// Проверка существования вакансии
	existing, err := s.repo.FindByID(ctx, vacancy.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("vacancy not found")
	}

	// Валидация обязательных полей
	if vacancy.Title == "" {
		return errors.New("title is required")
	}
	if vacancy.Type == "" {
		return errors.New("type is required")
	}
	if vacancy.Format == "" {
		return errors.New("format is required")
	}
	if vacancy.SalaryType == "" {
		return errors.New("salary_type is required")
	}

	// Валидация типа зарплаты
	switch vacancy.SalaryType {
	case entities.SalaryTypeRange:
		if vacancy.SalaryFrom == nil || vacancy.SalaryTo == nil {
			return errors.New("salary_from and salary_to are required when salary_type is 'range'")
		}
		if *vacancy.SalaryFrom > *vacancy.SalaryTo {
			return errors.New("salary_from cannot be greater than salary_to")
		}
	case entities.SalaryTypeFixed:
		if vacancy.SalaryFixed == nil {
			return errors.New("salary_fixed is required when salary_type is 'fixed'")
		}
	default:
		return errors.New("invalid salary_type, must be 'range' or 'fixed'")
	}

	// Валидация типа занятости
	if vacancy.Type != entities.VacancyTypeFull &&
		vacancy.Type != entities.VacancyTypePartial &&
		vacancy.Type != entities.VacancyTypeInternship {
		return errors.New("invalid type")
	}

	// Валидация формата работы
	if vacancy.Format != entities.VacancyFormatOffice &&
		vacancy.Format != entities.VacancyFormatRemote &&
		vacancy.Format != entities.VacancyFormatHybrid {
		return errors.New("invalid format")
	}

	return s.repo.Update(ctx, vacancy)
}

// UpdateVacancyStatus обновляет статус вакансии
func (s *VacancyService) UpdateVacancyStatus(ctx context.Context, id string, status string) error {
	// Валидация статуса
	if status != entities.VacancyStatusActive &&
		status != entities.VacancyStatusPaused &&
		status != entities.VacancyStatusClosed {
		return errors.New("invalid status, must be 'Активна', 'Приостановлена', or 'Закрыта'")
	}

	// Проверка существования вакансии
	vacancy, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if vacancy == nil {
		return errors.New("vacancy not found")
	}

	return s.repo.UpdateStatus(ctx, id, status)
}

// DeleteVacancy удаляет вакансию
func (s *VacancyService) DeleteVacancy(ctx context.Context, id string) error {
	// Проверка существования вакансии
	vacancy, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if vacancy == nil {
		return errors.New("vacancy not found")
	}

	return s.repo.Delete(ctx, id)
}
