# Полная инструкция по сборке и запуску TaskManager

## Шаг 1: Остановка и очистка (если запущено)

```bash
# Остановить все контейнеры
docker-compose down

# Удалить все контейнеры и образы (если нужна полная пересборка)
docker-compose down -v
docker rmi taskmanager-api taskmanager-frontend 2>/dev/null || true
```

## Шаг 2: Проверка конфигурации

### 2.1 Проверьте .env файл

Файл `.env` должен содержать:

```env
JWT_SECRET=your-super-secret-jwt-key-minimum-32-characters-long-please-change-this
DATABASE_URL=postgres://taskmanager:taskmanager_password@localhost:5432/taskmanager?sslmode=disable
SERVER_ADDRESS=:8080
ENVIRONMENT=production
LOG_LEVEL=info
```

**ВАЖНО**: Замените `JWT_SECRET` на свой уникальный ключ (минимум 32 символа).

### 2.2 Проверьте appsettings.json

Файл `TaskManager.Client/wwwroot/appsettings.json` должен содержать:

```json
{
  "ApiBaseUrl": "/api/v1"
}
```

**ВАЖНО**: Должен быть относительный путь `/api/v1`, а НЕ `http://localhost:8080/api/v1`!

## Шаг 3: Сборка всех образов

```bash
# Сборка БЕЗ кэша (рекомендуется для первого раза)
docker-compose build --no-cache

# ИЛИ сборка с кэшем (быстрее, но может использовать старые файлы)
docker-compose build
```

Процесс займет 3-5 минут. Вы увидите:
- ✅ Building postgres (используется готовый образ)
- ✅ Building api (компиляция Go)
- ✅ Building frontend (компиляция Blazor WASM)

## Шаг 4: Запуск всех сервисов

```bash
# Запуск в фоновом режиме
docker-compose up -d

# ИЛИ запуск с выводом логов (для отладки)
docker-compose up
```

## Шаг 5: Проверка статуса

```bash
# Проверить статус всех контейнеров
docker-compose ps
```

Вы должны увидеть:
```
NAME                   STATUS                   PORTS
taskmanager_frontend   Up X seconds             0.0.0.0:8081->80/tcp
taskmanager_api        Up X seconds             0.0.0.0:8080->8080/tcp
taskmanager_postgres   Up X seconds (healthy)   0.0.0.0:5432->5432/tcp
```

## Шаг 6: Проверка работоспособности

### 6.1 Проверьте Backend API

```bash
curl http://localhost:8080/api/v1/health
```

Ожидаемый ответ:
```json
{"status":"ok","service":"taskmanager"}
```

### 6.2 Проверьте Frontend

```bash
curl -I http://localhost:8081
```

Ожидаемый ответ: `HTTP/1.1 200 OK`

### 6.3 Проверьте PostgreSQL

```bash
docker exec -it taskmanager_postgres psql -U taskmanager -d taskmanager -c "SELECT 1;"
```

## Шаг 7: Откройте приложение

Откройте браузер и перейдите на:

```
http://localhost:8081
```

## Шаг 8: Первая регистрация

1. Нажмите **"Зарегистрироваться"**
2. Заполните форму:
   - **Имя**: Ваше имя
   - **Email**: test@example.com
   - **Отдел**: Разработка
   - **Должность**: Разработчик
   - **Пароль**: SecurePass123!
   - **Подтвердите пароль**: SecurePass123!
3. Нажмите **"Зарегистрироваться"**
4. Вы автоматически войдете в систему!

## Просмотр логов

### Все сервисы сразу
```bash
docker-compose logs -f
```

### Только Frontend
```bash
docker-compose logs -f frontend
```

### Только API
```bash
docker-compose logs -f api
```

### Только PostgreSQL
```bash
docker-compose logs -f postgres
```

## Управление сервисами

### Перезапуск всех сервисов
```bash
docker-compose restart
```

### Перезапуск одного сервиса
```bash
docker-compose restart frontend
```

### Остановка всех сервисов
```bash
docker-compose down
```

### Остановка с удалением данных БД
```bash
docker-compose down -v
```

