'use client'

import { Video } from '@/lib/api/video'

interface VideoPlayerProps {
  video: Video
}

export default function VideoPlayer({ video }: VideoPlayerProps) {
  const renderMedia = () => {
    switch (video.media_type) {
      case 'video':
        return (
          <video
            src={video.media_url}
            poster={video.thumbnail_url}
            controls
            className="w-full h-full object-contain"
          />
        )
      case 'gif':
        return (
          <img
            src={video.media_url}
            alt={video.title}
            className="w-full h-full object-contain"
          />
        )
      case 'photo':
        return (
          <img
            src={video.media_url}
            alt={video.title}
            className="w-full h-full object-contain"
          />
        )
      default:
        return <div className="w-full h-full bg-gray-200 flex items-center justify-center">Неизвестный тип медиа</div>
    }
  }

  return (
    <div className="relative w-full bg-black rounded-lg overflow-hidden">
      <div className="relative w-full aspect-video bg-black">
        {renderMedia()}
      </div>
      <div className="p-4 text-white">
        <h3 className="font-bold text-lg mb-2">{video.title}</h3>
        {video.description && (
          <p className="text-sm text-gray-300 mb-2">{video.description}</p>
        )}
        <div className="flex items-center justify-between text-sm">
          <span>@{video.user.username}</span>
          <span>{video.views} просмотров</span>
        </div>
      </div>
    </div>
  )
}
