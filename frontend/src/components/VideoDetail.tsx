'use client'

import { useState, useEffect, useRef } from 'react'
import { Video, videoApi } from '@/lib/api/video'
import { useAuthStore } from '@/store/authStore'
import { formatDistanceToNow } from 'date-fns'

interface VideoDetailProps {
  video: Video
  onUpdate?: (video: Video) => void
}


export default function VideoDetail({ video, onUpdate }: VideoDetailProps) {
  const { isAuthenticated, user } = useAuthStore()
  const [isLiked, setIsLiked] = useState(false)
  const [likesCount, setLikesCount] = useState(video.likes?.length || 0)
  const [loading, setLoading] = useState(false)
  const videoRef = useRef<HTMLVideoElement>(null)

  useEffect(() => {
    // Проверяем, лайкнул ли текущий пользователь видео
    if (isAuthenticated && user && video.likes) {
      const liked = video.likes.some((like: any) => like.user_id === user.id)
      setIsLiked(liked)
    }
    setLikesCount(video.likes?.length || 0)
  }, [video, isAuthenticated, user])

  const handleLike = async () => {
    if (!isAuthenticated) {
      return
    }

    setLoading(true)
    try {
      if (isLiked) {
        await videoApi.unlikeVideo(video.id)
        setIsLiked(false)
        setLikesCount((prev) => Math.max(0, prev - 1))
      } else {
        await videoApi.likeVideo(video.id)
        setIsLiked(true)
        setLikesCount((prev) => prev + 1)
      }
    } catch (error) {
      console.error('Error toggling like:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleRepeat = () => {
    if (videoRef.current) {
      videoRef.current.currentTime = 0
      videoRef.current.play()
    }
  }


  const renderMedia = () => {
    switch (video.media_type) {
      case 'video':
        return (
          <video
            ref={videoRef}
            src={video.media_url}
            poster={video.thumbnail_url}
            controls
            className="w-full h-full object-contain rounded-lg"
            autoPlay
          />
        )
      case 'gif':
        return (
          <img
            src={video.media_url}
            alt={video.title}
            className="w-full h-full object-contain rounded-lg"
          />
        )
      case 'photo':
        return (
          <img
            src={video.media_url}
            alt={video.title}
            className="w-full h-full object-contain rounded-lg"
          />
        )
      default:
        return (
          <div className="w-full h-full bg-gray-200 flex items-center justify-center rounded-lg">
            Неизвестный тип медиа
          </div>
        )
    }
  }

  return (
    <div className="bg-white rounded-lg shadow-lg overflow-hidden select-none">
      {/* Медиа */}
      <div className="relative bg-black w-full aspect-video flex items-center justify-center touch-none">
        {renderMedia()}
      </div>

      {/* Информация о видео */}
      <div className="p-6">
        <div className="flex items-start justify-between mb-4">
          <div className="flex-1">
            {video.category && (
              <span className="inline-block bg-blue-100 text-blue-800 text-sm font-semibold px-3 py-1 rounded mb-2">
                {video.category.name}
              </span>
            )}
            <h1 className="text-2xl font-bold mb-2">{video.title}</h1>
            {video.description && (
              <p className="text-gray-600 mb-4">{video.description}</p>
            )}
            {video.tags && (
              <div className="flex flex-wrap gap-2 mb-4">
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
            <div className="flex items-center gap-4 text-sm text-gray-500">
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
          </div>
        </div>

        {/* Действия */}
        <div className="flex items-center gap-4 mb-6 pb-6 border-b">
          <button
            onClick={handleLike}
            disabled={!isAuthenticated || loading}
            className={`flex items-center gap-2 px-4 py-2 rounded-lg font-semibold transition-colors ${
              isLiked
                ? 'bg-red-500 text-white hover:bg-red-600'
                : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            } disabled:opacity-50 disabled:cursor-not-allowed`}
          >
            <svg
              className="w-5 h-5"
              fill={isLiked ? 'currentColor' : 'none'}
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
              />
            </svg>
            {likesCount}
          </button>
          {video.media_type === 'video' && (
            <button
              onClick={handleRepeat}
              className="flex items-center gap-2 px-4 py-2 rounded-lg font-semibold transition-colors bg-gray-200 text-gray-700 hover:bg-gray-300"
              aria-label="Повторить видео"
            >
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
                />
              </svg>
              Повторить
            </button>
          )}
        </div>

      </div>
    </div>
  )
}
