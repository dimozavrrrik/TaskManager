# CRM-сервис для управления задачами команды

Stateless микросервис для управления задачами команд до 200 человек, построенный на Go и PostgreSQL.

## Возможности

- **Управление сотрудниками**: CRUD операции для профилей сотрудников
- **Управление задачами**: Полный жизненный цикл задачи с несколькими статусами
- **Назначение ролей**: Роли исполнителя, ответственного и заказчика
- **Учет времени**: Несколько записей времени на задачу с автоматическим суммированием
- **Коммуникация в задачах**: Внутренняя переписка с автоматическими системными сообщениями
- **Архивация задач**: Архивирование завершенных задач
- **Отслеживание статусов**: Автоматические системные сообщения при смене статуса

## Статусы задач

- Новый (new)
- В работе (in_progress)
- На код ревью (code_review)
- На тестировании (testing)
- Возвращено с ошибкой (returned_with_errors)
- Закрыта (closed)

## Технологический стек

### Backend
- **Язык**: Go 1.22
- **База данных**: PostgreSQL 16
- **Router**: Gorilla Mux
- **Аутентификация**: JWT (golang-jwt/jwt/v5)
- **Хеширование паролей**: bcrypt
- **Валидация**: go-playground/validator
- **Архитектура**: Микросервис, Stateless, RESTful API

### Frontend
- **Фреймворк**: Blazor WebAssembly (.NET 8)
- **UI библиотека**: MudBlazor (Material Design)
- **Хранение токенов**: Blazored.LocalStorage
- **HTTP клиент**: System.Net.Http
- **PWA**: Service Worker поддержка

### Инфраструктура
- **Контейнеризация**: Docker, Docker Compose
- **Web сервер**: Nginx (для frontend)
- **Reverse Proxy**: Nginx (API routing)

## Структура проекта

```
TaskManager/
├── cmd/api/                      # Точка входа приложения
│   └── main.go
├── internal/                     # Приватный код приложения
│   ├── config/                   # Управление конфигурацией
│   ├── database/                 # Подключение к БД и миграции
│   ├── domain/                   # Domain-сущности
│   ├── dto/                      # Data Transfer Objects
│   ├── handler/                  # HTTP handlers
│   ├── middleware/               # HTTP middleware
│   ├── repository/               # Слой доступа к данным
│   ├── router/                   # HTTP routing
│   └── service/                  # Слой бизнес-логики
├── pkg/                          # Переиспользуемые пакеты
│   ├── errors/                   # Пользовательские типы ошибок
│   ├── logger/                   # Структурированное логирование
│   └── validator/                # Валидация запросов
├── Dockerfile                    # Multi-stage Docker build
├── docker-compose.yml            # Конфигурация Docker Compose
├── Makefile                      # Команды для сборки и разработки
└── README.md
```

## Требования

- Docker 20.10+
- Docker Compose 1.29+
- Go 1.22+ (для локальной разработки)
- Make (опционально)

## Быстрый старт

### Использование Docker Compose (рекомендуется)

1. Клонируйте репозиторий:
```bash
git clone git@github.com:dimozavrrrik/TaskManager.git
cd TaskManager
```

2. Скопируйте файл с переменными окружения:
```bash
cp .env.example .env
```

3. Запустите сервисы:
```bash
docker-compose up -d
```

4. Проверьте работоспособность сервиса:
```bash
curl http://localhost:8080/api/v1/health
```

Ожидаемый ответ:
```json
{
  "status": "ok",
  "service": "taskmanager"
}
```

### Локальная разработка

1. Установите зависимости:
```bash
go mod download
```

2. Запустите PostgreSQL:
```bash
docker-compose up postgres -d
```

3. Установите переменные окружения:
```bash
export DATABASE_URL="postgres://taskmanager:taskmanager_password@localhost:5432/taskmanager?sslmode=disable"
export SERVER_ADDRESS=":8080"
export LOG_LEVEL="debug"
```

4. Примените миграции:
```bash
make migrate-up
```

5. Запустите приложение:
```bash
go run cmd/api/main.go
```

## Конфигурация

Переменные окружения:

