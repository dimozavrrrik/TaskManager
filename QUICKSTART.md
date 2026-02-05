# Быстрый старт TaskManager

## Запуск через Docker Compose (рекомендуется)

### 1. Предварительные требования

- Docker и Docker Compose установлены
- Порты 3000, 8080, 5432 свободны

### 2. Клонируйте и настройте

```bash
# Клонируйте репозиторий
git clone <repository-url>
cd TaskManager

# Скопируйте и настройте переменные окружения
cp .env.example .env

# ВАЖНО: Измените JWT_SECRET в .env файле на случайную строку (минимум 32 символа)
# Например: JWT_SECRET=my-super-secret-key-12345678901234567890
```

### 3. Запустите все сервисы

```bash
docker-compose up -d
```

Это запустит:
- PostgreSQL (порт 5432)
- Backend API (порт 8080)
- Frontend (порт 3000)

### 4. Проверьте что все работает

```bash
# Проверьте статус сервисов
docker-compose ps

# Проверьте Backend API
curl http://localhost:8080/api/v1/health

# Откройте Frontend в браузере
# http://localhost:3000
```

### 5. Начните работу

1. Откройте http://localhost:3000 в браузере
2. Нажмите "Зарегистрироваться"
3. Заполните форму регистрации
4. Войдите в систему
5. Создайте свою первую задачу!

## Полезные команды

```bash
# Просмотр логов
docker-compose logs -f

# Просмотр логов конкретного сервиса
docker-compose logs -f frontend
docker-compose logs -f api

# Остановка сервисов
docker-compose down

# Остановка с удалением данных БД
docker-compose down -v

# Перезапуск сервисов
docker-compose restart

# Пересборка образов
docker-compose up -d --build
```

## Устранение проблем

### Frontend не загружается

```bash
# Проверьте логи
docker-compose logs frontend

# Убедитесь что порт 3000 свободен
netstat -an | grep 3000  # Linux/Mac
netstat -an | findstr 3000  # Windows

# Перезапустите сервис
docker-compose restart frontend
```

### API не отвечает

```bash
# Проверьте что JWT_SECRET установлен
docker-compose config | grep JWT_SECRET

# Проверьте логи
docker-compose logs api

# Перезапустите API
docker-compose restart api
```

### База данных не подключается

```bash
# Проверьте что PostgreSQL запущен
docker-compose ps postgres

# Проверьте логи
docker-compose logs postgres

# Подключитесь к БД вручную
docker exec -it taskmanager_postgres psql -U taskmanager -d taskmanager
```

## Структура портов

| Сервис | Порт | URL |
|--------|------|-----|
| Frontend | 3000 | http://localhost:3000 |
| Backend API | 8080 | http://localhost:8080 |
| PostgreSQL | 5432 | localhost:5432 |
| PgAdmin (dev) | 5050 | http://localhost:5050 |

## Следующие шаги

1. Ознакомьтесь с полной документацией в [README.md](README.md)
2. Изучите [API документацию](README.md#документация-api)
3. Настройте [production deployment](README.md#развертывание)

## Разработка

### Локальный запуск Backend

```bash
# Запустите только PostgreSQL
docker-compose up postgres -d

# Установите переменные окружения
export DATABASE_URL="postgres://taskmanager:taskmanager_password@localhost:5432/taskmanager?sslmode=disable"
export JWT_SECRET="dev-secret-key-min-32-chars-long-for-testing-only"
export SERVER_ADDRESS=":8080"

# Запустите API
go run cmd/api/main.go
```

### Локальный запуск Frontend

```bash
cd TaskManager.Client
dotnet run
```

Frontend будет доступен на http://localhost:5000
