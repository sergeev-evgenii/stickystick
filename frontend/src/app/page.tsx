'use client'

import Link from 'next/link'
import { useAuthStore } from '@/store/authStore'
import VideoFeed from '@/components/VideoFeed'

export default function Home() {
  const { user, isAuthenticated, logout } = useAuthStore()

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Header */}
      <header className="bg-white shadow-sm sticky top-0 z-10">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <Link href="/" className="text-2xl font-bold text-blue-600">
              Sticky Stick
            </Link>
            <nav className="flex items-center gap-4">
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
        </div>
      </header>

      {/* Video Feed */}
      <VideoFeed />
    </div>
  )
}