| Переменная         | Описание                          | По умолчанию |
|-------------------|-----------------------------------|--------------|
| SERVER_ADDRESS    | Адрес для запуска сервера         | :8080        |
| DATABASE_URL      | Строка подключения к PostgreSQL   | -            |
| LOG_LEVEL         | Уровень логирования               | info         |
| ENVIRONMENT       | Окружение (development/production)| development  |
| DB_MAX_OPEN_CONNS | Макс. количество подключений к БД | 25           |
| DB_MAX_IDLE_CONNS | Макс. количество idle подключений| 5            |
| DB_MAX_IDLE_TIME  | Макс. время idle подключения      | 15m          |
| JWT_SECRET        | Секретный ключ для подписи JWT (мин. 32 символа) | - |
| JWT_ACCESS_EXPIRY_MIN | Время жизни access токена (минуты) | 15 |
| JWT_REFRESH_EXPIRY_DAYS | Время жизни refresh токена (дни) | 7 |

**ВАЖНО**: В production обязательно установите надежный `JWT_SECRET` (минимум 32 случайных символа)!

## Документация API

### Базовый URL

```
http://localhost:8080/api/v1
```

### Аутентификация

API использует JWT (JSON Web Tokens) для аутентификации:

- **Access Token**: Срок действия 15 минут, передается в заголовке `Authorization: Bearer <token>`
- **Refresh Token**: Срок действия 7 дней, хранится в базе данных для возможности отзыва
- **Пароли**: Хешируются с помощью bcrypt перед сохранением

#### Регистрация и вход

**Регистрация нового пользователя:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Иван Иванов",
    "department": "Разработка",
    "position": "Разработчик",
    "email": "ivan@example.com",
    "password": "SecurePass123!"
  }'
```

**Вход в систему:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "ivan@example.com",
    "password": "SecurePass123!"
  }'
```

Ответ содержит access_token, refresh_token и информацию о пользователе.

**Обновление токена:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "your-refresh-token"
  }'
```

**Выход из системы:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "your-refresh-token"
  }'
```

### Endpoints

#### Health Check

```http
GET /health
```

Ответ:
```json
{
  "status": "ok",
  "service": "taskmanager"
}
```

#### Задачи

**Создание задачи**
```http
POST /tasks
Authorization: Bearer <access-token>

{
  "title": "Реализовать функцию X",
  "description": "Добавить новую функцию",
  "status": "pending",
  "priority": "medium"
}
```

**Список всех задач**
```http
GET /tasks?page=1&page_size=20&status=in_progress
```

**Получение задачи по ID**
```http
GET /tasks/{id}
```

**Обновление статуса задачи**
```http
PATCH /tasks/{id}/status

{
  "status": "in_progress"
}
```

Автоматически создается системное сообщение:
```
"Task status changed from 'new' to 'in_progress'"
```

**Архивация задачи**
```http
PATCH /tasks/{id}/archive
```

**Получение участников задачи**
```http
GET /tasks/{id}/participants
```

**Добавление участника**
```http
POST /tasks/{id}/participants

{
  "employee_id": "uuid",
  "role": "responsible"
}
```

**Получение задач сотрудника**
```http
GET /employees/{id}/tasks?page=1&page_size=20
```

Возвращает все задачи, где сотрудник является участником (главный экран).

### Формат ответов

**Успешный ответ**:
```json
{
  "success": true,
  "data": { ... }
}
```

**Ответ с ошибкой**:
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      {
        "field": "title",
        "message": "This field is required"
      }
    ]
  }
}
```

**Ответ с пагинацией**:
```json
{
  "success": true,
  "data": {
    "data": [ ... ],
    "total": 100,
    "page": 1,
    "page_size": 20,
    "total_pages": 5
  }
}
```

### Коды ошибок

- `VALIDATION_ERROR` (400): Некорректные данные в запросе
- `BAD_REQUEST` (400): Неправильный формат запроса
- `NOT_FOUND` (404): Ресурс не найден
- `CONFLICT` (409): Конфликт ресурсов
- `UNAUTHORIZED` (401): Требуется аутентификация
- `INTERNAL_ERROR` (500): Ошибка сервера

## Схема базы данных

### Основные таблицы

1. **employees** - Профили сотрудников
   - id, name, department, position, email
   - Мягкое удаление (deleted_at)

2. **tasks** - Управление задачами
   - id, title, description, status, priority, created_by, archived
   - Статусы: new, in_progress, code_review, testing, returned_with_errors, closed

3. **task_participants** - Назначение ролей
   - id, task_id, employee_id, role
   - Роли: executor, responsible, customer
   - Уникальное ограничение: (task_id, employee_id, role)

4. **task_messages** - Коммуникация в задачах
   - id, task_id, author_id, content, is_system_message
   - NULL author_id для системных сообщений

5. **time_entries** - Учет времени
   - id, task_id, employee_id, hours, description, entry_date

### Представления (Views)

- **task_time_summary**: Суммирование времени по задачам
- **employee_tasks_view**: Данные для главного экрана (задачи с ролями)

### Ключевые особенности

- UUID первичные ключи
- Мягкое удаление для восстановления данных
- Автоматическое обновление timestamps (триггеры)
- Стратегические индексы на внешних ключах
- Каскадное удаление для дочерних записей

## Разработка

### Запуск тестов

```bash
# Запустить все тесты
make test

