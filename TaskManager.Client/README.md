# TaskManager Frontend

Blazor WebAssembly приложение для управления задачами.

## Технологии

- Blazor WebAssembly (.NET 8)
- MudBlazor UI компоненты
- Blazored.LocalStorage для хранения JWT токенов

## Разработка

### Локальный запуск

```bash
dotnet run
```

Приложение будет доступно на `http://localhost:5000`

### Настройка API URL

Измените `wwwroot/appsettings.json`:

```json
{
  "ApiBaseUrl": "http://localhost:8080/api/v1"
}
```

## Production

### Docker

```bash
docker build -t taskmanager-frontend .
docker run -p 3000:80 taskmanager-frontend
```

### Docker Compose

Из корневой директории проекта:

```bash
docker-compose up -d
```

Frontend будет доступен на `http://localhost:3000`

## Функционал

- Аутентификация (JWT)
- Управление задачами
- Просмотр сотрудников
- Dashboard с статистикой
- Responsive дизайн

## Структура

```
TaskManager.Client/
├── Models/          # DTOs для API
├── Services/        # HTTP сервисы
├── Pages/           # Blazor страницы
├── Shared/          # Общие компоненты
└── wwwroot/         # Статические файлы
```
