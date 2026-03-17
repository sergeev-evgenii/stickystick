'use client'

import './globals.css'
import SettingsLoader from '@/components/SettingsLoader'

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
        <link rel="icon" href="/favicon.jpg" type="image/jpeg" />
      </head>
      <body className="font-sans antialiased">
        <SettingsLoader />
        {children}
      </body>
    </html>
  )
}
