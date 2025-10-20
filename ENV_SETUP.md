# Настройка переменных окружения

## MongoDB Atlas

Для подключения к MongoDB Atlas создайте файл `.env` в корне проекта со следующим содержимым:

```env
# MongoDB Atlas Configuration
MONGODB_URI=mongodb+srv://alibekdias36_db_user:di%40s_o5@cluster0.lkhseyf.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0
MONGODB_DB=student_job_finder

# Server Configuration
PORT=8080
```

## Альтернативный способ

Вы также можете установить переменные окружения напрямую в системе:

### Linux/macOS:
```bash
export MONGODB_URI="mongodb+srv://alibekdias36_db_user:di%40s_o5@cluster0.lkhseyf.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
export MONGODB_DB="student_job_finder"
export PORT="8080"
```

### Windows:
```cmd
set MONGODB_URI=mongodb+srv://alibekdias36_db_user:di%40s_o5@cluster0.lkhseyf.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0
set MONGODB_DB=student_job_finder
set PORT=8080
```

## Запуск приложения

После настройки переменных окружения запустите приложение:

```bash
go run main.go
```

Приложение будет доступно по адресу: http://localhost:8080
