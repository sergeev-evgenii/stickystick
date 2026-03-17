'use client'

import { useState, useCallback } from 'react'
import Link from 'next/link'
import { useAuthStore } from '@/store/authStore'
import VideoFeed from '@/components/VideoFeed'
import { analyticsApi } from '@/lib/api/analytics'

function randomBetween(min: number, max: number): number {
  return Math.floor(Math.random() * (max - min + 1)) + min
}

export default function Home() {
  const { user, isAuthenticated, logout } = useAuthStore()
  const [serviceUnavailable, setServiceUnavailable] = useState(false)
  const [waiting, setWaiting] = useState(false)

  const handleGenerateVideo = useCallback(async () => {
    if (waiting || serviceUnavailable) return
    try {
      await analyticsApi.logGenerateVideoClick()
    } catch (_) {
      // логируем клик в любом случае; ошибка сети — не критична
    }
    setWaiting(true)
    const delayMs = randomBetween(5, 10) * 1000
    setTimeout(() => {
      setWaiting(false)
      setServiceUnavailable(true)
      setTimeout(() => setServiceUnavailable(false), 3000)
    }, delayMs)
  }, [waiting, serviceUnavailable])
  const telegramUrl = process.env.NEXT_PUBLIC_TELEGRAM_URL
  const vkUrl = process.env.NEXT_PUBLIC_VK_URL
  const maxUrl = process.env.NEXT_PUBLIC_MAX_URL

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Header */}
      <header className="bg-white shadow-sm sticky top-0 z-10">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <Link href="/" className="flex flex-col font-verveine tracking-wider">
              <span className="text-4xl font-bold text-amber-500 uppercase">Sticky Stick</span>
              <span className="text-lg font-semibold text-amber-500 uppercase">липкая палка</span>
              <span className="text-base font-semibold text-black">просто мемы. просто смешно. иногда эпично</span>
            </Link>
            <nav className="flex items-center gap-4">
              {telegramUrl && (
                <a
                  href={telegramUrl}
                  target="_blank"
                  rel="noreferrer"
                  className="text-gray-600 hover:text-sky-600"
                  aria-label="Telegram"
                  title="Telegram"
                >
                  <svg className="w-6 h-6" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
                    <path d="M9.993 15.674 9.6 20.2c.563 0 .807-.243 1.1-.533l2.64-2.52 5.47 4.002c1.003.553 1.717.262 1.982-.93l3.59-16.82c.338-1.56-.563-2.17-1.55-1.804L1.726 9.17c-1.52.59-1.497 1.437-.259 1.82l5.65 1.763L20.24 4.73c.62-.41 1.185-.183.72.227L9.993 15.674z" />
                  </svg>
                </a>
              )}
              {vkUrl && (
                <a
                  href={vkUrl}
                  target="_blank"
                  rel="noreferrer"
                  className="text-gray-600 hover:text-blue-600"
                  aria-label="VK"
                  title="VK"
                >
                  <svg className="w-6 h-6" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
                    <path d="M12.784 16.58c.4 0 .57-.28.57-.66 0-.4-.01-.93.01-1.35.03-.78.25-1.32.86-1.32.59 0 1.14.63 1.62 1.22.53.65 1.05 1.32 1.62 1.9.42.42.76.21 1.03-.06l1.08-1.1c.26-.28.31-.63.07-.94-.31-.41-.77-.94-1.24-1.5-.55-.65-1.1-1.33-1.32-1.76-.22-.41-.12-.7.14-1.05.53-.7 1.2-1.55 1.72-2.38.26-.42.44-.78.6-1.19.17-.45-.03-.65-.48-.65h-1.58c-.41 0-.6.15-.72.46-.22.6-.56 1.3-1.08 2.12-.55.86-1.1 1.6-1.52 1.6-.29 0-.41-.32-.43-.82-.03-.71.02-1.7.02-2.32 0-.79-.24-1.04-.69-1.14-.22-.05-.53-.09-1.13-.09-.77 0-1.4.01-1.76.2-.25.13-.41.4-.3.41.37.04.5.23.58.51.12.41.12 1.33.12 1.98 0 .46.03 1.08-.1 1.33-.1.19-.2.26-.36.26-.38 0-.95-.72-1.51-1.6-.52-.82-.92-1.72-1.11-2.23-.12-.3-.27-.44-.68-.44H3.63c-.45 0-.6.23-.48.61.18.56.5 1.25.8 1.82.7 1.33 1.63 2.67 2.74 3.86 1.27 1.35 2.78 2.2 3.86 2.2z" />
                  </svg>
                </a>
              )}
              {maxUrl && (
                <a
                  href={maxUrl}
                  target="_blank"
                  rel="noreferrer"
                  className="text-gray-600 hover:text-violet-600"
                  aria-label="Max"
                  title="Max"
                >
                  <svg className="w-6 h-6" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
                    <path d="M4 6h3.2l2.1 6.2L11.4 6H14l2.1 6.2L18.2 6H21l-3.2 12h-2.7l-2.4-6.7L10.3 18H7.6L4 6z" />
                  </svg>
                </a>
              )}
              {isAuthenticated ? (
                <>
                  <Link
                    href="/upload"
                    className="bg-blue-600 text-white px-4 py-2 rounded-lg font-semibold hover:bg-blue-700 transition-colors"
                  >
                    Загрузить
                  </Link>
                  {user?.is_admin && (
                    <Link
                      href="/moderation"
                      className="bg-orange-600 text-white px-4 py-2 rounded-lg font-semibold hover:bg-orange-700 transition-colors"
                    >
                      Модерация
                    </Link>
                  )}
                  <span className="text-gray-600">@{user?.username}</span>
                  <button
                    onClick={logout}
                    className="text-gray-600 hover:text-gray-800 font-semibold"
                  >
                    Выйти
                  </button>
                </>
              ) : (
                <>
                  <Link
                    href="/login"
                    className="text-gray-600 hover:text-gray-800 font-semibold"
                  >
                    Войти
                  </Link>
                  <Link
                    href="/register"
                    className="bg-green-600 text-white px-4 py-2 rounded-lg font-semibold hover:bg-green-700 transition-colors"
                  >
                    Регистрация
                  </Link>
                </>
              )}
            </nav>
          </div>
          {/* Кнопка «Сгенерировать своё видео» — справа, под навигацией */}
          <div className="flex flex-col items-end gap-1 mt-3 pt-3 border-t border-gray-100">
            <button
              type="button"
              onClick={handleGenerateVideo}
              disabled={waiting || serviceUnavailable}
              className="bg-amber-500 text-white px-6 py-2 rounded-lg font-semibold hover:bg-amber-600 disabled:opacity-60 disabled:cursor-not-allowed transition-colors"
            >
              Сгенерировать своё видео
            </button>
            {waiting && (
              <p className="text-gray-600 font-medium text-sm">Ожидание…</p>
            )}
            {serviceUnavailable && !waiting && (
              <p className="text-amber-700 font-medium text-sm">Service temporary unavailable</p>
            )}
          </div>
        </div>
      </header>

      {/* Video Feed */}
      <VideoFeed />
    </div>
  )
}
