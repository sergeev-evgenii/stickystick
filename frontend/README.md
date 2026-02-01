# Sticky Stick Frontend

Frontend для сервиса коротких видео на Next.js.

## Структура проекта

```
frontend/
├── src/
│   ├── app/                 # Next.js App Router
│   │   ├── layout.tsx
│   │   ├── page.tsx
│   │   └── globals.css
│   ├── components/          # React компоненты
│   │   └── VideoPlayer.tsx
│   ├── lib/                 # Утилиты и API клиенты
│   │   ├── api.ts
│   │   └── api/
│   │       ├── auth.ts
│   │       └── video.ts
│   └── store/               # State management (Zustand)
│       └── authStore.ts
├── public/                  # Статические файлы
├── .env.example
├── package.json
├── tsconfig.json
├── next.config.js
└── tailwind.config.js
```

## Установка

1. Установите зависимости:
```bash
npm install
# или
yarn install
```

2. Создайте файл `.env.local` на основе `.env.example`:
```bash
cp .env.example .env.local
```

3. Настройте переменные окружения в `.env.local`

4. Запустите dev сервер:
```bash
npm run dev
# или
yarn dev
```

Приложение будет доступно по адресу [http://localhost:3000](http://localhost:3000)

## Технологии

- Next.js 14 (App Router)
- React 18
- TypeScript
- Tailwind CSS
- Zustand (state management)
- Axios (HTTP client)

## Основные функции

- Авторизация и регистрация
- Просмотр ленты видео
- Загрузка видео
- Лайки и комментарии
- Профили пользователей
