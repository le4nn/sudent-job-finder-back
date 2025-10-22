#!/bin/bash

# Цвета для вывода
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8081/api"

echo -e "${BLUE}=== Тестирование Vacancy API ===${NC}\n"

# 1. Проверка health endpoint
echo -e "${BLUE}1. Проверка health endpoint${NC}"
curl -s -X GET "$BASE_URL/health" | jq '.'
echo -e "\n"

# 2. Создание вакансии с диапазоном зарплаты
echo -e "${BLUE}2. Создание вакансии (диапазон зарплаты)${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/vacancies" \
  -H "Content-Type: application/json" \
  -d '{
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
  }')

echo "$RESPONSE" | jq '.'
VACANCY_ID=$(echo "$RESPONSE" | jq -r '.data.id')
echo -e "${GREEN}Создана вакансия с ID: $VACANCY_ID${NC}\n"

# 3. Создание вакансии с фиксированной зарплатой
echo -e "${BLUE}3. Создание вакансии (фиксированная зарплата)${NC}"
RESPONSE2=$(curl -s -X POST "$BASE_URL/vacancies" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Backend Developer",
    "type": "Полная",
    "format": "Офис",
    "location": "Астана",
    "salary_type": "fixed",
    "salary_fixed": 400000,
    "skills": ["Go", "MongoDB", "Docker"],
    "description": "Требуется опытный backend разработчик",
    "responsibilities": [
      "Разработка API",
      "Работа с базами данных"
    ],
    "requirements": [
      "Опыт работы с Go от 1 года",
      "Знание MongoDB"
    ],
    "benefits": [
      "Корпоративное обучение",
      "Премии"
    ],
    "deadline": "2024-12-31T00:00:00Z"
  }')

echo "$RESPONSE2" | jq '.'
VACANCY_ID_2=$(echo "$RESPONSE2" | jq -r '.data.id')
echo -e "${GREEN}Создана вакансия с ID: $VACANCY_ID_2${NC}\n"

# 4. Получение всех вакансий
echo -e "${BLUE}4. Получение всех вакансий${NC}"
curl -s -X GET "$BASE_URL/vacancies" | jq '.'
echo -e "\n"

# 5. Получение активных вакансий
echo -e "${BLUE}5. Получение активных вакансий${NC}"
curl -s -X GET "$BASE_URL/vacancies?status=Активна" | jq '.'
echo -e "\n"

# 6. Получение конкретной вакансии по ID
if [ ! -z "$VACANCY_ID" ] && [ "$VACANCY_ID" != "null" ]; then
  echo -e "${BLUE}6. Получение вакансии по ID: $VACANCY_ID${NC}"
  curl -s -X GET "$BASE_URL/vacancies/$VACANCY_ID" | jq '.'
  echo -e "\n"

  # 7. Обновление вакансии
  echo -e "${BLUE}7. Обновление вакансии $VACANCY_ID${NC}"
  curl -s -X PUT "$BASE_URL/vacancies/$VACANCY_ID" \
    -H "Content-Type: application/json" \
    -d '{
      "title": "Senior Frontend Developer",
      "type": "Полная",
      "format": "Удалённо",
      "location": "Алматы",
      "salary_type": "range",
      "salary_from": 300000,
      "salary_to": 500000,
      "skills": ["React", "TypeScript", "Next.js"],
      "description": "Ищем Senior фронтенд разработчика",
      "responsibilities": [
        "Разработка сложных интерфейсов",
        "Менторство джуниоров"
      ],
      "requirements": [
        "Опыт работы с React от 3 лет",
        "Опыт с Next.js"
      ],
      "benefits": [
        "Полностью удаленная работа",
        "ДМС"
      ],
      "status": "Активна",
      "deadline": "2024-12-31T00:00:00Z"
    }' | jq '.'
  echo -e "\n"

  # 8. Изменение статуса вакансии
  echo -e "${BLUE}8. Изменение статуса вакансии на 'Приостановлена'${NC}"
  curl -s -X PATCH "$BASE_URL/vacancies/$VACANCY_ID/status" \
    -H "Content-Type: application/json" \
    -d '{"status": "Приостановлена"}' | jq '.'
  echo -e "\n"

  # 9. Проверка изменения статуса
  echo -e "${BLUE}9. Проверка изменения статуса${NC}"
  curl -s -X GET "$BASE_URL/vacancies/$VACANCY_ID" | jq '.data.status'
  echo -e "\n"

  # 10. Удаление вакансии
  echo -e "${BLUE}10. Удаление вакансии $VACANCY_ID${NC}"
  curl -s -X DELETE "$BASE_URL/vacancies/$VACANCY_ID" | jq '.'
  echo -e "\n"

  # 11. Проверка удаления
  echo -e "${BLUE}11. Проверка удаления (должно вернуть 404)${NC}"
  curl -s -X GET "$BASE_URL/vacancies/$VACANCY_ID" | jq '.'
  echo -e "\n"
fi

# 12. Тест валидации - создание без обязательных полей
echo -e "${BLUE}12. Тест валидации - создание без title (должна быть ошибка)${NC}"
curl -s -X POST "$BASE_URL/vacancies" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "Полная",
    "format": "Офис",
    "salary_type": "fixed",
    "salary_fixed": 300000
  }' | jq '.'
echo -e "\n"

# 13. Тест валидации - неверный salary_type
echo -e "${BLUE}13. Тест валидации - range без salary_from и salary_to (должна быть ошибка)${NC}"
curl -s -X POST "$BASE_URL/vacancies" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Vacancy",
    "type": "Полная",
    "format": "Офис",
    "salary_type": "range"
  }' | jq '.'
echo -e "\n"

echo -e "${GREEN}=== Тестирование завершено ===${NC}"
