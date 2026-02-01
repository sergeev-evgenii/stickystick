'use client'

import { useEffect, useState, Suspense } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { videoApi, Video } from '@/lib/api/video'
import { categoryApi, Category } from '@/lib/api/category'
import VideoCard from './VideoCard'

function VideoFeedContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const [videos, setVideos] = useState<Video[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [selectedCategory, setSelectedCategory] = useState<number | undefined>(undefined)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [offset, setOffset] = useState(0)
  const [hasMore, setHasMore] = useState(true)

  useEffect(() => {
    // Загружаем категории
    categoryApi.getAll().then((cats) => {
      setCategories(cats)
      
      // Проверяем параметр категории из URL
      const categoryParam = searchParams.get('category')
      if (categoryParam) {
        const category = cats.find(c => c.slug === categoryParam || c.id.toString() === categoryParam)
        if (category) {
          setSelectedCategory(category.id)
        }
      }
    }).catch(console.error)
  }, [searchParams])

  const loadVideos = async (reset = false) => {
    try {
      setLoading(true)
      const currentOffset = reset ? 0 : offset
      const newVideos = await videoApi.getFeed(20, currentOffset, selectedCategory)

      if (reset) {
        setVideos(newVideos)
      } else {
        setVideos((prev) => [...prev, ...newVideos])
      }

      setOffset(currentOffset + newVideos.length)
      setHasMore(newVideos.length === 20)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка при загрузке видео')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    setOffset(0)
    loadVideos(true)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedCategory])

  if (loading && videos.length === 0) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Загрузка видео...</p>
        </div>
      </div>
    )
  }

  if (error && videos.length === 0) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-600 text-lg mb-4">{error}</p>
          <button
            onClick={() => loadVideos(true)}
            className="bg-blue-600 text-white px-6 py-3 rounded-lg font-semibold hover:bg-blue-700"
          >
            Попробовать снова
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-100 py-8">
      <div className="container mx-auto px-4">
        {/* Фильтр по категориям */}
        {categories.length > 0 && (
          <div className="mb-6">
            <div className="flex flex-wrap gap-2">
              <button
                onClick={() => setSelectedCategory(undefined)}
                className={`px-4 py-2 rounded-lg font-semibold transition-colors ${
                  selectedCategory === undefined
                    ? 'bg-blue-600 text-white'
                    : 'bg-white text-gray-700 hover:bg-gray-100'
                }`}
              >
                Все
              </button>
              {categories.map((category) => (
                <button
                  key={category.id}
                  onClick={() => setSelectedCategory(category.id)}
                  className={`px-4 py-2 rounded-lg font-semibold transition-colors ${
                    selectedCategory === category.id
                      ? 'bg-blue-600 text-white'
                      : 'bg-white text-gray-700 hover:bg-gray-100'
                  }`}
                >
                  {category.name}
                </button>
              ))}
            </div>
          </div>
        )}

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {videos.map((video) => (
            <VideoCard
              key={video.id}
              video={video}
              onClick={() => router.push(`/videos?id=${video.id}`)}
            />
          ))}
        </div>

        {hasMore && (
          <div className="text-center mt-8">
            <button
              onClick={() => loadVideos(false)}
              disabled={loading}
              className="bg-blue-600 text-white px-6 py-3 rounded-lg font-semibold hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading ? 'Загрузка...' : 'Загрузить еще'}
            </button>
          </div>
        )}

        {!hasMore && videos.length > 0 && (
          <div className="text-center mt-8 text-gray-500">
            Все видео загружены
          </div>
        )}
      </div>
    </div>
  )
}

export default function VideoFeed() {
  return (
    <Suspense fallback={
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Загрузка...</p>
        </div>
      </div>
    }>
      <VideoFeedContent />
    </Suspense>
  )
}
