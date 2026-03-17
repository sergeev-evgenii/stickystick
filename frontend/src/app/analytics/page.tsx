'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { adminApi, AnalyticsResponse, ActivityLogItem } from '@/lib/api/admin'
import { useAuthStore } from '@/store/authStore'
import { formatDistanceToNow } from 'date-fns'

const ACTION_LABELS: Record<string, string> = {
  login: 'Вход',
  register: 'Регистрация',
  video_view: 'Просмотр видео',
  like: 'Лайк',
  unlike: 'Снятие лайка',
  upload: 'Загрузка',
  feed_view: 'Просмотр ленты',
  generate_video_click: 'Нажатие «Сгенерировать своё видео»',
}

export default function AnalyticsPage() {
  const router = useRouter()
  const { user, isAuthenticated } = useAuthStore()
  const telegramUrl = process.env.NEXT_PUBLIC_TELEGRAM_URL
  const vkUrl = process.env.NEXT_PUBLIC_VK_URL
  const maxUrl = process.env.NEXT_PUBLIC_MAX_URL
  const [mounted, setMounted] = useState(false)
  const [data, setData] = useState<AnalyticsResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [period, setPeriod] = useState<'24h' | '7d' | '30d'>('24h')

  useEffect(() => { setMounted(true) }, [])

  useEffect(() => {
    if (!mounted) return
    if (!isAuthenticated) {
      router.push('/login')
      return
    }
    if (!user?.is_admin) {
      router.push('/')
      return
    }
    loadAnalytics()
  }, [mounted, isAuthenticated, user, router, period])

  const loadAnalytics = async () => {
    try {
      setLoading(true)
      setError('')
      const result = await adminApi.getAnalytics(period, 200, 0)
      setData(result)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка загрузки аналитики')
    } finally {
      setLoading(false)
    }
  }


  if (!mounted) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto" />
          <p className="mt-4 text-gray-600">Загрузка...</p>
        </div>
      </div>
    )
  }

  if (!isAuthenticated || !user?.is_admin) {
    return null
  }

  if (loading && !data) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto" />
          <p className="mt-4 text-gray-600">Загрузка аналитики...</p>
        </div>
      </div>
    )
  }

  if (error && !data) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-600 text-lg mb-4">{error}</p>
          <Link href="/" className="text-blue-600 hover:underline">На главную</Link>
        </div>
      </div>
    )
  }

  const stats = data?.stats
  const activity = data?.activity ?? []

  return (
    <div className="min-h-screen bg-gray-100">
      <header className="bg-white shadow-sm sticky top-0 z-10">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <Link href="/" className="flex flex-col font-verveine tracking-wider">
              <span className="text-4xl font-bold text-amber-500 uppercase">Sticky Stick</span>
              <span className="text-lg font-semibold text-amber-500 uppercase">липкая палка</span>
              <span className="text-base font-semibold text-black">просто мемы. просто смешно. иногда эпично</span>
            </Link>
            <nav className="flex items-center gap-4">
              {telegramUrl && (
                <a
                  href={telegramUrl}
                  target="_blank"
                  rel="noreferrer"
                  className="text-gray-600 hover:text-sky-600"
                  aria-label="Telegram"
                  title="Telegram"
                >
                  <svg className="w-6 h-6" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
                    <path d="M9.993 15.674 9.6 20.2c.563 0 .807-.243 1.1-.533l2.64-2.52 5.47 4.002c1.003.553 1.717.262 1.982-.93l3.59-16.82c.338-1.56-.563-2.17-1.55-1.804L1.726 9.17c-1.52.59-1.497 1.437-.259 1.82l5.65 1.763L20.24 4.73c.62-.41 1.185-.183.72.227L9.993 15.674z" />
                  </svg>
                </a>
              )}
              {vkUrl && (
                <a
                  href={vkUrl}
                  target="_blank"
                  rel="noreferrer"
                  className="text-gray-600 hover:text-blue-600"
                  aria-label="VK"
                  title="VK"
                >
                  <svg className="w-6 h-6" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
                    <path d="M12.784 16.58c.4 0 .57-.28.57-.66 0-.4-.01-.93.01-1.35.03-.78.25-1.32.86-1.32.59 0 1.14.63 1.62 1.22.53.65 1.05 1.32 1.62 1.9.42.42.76.21 1.03-.06l1.08-1.1c.26-.28.31-.63.07-.94-.31-.41-.77-.94-1.24-1.5-.55-.65-1.1-1.33-1.32-1.76-.22-.41-.12-.7.14-1.05.53-.7 1.2-1.55 1.72-2.38.26-.42.44-.78.6-1.19.17-.45-.03-.65-.48-.65h-1.58c-.41 0-.6.15-.72.46-.22.6-.56 1.3-1.08 2.12-.55.86-1.1 1.6-1.52 1.6-.29 0-.41-.32-.43-.82-.03-.71.02-1.7.02-2.32 0-.79-.24-1.04-.69-1.14-.22-.05-.53-.09-1.13-.09-.77 0-1.4.01-1.76.2-.25.13-.41.4-.3.41.37.04.5.23.58.51.12.41.12 1.33.12 1.98 0 .46.03 1.08-.1 1.33-.1.19-.2.26-.36.26-.38 0-.95-.72-1.51-1.6-.52-.82-.92-1.72-1.11-2.23-.12-.3-.27-.44-.68-.44H3.63c-.45 0-.6.23-.48.61.18.56.5 1.25.8 1.82.7 1.33 1.63 2.67 2.74 3.86 1.27 1.35 2.78 2.2 3.86 2.2z" />
                  </svg>
                </a>
              )}
              {maxUrl && (
                <a
                  href={maxUrl}
                  target="_blank"
                  rel="noreferrer"
                  className="text-gray-600 hover:text-violet-600"
                  aria-label="Max"
                  title="Max"
                >
                  <svg className="w-6 h-6" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
                    <path d="M4 6h3.2l2.1 6.2L11.4 6H14l2.1 6.2L18.2 6H21l-3.2 12h-2.7l-2.4-6.7L10.3 18H7.6L4 6z" />
                  </svg>
                </a>
              )}
              <span className="text-indigo-600 font-semibold">Аналитика</span>
              <Link href="/moderation" className="text-orange-600 hover:text-orange-700 font-semibold">
                Модерация
              </Link>
              <Link href="/" className="text-gray-600 hover:text-gray-800 font-semibold">
                На главную
              </Link>
            </nav>
          </div>
        </div>
      </header>

      <main className="container mx-auto px-4 py-8">
        <h1 className="text-2xl font-bold mb-6">Аналитика</h1>

        <div className="flex flex-wrap gap-2 mb-6">
          {(['24h', '7d', '30d'] as const).map((p) => (
            <button
              key={p}
              onClick={() => setPeriod(p)}
              className={`px-4 py-2 rounded-lg font-semibold transition-colors ${
                period === p ? 'bg-blue-600 text-white' : 'bg-white text-gray-700 hover:bg-gray-100'
              }`}
            >
              {p === '24h' ? '24 часа' : p === '7d' ? '7 дней' : '30 дней'}
            </button>
          ))}
        </div>

        {stats && (
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-7 gap-4 mb-8">
            <StatCard title="Уникальных IP" value={stats.unique_visitors} />
            <StatCard title="Просмотров видео" value={stats.total_video_views} />
            <StatCard title="Входов" value={stats.total_logins} />
            <StatCard title="Регистраций" value={stats.total_registrations} />
            <StatCard title="Лайков" value={stats.total_likes} />
            <StatCard title="Загрузок" value={stats.total_uploads} />
            <StatCard title="«Сгенерировать своё видео»" value={stats.total_generate_video_clicks ?? 0} />
          </div>
        )}

        <h2 className="text-xl font-bold mb-4">Последняя активность</h2>
        <div className="bg-white rounded-lg shadow overflow-hidden overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">IP</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Пользователь</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Действие</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Ресурс</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Когда</th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">User-Agent</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {activity.map((item) => (
                <ActivityRow key={item.id} item={item} />
              ))}
            </tbody>
          </table>
        </div>
        {activity.length === 0 && !loading && (
          <p className="text-gray-500 mt-4">Нет записей за выбранный период.</p>
        )}

      </main>
    </div>
  )
}