### Пересборка конкретного сервиса
```bash
# Только frontend
docker-compose build --no-cache frontend
docker-compose up -d frontend

# Только API
docker-compose build --no-cache api
docker-compose up -d api
```

## Устранение проблем

### Проблема: "Failed to fetch" при регистрации

**Причина**: Frontend не может подключиться к API.

**Решение 1**: Проверьте appsettings.json внутри контейнера
```bash
docker exec taskmanager_frontend sh -c "cat /usr/share/nginx/html/appsettings.json"
```

Должно быть: `{"ApiBaseUrl": "/api/v1"}`

Если НЕТ - пересоберите frontend:
```bash
# 1. Убедитесь что TaskManager.Client/wwwroot/appsettings.json содержит "/api/v1"
cat TaskManager.Client/wwwroot/appsettings.json

# 2. Пересоберите
docker-compose build --no-cache frontend
docker-compose up -d frontend
```

**Решение 2**: Проверьте что API работает
```bash
curl http://localhost:8080/api/v1/health
```

**Решение 3**: Проверьте логи nginx
```bash
docker-compose logs frontend | grep api
```

### Проблема: Порт 8080 или 8081 занят

```bash
# Найти процесс
netstat -ano | findstr :8080
netstat -ano | findstr :8081

# Убить процесс (замените PID на ваш)
taskkill /PID <PID> /F

# Или измените порт в docker-compose.yml
```

### Проблема: База данных не подключается

```bash
# Проверьте логи PostgreSQL
docker-compose logs postgres

# Проверьте подключение
docker exec -it taskmanager_postgres psql -U taskmanager -c "\l"

# Пересоздайте БД
docker-compose down -v
docker-compose up -d
```

### Проблема: Ошибка компиляции Go

```bash
# Убедитесь что go.mod и go.sum актуальны
cd cmd/api
go mod tidy
cd ../..

# Пересоберите
docker-compose build --no-cache api
```

### Проблема: Ошибка компиляции Blazor

```bash
# Проверьте что все пакеты установлены
cd TaskManager.Client
cat TaskManager.Client.csproj

# Убедитесь что есть:
# Microsoft.AspNetCore.Components.Authorization

# Пересоберите
cd ..
docker-compose build --no-cache frontend
```

## Полная пересборка с нуля

Если что-то пошло не так и нужно начать заново:

```bash
# 1. Остановить и удалить ВСЕ
docker-compose down -v
docker rmi taskmanager-api taskmanager-frontend 2>/dev/null || true
docker system prune -f

# 2. Проверить конфигурацию
cat TaskManager.Client/wwwroot/appsettings.json
# Должно быть: {"ApiBaseUrl": "/api/v1"}

cat .env
# Должен содержать JWT_SECRET

# 3. Собрать заново без кэша
docker-compose build --no-cache

# 4. Запустить
docker-compose up -d

# 5. Проверить статус
docker-compose ps

# 6. Проверить логи
docker-compose logs -f
```

## Тестирование через curl

После успешного запуска можно протестировать API напрямую:

```bash
# 1. Регистрация
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "department": "IT",
    "position": "Developer",
    "email": "test@example.com",
    "password": "SecurePass123!"
  }'

# Сохраните access_token из ответа

# 2. Создание задачи
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <ваш-access-token>" \
  -d '{
    "title": "Тестовая задача",
    "description": "Описание задачи",
    "status": "pending",
    "priority": "medium"
  }'
```

## Порты

| Сервис | Внутренний порт | Внешний порт | URL |
|--------|----------------|--------------|-----|
| Frontend | 80 | 8081 | http://localhost:8081 |
| Backend API | 8080 | 8080 | http://localhost:8080 |
| PostgreSQL | 5432 | 5432 | localhost:5432 |

## Production deployment

Для production:

1. **Смените JWT_SECRET** на случайную строку (32+ символа)
2. **Настройте HTTPS** через reverse proxy
3. **Измените пароль БД** в docker-compose.yml
4. **Настройте резервное копирование**
5. **Настройте мониторинг**

---

**Готово! Приложение должно работать на http://localhost:8081** 🎉
