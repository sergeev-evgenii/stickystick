'use client'

import { useEffect, useState, useRef, Suspense } from 'react'
import { useParams, useRouter, useSearchParams } from 'next/navigation'
import Link from 'next/link'
import { videoApi, Video } from '@/lib/api/video'
import { getMediaSrc } from '@/lib/utils/mediaUrl'
import VideoDetail from '@/components/VideoDetail'

const VIDEO_NAV_STACK_KEY = 'videoNavStack'

function getNavStack(): number[] {
  if (typeof window === 'undefined') return []
  try {
    const raw = sessionStorage.getItem(VIDEO_NAV_STACK_KEY)
    return raw ? JSON.parse(raw) : []
  } catch {
    return []
  }
}

function pushNavStack(id: number) {
  const stack = getNavStack()
  stack.push(id)
  sessionStorage.setItem(VIDEO_NAV_STACK_KEY, JSON.stringify(stack))
}

function popNavStack(): number | null {
  const stack = getNavStack()
  if (stack.length === 0) return null
  const id = stack.pop()!
  sessionStorage.setItem(VIDEO_NAV_STACK_KEY, JSON.stringify(stack))
  return id
}

function VideoPageContent() {
  const params = useParams()
  const router = useRouter()
  const searchParams = useSearchParams()
  const idParam = params?.id as string | undefined
  const [video, setVideo] = useState<Video | null>(null)
  const [videos, setVideos] = useState<Video[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const touchStartY = useRef<number | null>(null)
  const touchEndY = useRef<number | null>(null)
  const minSwipeDistance = 50

  // Сбрасываем стек, если зашли на видео не через кнопки Назад/Вперёд (например с ленты)
  useEffect(() => {
    if (typeof window === 'undefined') return
    const nav = searchParams.get('nav')
    if (nav !== 'next' && nav !== 'prev') {
      sessionStorage.removeItem(VIDEO_NAV_STACK_KEY)
    }
  }, [searchParams])

  useEffect(() => {
    const fetchVideo = async () => {
      if (!idParam) {
        setError('ID видео не указан')
        setLoading(false)
        return
      }
      const id = parseInt(idParam, 10)
      if (isNaN(id)) {
        setError('Неверный ID видео')
        setLoading(false)
        return
      }
      try {
        const [currentVideo, feedVideos] = await Promise.all([
          videoApi.getVideo(id),
          videoApi.getFeed(10, 0), // лента непросмотренных; при просмотре всего бэкенд сбрасывает и отдаёт в случайном порядке
        ])
        setVideo(currentVideo)
        setVideos(feedVideos ?? [])
      } catch (err: any) {
        setError(err.response?.data?.error || 'Видео не найдено')
      } finally {
        setLoading(false)
      }
    }
    fetchVideo()
  }, [idParam])

  // SEO: meta title, description, Schema.org VideoObject
  useEffect(() => {
    if (!video) return
    const title = `${video.title} — Смешные короткие видео | Sticky Stick`
    const description =
      video.description?.trim() ||
      `Смотри смешное видео: ${video.title}. Короткие смешные ролики, юмор, приколы.`
    document.title = title

    let metaDesc = document.querySelector('meta[name="description"]')
    if (!metaDesc) {
      metaDesc = document.createElement('meta')
      metaDesc.setAttribute('name', 'description')
      document.head.appendChild(metaDesc)
    }
    metaDesc.setAttribute('content', description)

    const scriptId = 'video-schema'
    const existing = document.getElementById(scriptId)
    if (existing) existing.remove()
    const script = document.createElement('script')
    script.id = scriptId
    script.type = 'application/ld+json'
    script.textContent = JSON.stringify({
      '@context': 'https://schema.org',
      '@type': 'VideoObject',
      name: video.title,
      description: description,
      thumbnailUrl: getMediaSrc(video.thumbnail_url || video.media_url),
      contentUrl: getMediaSrc(video.media_url),
      uploadDate: video.created_at,
      duration: video.duration ? `PT${video.duration}S` : undefined,
      ...(video.tags && { keywords: video.tags.split(',').map((t: string) => t.trim()).filter(Boolean).join(', ') }),
    })
    document.head.appendChild(script)
    return () => {
      const s = document.getElementById(scriptId)
      if (s) s.remove()
    }
  }, [video])

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
    if (Math.abs(distance) < minSwipeDistance) return
    const currentIndex = videos.findIndex((v) => v.id === video.id)
    if (distance > 0) {
      pushNavStack(video.id)
      const nextIndex = currentIndex < videos.length - 1 ? currentIndex + 1 : 0
      router.push(`/videos/${videos[nextIndex].id}/?nav=next`)
    } else {
      const backId = popNavStack()
      if (backId != null) {
        router.push(`/videos/${backId}/?nav=prev`)
      } else {
        const prevIndex = currentIndex > 0 ? currentIndex - 1 : videos.length - 1
        router.push(`/videos/${videos[prevIndex].id}/`)
      }
    }
  }

  const handleNext = () => {
    if (!video || videos.length === 0) return
    pushNavStack(video.id)
    const currentIndex = videos.findIndex((v) => v.id === video.id)
    const nextIndex = currentIndex < videos.length - 1 ? currentIndex + 1 : 0
    router.push(`/videos/${videos[nextIndex].id}/?nav=next`)
  }
  const handlePrev = () => {
    if (!video || videos.length === 0) return
    const backId = popNavStack()
    if (backId != null) {
      router.push(`/videos/${backId}/?nav=prev`)
      return
    }
    const currentIndex = videos.findIndex((v) => v.id === video.id)
    const prevIndex = currentIndex > 0 ? currentIndex - 1 : videos.length - 1
    router.push(`/videos/${videos[prevIndex].id}/`)
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto" />
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
          <Link href="/" className="bg-blue-600 text-white px-6 py-3 rounded-lg font-semibold hover:bg-blue-700 transition-colors">
            Вернуться на главную
          </Link>
        </div>
      </div>
    )
  }

  const hasVideos = videos.length > 0

  return (
    <div
      className="h-screen bg-black flex flex-col overflow-hidden md:bg-gray-100 md:flex-row md:items-center md:justify-center"
      onTouchStart={handleTouchStart}
      onTouchMove={handleTouchMove}
      onTouchEnd={handleTouchEnd}
    >
      <Link
        href="/"
        className="absolute top-4 left-4 z-30 px-3 py-2 rounded-lg bg-black/50 text-white text-sm font-semibold hover:bg-black/70 transition-colors md:top-6 md:left-6 md:bg-white md:text-gray-800 md:hover:bg-gray-100"
        aria-label="Назад к ленте"
      >
        ← Назад
      </Link>

      <div className="flex-1 flex items-center justify-center min-h-0 w-full p-0 md:p-4 md:max-w-4xl md:max-h-[90vh] md:w-full">
        <div className="relative w-full h-full max-h-full md:max-h-[90vh] flex items-center justify-center">
          <VideoDetail video={video} onUpdate={setVideo} />

          {hasVideos && (
            <>
              <button
                onClick={handlePrev}
                className="absolute left-2 top-[38%] -translate-y-1/2 z-20 bg-white/90 rounded-full p-2 shadow-lg hover:bg-white transition-colors"
                aria-label="Предыдущее видео"
              >
                <svg className="w-5 h-5 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                </svg>
              </button>
              <button
                onClick={handleNext}
                className="absolute right-2 top-[38%] -translate-y-1/2 z-20 bg-white/90 rounded-full p-2 shadow-lg hover:bg-white transition-colors"
                aria-label="Следующее видео"
              >
                <svg className="w-5 h-5 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                </svg>
              </button>
            </>
          )}
        </div>
      </div>
    </div>
  )
}

export default function VideoPageClient() {
  return (
    <Suspense
      fallback={
        <div className="min-h-screen bg-gray-100 flex items-center justify-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600" />
        </div>
      }
    >
      <VideoPageContent />
    </Suspense>
  )
}
