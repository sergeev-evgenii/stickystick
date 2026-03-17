# Ручные миграции БД

Миграции применяются **только вручную**. Не используйте `RUN_MIGRATIONS=1` при старте приложения.

## Порядок выполнения

Выполняйте файлы по порядку (один раз; повторный запуск безопасен — скрипты идемпотентны):

| Файл | Назначение |
|------|------------|
| `001_rename_video_url.sql` | Переименование `video_url` → `media_url`, добавление `media_type` в `videos` |
| `002_video_reports.sql` | Создание таблицы `video_reports` (жалобы на видео) |
| `003_schema_extras.sql` | Колонки `moderation_status` в `videos`, `is_admin` в `users` |

## Как выполнить

### Вариант 1: через psql

Подставьте свои хост, порт, пользователя и базу (например, для основной БД `stickystick`):

```bash
cd backend/migrations

export PGHOST=localhost
export PGPORT=5433
export PGUSER=stickystick
export PGPASSWORD=stickystick
export PGDATABASE=stickystick

psql -f 001_rename_video_url.sql
psql -f 002_video_reports.sql
psql -f 003_schema_extras.sql
```

Или одной строкой подключения:

```bash
psql "postgres://stickystick:stickystick@localhost:5433/stickystick?sslmode=disable" -f 001_rename_video_url.sql
psql "postgres://stickystick:stickystick@localhost:5433/stickystick?sslmode=disable" -f 002_video_reports.sql
psql "postgres://stickystick:stickystick@localhost:5433/stickystick?sslmode=disable" -f 003_schema_extras.sql
```

### Вариант 2: скрипт run_migrations.sh

Задайте переменную `DATABASE_URL` (URL подключения к PostgreSQL) и запустите:

```bash
cd backend/migrations
export DATABASE_URL="postgres://stickystick:stickystick@localhost:5433/stickystick?sslmode=disable"
./run_migrations.sh
```

Скрипт выполнит все `*.sql` в порядке по имени.

## На сервере (production)

Скопируйте нужные файлы на сервер и выполните их в той же базе, к которой подключается бэкенд (например, внутри контейнера или с хоста с доступом к БД):

```bash
# Пример: выполнение внутри контейнера postgres
docker exec -i stickystick_postgres psql -U stickystick -d stickystick < 002_video_reports.sql
```

После применения миграций перезапускать приложение не обязательно (таблицы/колонки подхватятся при следующем запросе).