# Запустить с покрытием
make test-coverage
```

### Качество кода

```bash
# Форматирование кода
make fmt

# Запустить линтер
make lint
```

### Миграции базы данных

```bash
# Применить миграции
make migrate-up

# Откатить миграции
make migrate-down
```

## Команды Makefile

```bash
make build          # Собрать бинарный файл
make run            # Запустить приложение локально
make test           # Запустить тесты
make docker-build   # Собрать Docker образ
make docker-up      # Запустить сервисы Docker Compose
make docker-down    # Остановить сервисы Docker Compose
make docker-logs    # Просмотр логов Docker Compose
make migrate-up     # Применить миграции БД
make migrate-down   # Откатить миграции БД
make lint           # Запустить линтер
make fmt            # Форматировать код
make deps           # Загрузить зависимости
make clean          # Очистить артефакты сборки
```

## Архитектурные решения

### 1. Stateless дизайн
- Нет хранения состояния на стороне сервера
- Все состояние сохраняется в PostgreSQL
- Возможность горизонтального масштабирования

### 2. Repository Pattern
- Абстракция слоя доступа к данным
- Легкое тестирование с помощью моков
- Возможность замены базы данных

### 3. Dependency Injection
- Инъекция через конструкторы
- Зависимости на основе интерфейсов
- Способствует тестируемости и слабой связанности

### 4. Управление транзакциями
- Слой Service обрабатывает транзакции
- Атомарные сложные операции (например, создание задачи с участниками + системное сообщение)

### 5. Системные сообщения
- Автоматическое создание при смене статуса
- NULL author_id для системных сообщений
- Полный аудит-трейл

### 6. Мягкое удаление
- Колонка deleted_at
- Возможность восстановления данных
- Фильтрация в запросах

### 7. Структурированное логирование
- JSON формат для машинной обработки
- Отслеживание Request ID
- Уровни логов: debug, info, warn, error

### 8. Обработка ошибок
- Пользовательский тип AppError
- Маппинг на HTTP статус-коды
- Детальные ошибки валидации

### 9. Цепочка Middleware
- Recovery → Logging → CORS
- Генерация Request ID
- Восстановление после паник со stack traces

## Мониторинг и наблюдаемость

- **Health Check**: Endpoint `/api/v1/health`
- **Структурированные JSON логи**: Машинно-читаемый вывод логов
- **Отслеживание Request ID**: Трассировка запросов через систему
- **Логирование времени ответа**: Мониторинг производительности

## Соображения безопасности

- Предотвращение SQL-инъекций (параметризованные запросы)
- Валидация входных данных на всех endpoints
- Мягкое удаление для восстановления данных
- CORS middleware для интеграции с фронтендом
- Восстановление после паник для предотвращения DoS

## Соображения производительности

- Пулинг подключений к БД (настраиваемый)
- Индексы на внешних ключах и часто запрашиваемых колонках
- Представления БД для сложных запросов
- Пагинация на endpoints со списками

## Развертывание

### Продакшн развертывание

1. Соберите продакшн образ:
```bash
docker build -t taskmanager:latest .
```

2. Запустите с продакшн настройками:
```bash
docker run -e ENVIRONMENT=production \
           -e DATABASE_URL=<prod-db-url> \
           -e LOG_LEVEL=info \
           -p 8080:8080 \
           taskmanager:latest
