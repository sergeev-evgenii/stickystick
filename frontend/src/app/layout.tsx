'use client'

import { Inter } from 'next/font/google'
import './globals.css'

const inter = Inter({ subsets: ['latin'] })

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="ru">
      <head>
        <title>Sticky Stick</title>
        <meta name="description" content="Сервис для коротких видео" />
      </head>
      <body className={inter.className}>{children}</body>
    </html>
  )
}
