'use client'

import { useEffect, useState, useRef, Suspense } from 'react'
import { useSearchParams, useRouter } from 'next/navigation'
import { videoApi, Video } from '@/lib/api/video'
import VideoDetail from '@/components/VideoDetail'
import Link from 'next/link'

function VideoPageContent() {
  const searchParams = useSearchParams()
  const router = useRouter()
  const [video, setVideo] = useState<Video | null>(null)
  const [videos, setVideos] = useState<Video[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const touchStartY = useRef<number | null>(null)
  const touchEndY = useRef<number | null>(null)
  const minSwipeDistance = 50

  useEffect(() => {
    const fetchVideo = async () => {
      try {
        const idParam = searchParams.get('id')
        if (!idParam) {
          setError('ID видео не указан')
          setLoading(false)
          return
        }

        const id = parseInt(idParam)
        if (isNaN(id)) {
          setError('Неверный ID видео')
          setLoading(false)
          return
        }

        // Загружаем текущее видео и список видео для навигации
        const [currentVideo, feedVideos] = await Promise.all([
          videoApi.getVideo(id),
          videoApi.getFeed(50, 0) // Загружаем больше видео для навигации
        ])
        
        setVideo(currentVideo)
        setVideos(feedVideos)
      } catch (err: any) {
        setError(err.response?.data?.error || 'Видео не найдено')
      } finally {
        setLoading(false)
      }
    }

    fetchVideo()
  }, [searchParams])

  const handleTouchStart = (e: React.TouchEvent) => {
    touchEndY.current = null
    touchStartY.current = e.targetTouches[0].clientY
  }

  const handleTouchMove = (e: React.TouchEvent) => {
    touchEndY.current = e.targetTouches[0].clientY
  }

  const handleTouchEnd = () => {
    if (!touchStartY.current || !touchEndY.current || !video || videos.length === 0) return
    
    const distance = touchStartY.current - touchEndY.current
    
    if (Math.abs(distance) < minSwipeDistance) {
      return
    }

    const currentIndex = videos.findIndex(v => v.id === video.id)
    
    if (distance > 0) {
      // Свайп вверх - следующее видео (зацикливание)
      const nextIndex = currentIndex < videos.length - 1 ? currentIndex + 1 : 0
      const nextVideo = videos[nextIndex]
      router.push(`/videos?id=${nextVideo.id}`)
    } else {
      // Свайп вниз - предыдущее видео (зацикливание)
      const prevIndex = currentIndex > 0 ? currentIndex - 1 : videos.length - 1
      const prevVideo = videos[prevIndex]
      router.push(`/videos?id=${prevVideo.id}`)
    }
  }

  const handleNext = () => {
    if (!video || videos.length === 0) return
    const currentIndex = videos.findIndex(v => v.id === video.id)
    // Зацикливание: если последнее видео, переходим к первому
    const nextIndex = currentIndex < videos.length - 1 ? currentIndex + 1 : 0
    const nextVideo = videos[nextIndex]
    router.push(`/videos?id=${nextVideo.id}`)
  }

  const handlePrev = () => {
    if (!video || videos.length === 0) return
    const currentIndex = videos.findIndex(v => v.id === video.id)
    // Зацикливание: если первое видео, переходим к последнему
    const prevIndex = currentIndex > 0 ? currentIndex - 1 : videos.length - 1
    const prevVideo = videos[prevIndex]
    router.push(`/videos?id=${prevVideo.id}`)
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Загрузка видео...</p>
        </div>
      </div>
    )
  }

  if (error || !video) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-600 text-lg mb-4">{error || 'Видео не найдено'}</p>
          <Link
            href="/"
            className="bg-blue-600 text-white px-6 py-3 rounded-lg font-semibold hover:bg-blue-700 transition-colors"
          >
            Вернуться на главную
          </Link>
        </div>
      </div>
    )
  }

  const currentIndex = video ? videos.findIndex(v => v.id === video.id) : -1
  const hasVideos = videos.length > 0

  return (
    <div 
      className="min-h-screen bg-gray-100 py-8"
      onTouchStart={handleTouchStart}
      onTouchMove={handleTouchMove}
      onTouchEnd={handleTouchEnd}
    >
      <div className="container mx-auto px-4">
        <Link
          href="/"
          className="inline-block mb-4 text-blue-600 hover:text-blue-800 font-semibold"
        >
          ← Назад к ленте
        </Link>
        <div className="max-w-4xl mx-auto relative">
          <VideoDetail video={video} onUpdate={setVideo} />
          
          {/* Кнопка назад - всегда показываем, если есть видео */}
          {hasVideos && (
            <button
              onClick={handlePrev}
              className="absolute left-2 top-1/2 -translate-y-1/2 z-20 bg-white rounded-full p-2 shadow-lg hover:bg-gray-100 transition-colors"
              aria-label="Предыдущее видео"
            >
              <svg className="w-5 h-5 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
            </button>
          )}
          
          {/* Кнопка вперед - всегда показываем, если есть видео */}
          {hasVideos && (
            <button
              onClick={handleNext}
              className="absolute right-2 top-1/2 -translate-y-1/2 z-20 bg-white rounded-full p-2 shadow-lg hover:bg-gray-100 transition-colors"
              aria-label="Следующее видео"
            >
              <svg className="w-5 h-5 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
            </button>
          )}
        </div>
      </div>
    </div>
  )
}

export default function VideoPage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Загрузка...</p>
        </div>
      </div>
    }>
      <VideoPageContent />
    </Suspense>
  )
}
