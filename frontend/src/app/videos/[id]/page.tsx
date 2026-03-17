import VideoPageClient from './VideoPageClient'

export async function generateStaticParams() {
  try {
    const base = process.env.NEXT_PUBLIC_API_URL || 'https://stickystick.ru'
    const res = await fetch(`${base}/api/v1/videos?limit=10&offset=0`, { next: { revalidate: 3600 } })
    const videos = await res.json()
    if (Array.isArray(videos)) {
      return videos.map((v: { id: number }) => ({ id: String(v.id) }))
    }
  } catch {
    return Array.from({ length: 50 }, (_, i) => ({ id: String(i + 1) }))
  }
  return []
}

export default function VideoPage() {
  return <VideoPageClient />
}