function StatCard({ title, value }: { title: string; value: number }) {
  return (
    <div className="bg-white rounded-lg shadow p-4">
      <p className="text-sm text-gray-500 mb-1">{title}</p>
      <p className="text-2xl font-bold text-gray-900">{value}</p>
    </div>
  )
}

function ActivityRow({ item }: { item: ActivityLogItem }) {
  const actionLabel = ACTION_LABELS[item.action] || item.action
  const resource = item.resource_type && item.resource_id ? `${item.resource_type} #${item.resource_id}` : '—'
  const userAgent = item.user_agent ? (item.user_agent.length > 60 ? item.user_agent.slice(0, 60) + '…' : item.user_agent) : '—'

  return (
    <tr className="hover:bg-gray-50">
      <td className="px-4 py-3 text-sm font-mono text-gray-900">{item.ip}</td>
      <td className="px-4 py-3 text-sm text-gray-700">
        {item.user ? `@${item.user.username}` : '—'}
      </td>
      <td className="px-4 py-3 text-sm text-gray-700">{actionLabel}</td>
      <td className="px-4 py-3 text-sm text-gray-600">{resource}</td>
      <td className="px-4 py-3 text-sm text-gray-600">
        {formatDistanceToNow(new Date(item.created_at), { addSuffix: true })}
      </td>
      <td className="px-4 py-3 text-xs text-gray-500 max-w-xs truncate" title={item.user_agent}>
        {userAgent}
      </td>
    </tr>
  )
}
