package entities

import "time"

type Vacancy struct {
	ID              string    `json:"id" bson:"_id,omitempty"`
	Title           string    `json:"title" bson:"title"`
	Type            string    `json:"type" bson:"type"` // "Полная", "Частичная", "Стажировка"
	Format          string    `json:"format" bson:"format"` // "Офис", "Удалённо", "Гибрид"
	Location        string    `json:"location" bson:"location"`
	SalaryType      string    `json:"salary_type" bson:"salary_type"` // "range", "fixed"
	SalaryFrom      *int      `json:"salary_from,omitempty" bson:"salary_from,omitempty"`
	SalaryTo        *int      `json:"salary_to,omitempty" bson:"salary_to,omitempty"`
	SalaryFixed     *int      `json:"salary_fixed,omitempty" bson:"salary_fixed,omitempty"`
	Skills          []string  `json:"skills" bson:"skills"`
	Description     string    `json:"description" bson:"description"`
	Responsibilities []string `json:"responsibilities" bson:"responsibilities"`
	Requirements    []string  `json:"requirements" bson:"requirements"`
	Benefits        []string  `json:"benefits" bson:"benefits"`
	Status          string    `json:"status" bson:"status"` // "Активна", "Приостановлена", "Закрыта"
	ResponsesCount  int       `json:"responses_count" bson:"responses_count"`
	ViewsCount      int       `json:"views_count" bson:"views_count"`
	Deadline        time.Time `json:"deadline" bson:"deadline"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
}

// VacancyStatus константы для статусов вакансий
const (
	VacancyStatusActive    = "Активна"
	VacancyStatusPaused    = "Приостановлена"
	VacancyStatusClosed    = "Закрыта"
)

// VacancyType константы для типов занятости
const (
	VacancyTypeFull       = "Полная"
	VacancyTypePartial    = "Частичная"
	VacancyTypeInternship = "Стажировка"
)

// VacancyFormat константы для форматов работы
const (
	VacancyFormatOffice  = "Офис"
	VacancyFormatRemote  = "Удалённо"
	VacancyFormatHybrid  = "Гибрид"
)

// SalaryType константы для типов зарплаты
const (
	SalaryTypeRange = "range"
	SalaryTypeFixed = "fixed"
)
