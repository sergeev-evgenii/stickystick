/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: 'export', // Статический экспорт
  images: {
    unoptimized: true, // Для статического экспорта
  },
  trailingSlash: true, // Для правильной работы с nginx
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL || 'https://stickystick.ru',
    NEXT_PUBLIC_TELEGRAM_URL: process.env.NEXT_PUBLIC_TELEGRAM_URL || '',
    NEXT_PUBLIC_VK_URL: process.env.NEXT_PUBLIC_VK_URL || '',
    NEXT_PUBLIC_MAX_URL: process.env.NEXT_PUBLIC_MAX_URL || '',
  },
  // Разрешить запросы к dev-серверу с продакшен-домена (при отладке)
  allowedDevOrigins: ['https://stickystick.ru', 'http://stickystick.ru'],
}

module.exports = nextConfig
