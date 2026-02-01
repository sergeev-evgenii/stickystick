'use client'

import { Video } from '@/lib/api/video'
import { formatDistanceToNow } from 'date-fns'
import Link from 'next/link'

interface VideoCardProps {
  video: Video
  onClick?: () => void
}

export default function VideoCard({ video, onClick }: VideoCardProps) {
  const renderThumbnail = () => {
    if (video.thumbnail_url && video.media_type === 'video') {
      return (
        <img
          src={video.thumbnail_url}
          alt={video.title}
          className="w-full h-48 object-cover"
        />
      )
    }

    switch (video.media_type) {
      case 'video':
        return (
          <div className="w-full h-48 bg-gray-800 flex items-center justify-center">
            <svg
              className="w-16 h-16 text-white"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path d="M6.3 2.841A1.5 1.5 0 004 4.11V15.89a1.5 1.5 0 002.3 1.269l9.344-5.89a1.5 1.5 0 000-2.538L6.3 2.84z" />
            </svg>
          </div>
        )
      case 'gif':
      case 'photo':
        return (
          <img
            src={video.media_url}
            alt={video.title}
            className="w-full h-48 object-cover"
          />
        )
      default:
        return (
          <div className="w-full h-48 bg-gray-200 flex items-center justify-center">
            Медиа
          </div>
        )
    }
  }

  return (
    <div
      className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow cursor-pointer"
      onClick={onClick}
    >
      <div className="relative">
        {renderThumbnail()}
        {video.media_type === 'video' && video.duration > 0 && (
          <div className="absolute bottom-2 right-2 bg-black bg-opacity-75 text-white text-xs px-2 py-1 rounded">
            {Math.floor(video.duration / 60)}:
            {String(video.duration % 60).padStart(2, '0')}
          </div>
        )}
      </div>
      <div className="p-4">
        {video.category && (
          <span className="inline-block bg-blue-100 text-blue-800 text-xs font-semibold px-2 py-1 rounded mb-2">
            {video.category.name}
          </span>
        )}
        <h3 className="font-bold text-lg mb-2 line-clamp-2">{video.title}</h3>
        {video.description && (
          <p className="text-gray-600 text-sm mb-3 line-clamp-2">
            {video.description}
          </p>
        )}
        {video.tags && (
          <div className="flex flex-wrap gap-1 mb-3">
            {video.tags.split(',').slice(0, 3).map((tag, idx) => (
              <span
                key={idx}
                className="text-xs bg-gray-100 text-gray-700 px-2 py-1 rounded"
              >
                #{tag.trim()}
              </span>
            ))}
          </div>
        )}
        <div className="flex items-center justify-between text-sm text-gray-500">
          <Link
            href={`/users/${video.user.id}`}
            onClick={(e) => e.stopPropagation()}
            className="hover:text-blue-600"
          >
            @{video.user.username}
          </Link>
          <div className="flex items-center gap-4">
            <span>{video.views} просмотров</span>
            <span>
              {formatDistanceToNow(new Date(video.created_at), {
                addSuffix: true,
              })}
            </span>
          </div>
        </div>
        {video.likes && video.likes.length > 0 && (
          <div className="mt-2 flex items-center gap-1 text-sm text-gray-500">
            <svg className="w-4 h-4 text-red-500" fill="currentColor" viewBox="0 0 20 20">
              <path
                fillRule="evenodd"
                d="M3.172 5.172a4 4 0 015.656 0L10 6.343l1.172-1.171a4 4 0 115.656 5.656L10 17.657l-6.828-6.829a4 4 0 010-5.656z"
                clipRule="evenodd"
              />
            </svg>
            <span>{video.likes.length}</span>
          </div>
        )}
      </div>
    </div>
  )
}
