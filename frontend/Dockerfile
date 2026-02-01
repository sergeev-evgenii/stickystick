# Multi-stage build для оптимизации
FROM node:20-alpine AS builder

WORKDIR /app

# Копируем package файлы и устанавливаем зависимости
COPY package.json package-lock.json* ./
RUN npm ci

# Копируем исходный код и собираем приложение
COPY . .
RUN npm run build

# Production stage
FROM nginx:alpine

# Копируем собранные файлы из builder stage (Next.js создает директорию 'out' при статическом экспорте)
COPY --from=builder /app/out /usr/share/nginx/html

# Копируем конфигурацию nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 8090

CMD ["nginx", "-g", "daemon off;"]