```

### Соображения масштабирования

- **Горизонтальное масштабирование**: Stateless дизайн поддерживает несколько реплик
- **Пул подключений к БД**: Размер зависит от количества реплик
- **Read реплики**: Рассмотрите для тяжелых read workloads
- **Кеширование**: Добавьте Redis для часто запрашиваемых данных
- **Load Balancer**: Используйте nginx или облачный балансировщик

## Пример использования

1. **Регистрация нового пользователя**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Иван Иванов",
    "department": "Разработка",
    "position": "Разработчик",
    "email": "ivan@example.com",
    "password": "SecurePass123!"
  }'
```

Ответ содержит `access_token` и `refresh_token`.

2. **Вход в систему**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "ivan@example.com",
    "password": "SecurePass123!"
  }'
```

3. **Создание задачи** (требуется аутентификация):
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access-token>" \
  -d '{
    "title": "Реализовать функцию X",
    "description": "Добавить новую функцию",
    "status": "pending",
    "priority": "medium"
  }'
```

4. **Получение списка задач**:
```bash
curl http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer <access-token>"
```

5. **Обновление статуса задачи**:
```bash
curl -X PATCH http://localhost:8080/api/v1/tasks/<task-id>/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access-token>" \
  -d '{"status": "in_progress"}'
```

6. **Добавление участника к задаче**:
```bash
curl -X POST http://localhost:8080/api/v1/tasks/<task-id>/participants \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access-token>" \
  -d '{
    "employee_id": "<uuid>",
    "role": "executor"
  }'
```

## Устранение неполадок

### Проблемы с подключением к БД

```bash
# Проверить, что PostgreSQL запущен
docker-compose ps postgres

# Проверить логи PostgreSQL
docker-compose logs postgres

# Тестовое подключение
docker exec -it taskmanager_postgres psql -U taskmanager -d taskmanager -c "SELECT 1;"
```

### API не отвечает

```bash
# Проверить логи API
docker-compose logs api

# Проверить health endpoint
curl http://localhost:8080/api/v1/health

# Перезапустить сервисы
docker-compose restart
```

### Ошибки миграций

```bash
# Проверить синтаксис файла миграции
cat internal/database/migrations/001_init_schema.up.sql

# Применить миграцию вручную
make migrate-up

# Откатить при необходимости
make migrate-down
```

## Frontend приложение

### Доступ к приложению

После запуска через Docker Compose:
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080/api/v1

### Функционал

- **Аутентификация**: Регистрация и вход с JWT токенами
- **Dashboard**: Статистика по задачам
- **Управление задачами**:
  - Просмотр списка задач с фильтрацией
  - Создание новых задач
  - Просмотр деталей задачи
  - Изменение статуса задачи
  - Архивирование задач
- **Управление участниками**: Добавление участников к задачам
- **Просмотр сотрудников**: Список всех зарегистрированных сотрудников
- **Responsive дизайн**: Адаптивная верстка для всех устройств
- **PWA**: Возможность установки как desktop приложение

### Разработка Frontend

Для локальной разработки frontend:

```bash
cd TaskManager.Client
dotnet run
```

Frontend будет доступен на `http://localhost:5000`

Для изменения API endpoint отредактируйте `wwwroot/appsettings.json`:

```json
{
  "ApiBaseUrl": "http://localhost:8080/api/v1"
}
```

### Production сборка Frontend

```bash
cd TaskManager.Client
dotnet publish -c Release -o ./publish
```

Или через Docker:

```bash
cd TaskManager.Client
docker build -t taskmanager-frontend .
docker run -p 3000:80 taskmanager-frontend
```

## Полный стек запуска

Для запуска всего приложения (Backend + Frontend + Database):

```bash
# 1. Скопируйте .env файл
cp .env.example .env

# 2. Отредактируйте JWT_SECRET в .env файле (обязательно для production!)

# 3. Запустите все сервисы
docker-compose up -d

# 4. Проверьте статус
docker-compose ps

# 5. Откройте браузер
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080/api/v1/health
```

### Доступ к сервисам

| Сервис | URL | Описание |
|--------|-----|----------|
| Frontend | http://localhost:3000 | Blazor WebAssembly приложение |
| Backend API | http://localhost:8080 | REST API |
| PostgreSQL | localhost:5432 | База данных |
| PgAdmin | http://localhost:5050 | Управление БД (профиль dev) |

### Остановка сервисов

```bash
docker-compose down
```

Для удаления volumes (БД будет очищена):

```bash
docker-compose down -v
```

---

**Создано с ❤️ на Go и PostgreSQL**
