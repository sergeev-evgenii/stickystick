'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { videoApi, Video } from '@/lib/api/video'
import { useAuthStore } from '@/store/authStore'
import { formatDistanceToNow } from 'date-fns'

export default function ModerationPage() {
  const router = useRouter()
  const { user, isAuthenticated } = useAuthStore()
  const [videos, setVideos] = useState<Video[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [moderating, setModerating] = useState<number | null>(null)

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login')
      return
    }

    if (!user?.is_admin) {
      router.push('/')
      return
    }

    loadVideos()
  }, [isAuthenticated, user, router])

  const loadVideos = async () => {
    try {
      setLoading(true)
      const data = await videoApi.getPendingModeration(50, 0)
      setVideos(data)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка при загрузке видео')
    } finally {
      setLoading(false)
    }
  }

  const handleModerate = async (videoId: number, status: 'approved' | 'rejected') => {
    try {
      setModerating(videoId)
      await videoApi.moderateVideo(videoId, status)
      // Удаляем видео из списка после модерации
      setVideos((prev) => prev.filter((v) => v.id !== videoId))
    } catch (err: any) {
      alert(err.response?.data?.error || 'Ошибка при модерации')
    } finally {
      setModerating(null)
    }
  }

  const renderMedia = (video: Video) => {
    switch (video.media_type) {
      case 'video':
        return (
          <video
            src={video.media_url}
            poster={video.thumbnail_url}
            controls
            className="w-full h-auto max-h-96 rounded-lg"
          />
        )
      case 'gif':
      case 'photo':
        return (
          <img
            src={video.media_url}
            alt={video.title}
            className="w-full h-auto max-h-96 object-contain rounded-lg"
          />
        )
      default:
        return <div className="w-full h-64 bg-gray-200 flex items-center justify-center rounded-lg">Неизвестный тип медиа</div>
    }
  }

  if (!isAuthenticated || !user?.is_admin) {
    return null
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Загрузка видео на модерацию...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-100">
      <header className="bg-white shadow-sm py-4">
        <div className="container mx-auto px-4 flex justify-between items-center">
          <Link href="/" className="text-2xl font-bold text-blue-600">
            Sticky Stick
          </Link>
          <nav className="space-x-4">
            <Link href="/" className="text-gray-600 hover:text-blue-600">
              Главная
            </Link>
            <Link href="/moderation" className="text-blue-600 font-semibold">
              Модерация
            </Link>
            <span className="text-gray-600">Админ: {user.username}</span>
            <button
              onClick={() => {
                const { logout } = useAuthStore.getState()
                logout()
                router.push('/')
              }}
              className="text-red-600 hover:text-red-700 font-semibold"
            >
              Выйти
            </button>
          </nav>
        </div>
      </header>

      <div className="container mx-auto px-4 py-8">
        <div className="mb-6">
          <h1 className="text-3xl font-bold mb-2">Модерация видео</h1>
          <p className="text-gray-600">
            Видео на модерации: <span className="font-semibold">{videos.length}</span>
          </p>
        </div>

        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
            {error}
          </div>
        )}

        {videos.length === 0 ? (
          <div className="bg-white rounded-lg shadow p-8 text-center">
            <p className="text-gray-600 text-lg">Нет видео на модерации</p>
            <Link href="/" className="text-blue-600 hover:text-blue-800 mt-4 inline-block">
              Вернуться на главную
            </Link>
          </div>
        ) : (
          <div className="space-y-6">
            {videos.map((video) => (
              <div key={video.id} className="bg-white rounded-lg shadow-lg overflow-hidden">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6 p-6">
                  {/* Медиа */}
                  <div className="bg-black rounded-lg flex items-center justify-center">
                    {renderMedia(video)}
                  </div>

                  {/* Информация */}
                  <div className="flex flex-col">
                    <div className="mb-4">
                      <h2 className="text-2xl font-bold mb-2">{video.title}</h2>
                      {video.description && (
                        <p className="text-gray-600 mb-3">{video.description}</p>
                      )}
                      
                      <div className="flex items-center gap-4 text-sm text-gray-500 mb-3">
                        <span>@{video.user.username}</span>
                        <span>•</span>
                        <span>{video.views} просмотров</span>
                        <span>•</span>
                        <span>
                          {formatDistanceToNow(new Date(video.created_at), {
                            addSuffix: true,
                          })}
                        </span>
                      </div>

                      {video.category && (
                        <span className="inline-block bg-blue-100 text-blue-800 text-sm font-semibold px-3 py-1 rounded mb-2">
                          {video.category.name}
                        </span>
                      )}

                      {video.tags && (
                        <div className="flex flex-wrap gap-2 mt-2">
                          {video.tags.split(',').map((tag, idx) => (
                            <span
                              key={idx}
                              className="text-sm bg-gray-100 text-gray-700 px-3 py-1 rounded-full"
                            >
                              #{tag.trim()}
                            </span>
                          ))}
                        </div>
                      )}
                    </div>

                    {/* Кнопки модерации */}
                    <div className="mt-auto flex gap-3">
                      <button
                        onClick={() => handleModerate(video.id, 'approved')}
                        disabled={moderating === video.id}
                        className="flex-1 bg-green-600 text-white px-6 py-3 rounded-lg font-semibold hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                      >
                        {moderating === video.id ? 'Одобряю...' : '✓ Одобрить'}
                      </button>
                      <button
                        onClick={() => handleModerate(video.id, 'rejected')}
                        disabled={moderating === video.id}
                        className="flex-1 bg-red-600 text-white px-6 py-3 rounded-lg font-semibold hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                      >
                        {moderating === video.id ? 'Отклоняю...' : '✗ Отклонить'}
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
