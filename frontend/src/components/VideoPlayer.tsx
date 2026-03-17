'use client'

import { Video } from '@/lib/api/video'
import { getMediaSrc } from '@/lib/utils/mediaUrl'
import { useSettingsStore } from '@/store/settingsStore'

interface VideoPlayerProps {
  video: Video
}

export default function VideoPlayer({ video }: VideoPlayerProps) {
  const showViewCount = useSettingsStore((s) => s.showViewCount)
  const renderMedia = () => {
    switch (video.media_type) {
      case 'video':
        return (
          <video
            src={getMediaSrc(video.media_url)}
            poster={getMediaSrc(video.thumbnail_url)}
            controls
            className="w-full h-full object-contain"
          />
        )
      case 'gif':
        return (
          <img
            src={getMediaSrc(video.media_url)}
            alt={video.title}
            className="w-full h-full object-contain"
          />
        )
      case 'photo':
        return (
          <img
            src={getMediaSrc(video.media_url)}
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
          {showViewCount && <span>{video.views} просмотров</span>}
        </div>
      </div>
    </div>
  )
}
