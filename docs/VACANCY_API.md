# Vacancy API Documentation

## Описание
REST API для управления вакансиями в приложении поиска работы для студентов.

## Архитектура
Проект следует принципам **Clean Architecture** и **SOLID**:

- **Domain Layer** (`internal/domain/`): Сущности и интерфейсы репозиториев
- **Application Layer** (`internal/application/usecases/`): Бизнес-логика и use cases
- **Infrastructure Layer** (`internal/infrastructure/`): Реализация репозиториев (MongoDB)
- **Presentation Layer** (`internal/interfaces/http/handlers/`): HTTP handlers

## Сущность Vacancy

```json
{
  "id": "string (ObjectID)",
  "title": "string (required)",
  "type": "string (required)", // "Полная", "Частичная", "Стажировка"
  "format": "string (required)", // "Офис", "Удалённо", "Гибрид"
  "location": "string",
  "salary_type": "string (required)", // "range" или "fixed"
  "salary_from": "int (optional)",
  "salary_to": "int (optional)",
  "salary_fixed": "int (optional)",
  "skills": ["string"],
  "description": "string",
  "responsibilities": ["string"],
  "requirements": ["string"],
  "benefits": ["string"],
  "status": "string", // "Активна", "Приостановлена", "Закрыта"
  "responses_count": "int",
  "views_count": "int",
  "deadline": "timestamp",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

## API Endpoints

### 1. Создать вакансию
**POST** `/api/vacancies`

#### Request Body:
```json
{
  "title": "Frontend Developer",
  "type": "Полная",
  "format": "Гибрид",
  "location": "Алматы",
  "salary_type": "range",
  "salary_from": 200000,
  "salary_to": 350000,
  "skills": ["React", "TypeScript", "CSS"],
  "description": "Ищем опытного фронтенд разработчика",
  "responsibilities": [
    "Разработка пользовательских интерфейсов",
    "Оптимизация производительности"
  ],
  "requirements": [
    "Опыт работы с React от 2 лет",
    "Знание TypeScript"
  ],
  "benefits": [
    "Гибкий график",
    "ДМС"
  ],
  "deadline": "2024-12-31T00:00:00Z"
}
```

#### Response (201 Created):
```json
{
  "message": "vacancy created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "title": "Frontend Developer",
    ...
  }
}
```

#### Validation:
- `title`, `type`, `format`, `salary_type` - обязательны
- Если `salary_type = "range"`, то `salary_from` и `salary_to` обязательны
- Если `salary_type = "fixed"`, то `salary_fixed` обязателен
- `type` должен быть: "Полная", "Частичная" или "Стажировка"
- `format` должен быть: "Офис", "Удалённо" или "Гибрид"

---

### 2. Получить все вакансии
**GET** `/api/vacancies`

#### Query Parameters:
- `status` (optional): Фильтр по статусу ("Активна", "Приостановлена", "Закрыта")

#### Examples:
```
GET /api/vacancies
GET /api/vacancies?status=Активна
GET /api/vacancies?status=Приостановлена
```

#### Response (200 OK):
```json
{
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "title": "Frontend Developer",
      "type": "Полная",
      ...
    }
  ],
  "count": 1
}
```

---

### 3. Получить вакансию по ID
**GET** `/api/vacancies/:id`

#### Response (200 OK):
```json
{
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "title": "Frontend Developer",
    ...
  }
}
```

#### Response (404 Not Found):
```json
{
  "error": "vacancy not found"
}
```

---

### 4. Обновить вакансию
**PUT** `/api/vacancies/:id`

#### Request Body:
```json
{
  "title": "Senior Frontend Developer",
  "type": "Полная",
  "format": "Удалённо",
  "location": "Астана",
  "salary_type": "range",
  "salary_from": 300000,
  "salary_to": 500000,
  ...
}
```

#### Response (200 OK):
```json
{
  "message": "vacancy updated successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    ...
  }
}
```

---

### 5. Изменить статус вакансии
**PATCH** `/api/vacancies/:id/status`

#### Request Body:
```json
{
  "status": "Приостановлена"
}
```

#### Allowed statuses:
- "Активна"
- "Приостановлена"
- "Закрыта"

#### Response (200 OK):
```json
{
  "message": "vacancy status updated successfully",
  "status": "Приостановлена"
}
```

---

### 6. Удалить вакансию
**DELETE** `/api/vacancies/:id`

#### Response (200 OK):
```json
{
  "message": "vacancy deleted successfully"
}
```

## Error Responses

### 400 Bad Request
```json
{
  "error": "validation error message"
}
```

### 404 Not Found
```json
{
  "error": "vacancy not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "internal error message"
}
```

## Запуск проекта

1. Убедитесь, что `.env` файл настроен корректно
2. Установите зависимости:
   ```bash
   go mod download
   ```
3. Запустите сервер:
   ```bash
   go run main.go
   ```
4. Сервер будет доступен на `http://localhost:8081`

## Тестирование API

Используйте cURL, Postman или любой другой HTTP клиент для тестирования endpoints.

### Пример с cURL:

```bash
# Создать вакансию
curl -X POST http://localhost:8081/api/vacancies \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Backend Developer",
    "type": "Полная",
    "format": "Офис",
    "location": "Алматы",
    "salary_type": "fixed",
    "salary_fixed": 400000,
    "skills": ["Go", "MongoDB", "Docker"],
    "description": "Требуется опытный backend разработчик",
    "deadline": "2024-12-31T00:00:00Z"
  }'

# Получить все вакансии
curl http://localhost:8081/api/vacancies

# Получить активные вакансии
curl http://localhost:8081/api/vacancies?status=Активна

# Получить вакансию по ID
curl http://localhost:8081/api/vacancies/{id}

# Обновить статус
curl -X PATCH http://localhost:8081/api/vacancies/{id}/status \
  -H "Content-Type: application/json" \
  -d '{"status": "Приостановлена"}'

# Удалить вакансию
curl -X DELETE http://localhost:8081/api/vacancies/{id}
```

## Принципы SOLID в проекте

- **Single Responsibility**: Каждый компонент имеет одну ответственность
  - `VacancyRepository` - только операции с БД
  - `VacancyService` - только бизнес-логика
  - `VacancyHandler` - только обработка HTTP запросов

- **Open/Closed**: Легко расширяется без изменения существующего кода

- **Liskov Substitution**: Репозитории реализуют общий интерфейс

- **Interface Segregation**: Интерфейсы содержат только необходимые методы

- **Dependency Inversion**: Зависимости от абстракций (интерфейсов), а не от конкретных реализаций
