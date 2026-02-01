/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: 'export', // Статический экспорт
  images: {
    unoptimized: true, // Для статического экспорта
  },
  trailingSlash: true, // Для правильной работы с nginx
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:5000',
  },
}

module.exports = nextConfig
