#!/usr/bin/env bash
# Запуск всех SQL-миграций в порядке по имени.
# Требуется: DATABASE_URL или переменные PGHOST, PGPORT, PGUSER, PGPASSWORD, PGDATABASE.

set -e
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

if [ -n "$DATABASE_URL" ]; then
  RUN="psql \"$DATABASE_URL\""
else
  RUN="psql"
fi

for f in 001_rename_video_url.sql 002_video_reports.sql 003_schema_extras.sql; do
  if [ -f "$f" ]; then
    echo "Running $f ..."
    eval "$RUN -f $f"
  fi
done

echo "Done."
