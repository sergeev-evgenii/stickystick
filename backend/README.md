# Sticky Stick Backend

Backend сервиса для коротких видео на Go.

## Структура проекта

```
backend/
├── cmd/
│   └── api/
│       └── main.go          # Точка входа приложения
├── internal/
│   ├── config/              # Конфигурация
│   ├── handler/             # HTTP handlers
│   ├── models/              # Модели данных
│   ├── repository/          # Слой работы с БД
│   └── service/             # Бизнес-логика
├── migrations/              # Миграции БД (опционально)
├── .env.example
├── go.mod
└── README.md
```

## Установка

1. Установите зависимости:
```bash
go mod download
```

2. Создайте файл `.env` на основе `.env.example`:
```bash
cp .env.example .env
```

3. Настройте переменные окружения в `.env`

4. Запустите приложение:
```bash
go run cmd/api/main.go
```

## API Endpoints

### Auth
- `POST /api/v1/auth/register` - Регистрация
- `POST /api/v1/auth/login` - Вход

### Users
- `GET /api/v1/users/:id` - Получить профиль
- `PUT /api/v1/users/:id` - Обновить профиль

### Videos
- `GET /api/v1/videos` - Лента видео
- `GET /api/v1/videos/:id` - Получить видео
- `POST /api/v1/videos` - Загрузить видео (JSON, старый метод)
- `POST /api/v1/videos/upload` - Загрузить медиа файл (multipart/form-data)
  - Поддерживает: видео (mp4, mov, avi, webm, mkv), фото (jpg, jpeg, png, webp), GIF
  - Параметры: `file` (обязательно), `title` (обязательно), `description` (опционально), `thumbnail` (опционально, для видео)
- `DELETE /api/v1/videos/:id` - Удалить видео
- `POST /api/v1/videos/:id/like` - Лайкнуть видео
- `DELETE /api/v1/videos/:id/like` - Убрать лайк
- `POST /api/v1/videos/:id/comment` - Добавить комментарий

### Статические файлы
- `GET /uploads/*` - Доступ к загруженным медиа файлам

## Технологии

- Go 1.21+
- Gin (HTTP framework)
- GORM (ORM)
- PostgreSQL
- JWT (аутентификация)
