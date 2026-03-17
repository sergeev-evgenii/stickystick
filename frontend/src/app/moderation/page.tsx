'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { videoApi, Video } from '@/lib/api/video'
import { getMediaSrc } from '@/lib/utils/mediaUrl'
import { useAuthStore } from '@/store/authStore'
import { useSettingsStore } from '@/store/settingsStore'
import { formatDistanceToNow } from 'date-fns'

type Tab = 'pending' | 'approved' | 'hidden'

interface EditState {
  id: number
  title: string
  description: string
  tags: string
}

export default function ModerationPage() {
  const router = useRouter()
  const { user, isAuthenticated } = useAuthStore()
  const telegramUrl = process.env.NEXT_PUBLIC_TELEGRAM_URL
  const vkUrl = process.env.NEXT_PUBLIC_VK_URL
  const maxUrl = process.env.NEXT_PUBLIC_MAX_URL
  const [mounted, setMounted] = useState(false)
  const [tab, setTab] = useState<Tab>('pending')
  const [videos, setVideos] = useState<Video[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [actioning, setActioning] = useState<number | null>(null)
  const [edit, setEdit] = useState<EditState | null>(null)
  const [publishingVK, setPublishingVK] = useState<number | null>(null)
  const [publishingTelegram, setPublishingTelegram] = useState<number | null>(null)
  const [publishingMax, setPublishingMax] = useState<number | null>(null)
  const [settingsSaving, setSettingsSaving] = useState(false)
  const [defaultsSaving, setDefaultsSaving] = useState(false)
  const {
    showViewCount,
    defaultPublishVK,
    defaultPublishTelegram,
    defaultPublishMax,
    loaded: settingsLoaded,
    fetchSettings,
    setShowViewCount: setShowViewCountStore,
    setPublishDefaults,
  } = useSettingsStore()

  const [editDefaultVK, setEditDefaultVK] = useState('')
  const [editDefaultTelegram, setEditDefaultTelegram] = useState('')
  const [editDefaultMax, setEditDefaultMax] = useState('')

  useEffect(() => { setMounted(true) }, [])

  useEffect(() => {
    if (!mounted) return
    if (!isAuthenticated) { router.push('/login'); return }
    if (!user?.is_admin) { router.push('/'); return }
    loadVideos()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [mounted, isAuthenticated, user, tab])

  useEffect(() => {
    if (mounted && user?.is_admin && !settingsLoaded) fetchSettings()
  }, [mounted, user?.is_admin, settingsLoaded, fetchSettings])

  useEffect(() => {
    // синхронизируем локальные поля редактирования, когда подтянули настройки
    if (!settingsLoaded) return
    setEditDefaultVK(defaultPublishVK ?? '')
    setEditDefaultTelegram(defaultPublishTelegram ?? '')
    setEditDefaultMax(defaultPublishMax ?? '')
  }, [settingsLoaded, defaultPublishVK, defaultPublishTelegram, defaultPublishMax])

  const loadVideos = async () => {
    setLoading(true)
    setError('')
    try {
      let data: Video[] = []
      if (tab === 'pending') data = await videoApi.getPendingModeration(50, 0)
      else if (tab === 'approved') data = await videoApi.getApproved(50, 0)
      else data = await videoApi.getHidden(50, 0)
      setVideos(data ?? [])
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка загрузки')
    } finally {
      setLoading(false)
    }
  }

  const handleModerate = async (id: number, status: 'approved' | 'rejected') => {
    setActioning(id)
    try {
      await videoApi.moderateVideo(id, status)
      setVideos(prev => prev.filter(v => v.id !== id))
    } catch (err: any) {
      alert(err.response?.data?.error || 'Ошибка')
    } finally {
      setActioning(null)
    }
  }

  const handleHide = async (id: number) => {
    setActioning(id)
    try {
      await videoApi.hideVideo(id)
      setVideos(prev => prev.filter(v => v.id !== id))
    } catch (err: any) {
      alert(err.response?.data?.error || 'Ошибка')
    } finally {
      setActioning(null)
    }
  }

  const handleUnhide = async (id: number) => {
    setActioning(id)
    try {
      await videoApi.unhideVideo(id)
      setVideos(prev => prev.filter(v => v.id !== id))
    } catch (err: any) {
      alert(err.response?.data?.error || 'Ошибка')
    } finally {
      setActioning(null)
    }
  }

  const handlePublishVK = async (id: number) => {
    if (!confirm('Опубликовать в группу ВКонтакте?')) return
    setPublishingVK(id)
    try {
      const res = await videoApi.publishToVK(id)
      alert(`Опубликовано в ВК! ID поста: ${res.post_id}`)
    } catch (err: any) {
      alert(err.response?.data?.error || 'Ошибка публикации в ВК')
    } finally {
      setPublishingVK(null)
    }
  }

  const handlePublishTelegram = async (id: number) => {
    if (!confirm('Опубликовать в канал Telegram?')) return
    setPublishingTelegram(id)
    try {
      const res = await videoApi.publishToTelegram(id)
      alert(`Опубликовано в Telegram! ID сообщения: ${res.message_id}`)
    } catch (err: any) {
      alert(err.response?.data?.error || 'Ошибка публикации в Telegram')
    } finally {
      setPublishingTelegram(null)
    }
  }

  const handlePublishMax = async (id: number) => {
    if (!confirm('Опубликовать в мессенджер Max?')) return
    setPublishingMax(id)
    try {
      await videoApi.publishToMax(id)
      alert('Опубликовано в Max!')
    } catch (err: any) {
      alert(err.response?.data?.error || 'Ошибка публикации в Max')
    } finally {
      setPublishingMax(null)
    }
  }

  const handleToggleShowViewCount = async () => {
    setSettingsSaving(true)
    try {
      await setShowViewCountStore(!showViewCount)
    } catch (err: any) {
      alert(err.response?.data?.error || 'Ошибка сохранения настройки')
    } finally {
      setSettingsSaving(false)
    }
  }

  const handleSavePublishDefaults = async () => {
    setDefaultsSaving(true)
    try {
      await setPublishDefaults({
        vk: editDefaultVK,
        telegram: editDefaultTelegram,
        max: editDefaultMax,
      })
      alert('Сохранено')
    } catch (err: any) {
      alert(err.response?.data?.error || 'Ошибка сохранения дефолтных сообщений')
    } finally {
      setDefaultsSaving(false)
    }
  }

  const handleSaveEdit = async () => {
    if (!edit) return
    setActioning(edit.id)
    try {
      await videoApi.updateVideoFields(edit.id, {
        title: edit.title,
        description: edit.description,
        tags: edit.tags,
      })
      setVideos(prev => prev.map(v => v.id === edit.id
        ? { ...v, title: edit.title, description: edit.description, tags: edit.tags }
        : v
      ))
      setEdit(null)
    } catch (err: any) {
      alert(err.response?.data?.error || 'Ошибка сохранения')
    } finally {
      setActioning(null)
    }
  }

  const renderMedia = (video: Video) => {
    if (video.media_type === 'video') {
      return <video src={getMediaSrc(video.media_url)} poster={getMediaSrc(video.thumbnail_url)} controls className="w-full h-auto max-h-80 rounded-lg" />
    }
    return <img src={getMediaSrc(video.media_url)} alt={video.title} className="w-full h-auto max-h-80 object-contain rounded-lg" />
  }

  if (!mounted) return null
  if (!isAuthenticated || !user?.is_admin) return null

  const tabs: { id: Tab; label: string }[] = [
    { id: 'pending', label: 'На модерации' },
    { id: 'approved', label: 'В раздаче' },
    { id: 'hidden', label: 'Скрытые' },
  ]

  return (
    <div className="min-h-screen bg-gray-100">
      <header className="bg-white shadow-sm py-4 sticky top-0 z-10">
        <div className="container mx-auto px-4 flex justify-between items-center">
          <Link href="/" className="flex flex-col font-verveine tracking-wider">
            <span className="text-4xl font-bold text-amber-500 uppercase">Sticky Stick</span>
            <span className="text-lg font-semibold text-amber-500 uppercase">липкая палка</span>
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
            <Link href="/" className="text-gray-600 hover:text-blue-600">Главная</Link>
            <span className="text-blue-600 font-semibold">Модерация</span>
            <span className="text-gray-500">@{user.username}</span>
            <button
              onClick={() => { useAuthStore.getState().logout(); router.push('/') }}
              className="text-red-600 hover:text-red-700 font-semibold"
            >
              Выйти
            </button>
          </nav>
        </div>
      </header>

      <div className="container mx-auto px-4 py-8">
        {/* Вкладки */}
        <div className="flex gap-2 mb-6">
          {tabs.map(t => (
            <button
              key={t.id}
              onClick={() => setTab(t.id)}
              className={`px-5 py-2 rounded-lg font-semibold transition-colors ${
                tab === t.id
                  ? 'bg-blue-600 text-white'
                  : 'bg-white text-gray-700 hover:bg-gray-50 shadow-sm'
              }`}
            >
              {t.label}
            </button>
          ))}
        </div>

        {/* Настройки: показывать просмотры на пользовательской странице */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4 mb-6 flex items-center justify-between">
          <span className="text-sm font-medium text-gray-700">Показывать количество просмотров на пользовательской странице</span>
          <button
            type="button"
            role="switch"
            aria-checked={showViewCount}
            disabled={settingsSaving}
            onClick={handleToggleShowViewCount}
            className={`relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 ${
              showViewCount ? 'bg-blue-600' : 'bg-gray-200'
            }`}
          >
            <span
              className={`pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition ${
                showViewCount ? 'translate-x-5' : 'translate-x-1'
              }`}
            />
          </button>
        </div>

        {/* Настройки публикации: дефолтные сообщения */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4 mb-6">
          <p className="text-sm font-semibold text-gray-800 mb-3">Дефолтные сообщения для публикации (если комментарий не задан)</p>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label className="block text-xs font-semibold text-gray-600 mb-1">ВК</label>
              <textarea
                value={editDefaultVK}
                onChange={(e) => setEditDefaultVK(e.target.value)}
                rows={4}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div>
              <label className="block text-xs font-semibold text-gray-600 mb-1">Telegram</label>
              <textarea
                value={editDefaultTelegram}
                onChange={(e) => setEditDefaultTelegram(e.target.value)}
                rows={4}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div>
              <label className="block text-xs font-semibold text-gray-600 mb-1">Max</label>
              <textarea
                value={editDefaultMax}
                onChange={(e) => setEditDefaultMax(e.target.value)}
                rows={4}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>
          <div className="mt-4 flex justify-end">
            <button
              onClick={handleSavePublishDefaults}
              disabled={defaultsSaving}
              className="bg-blue-600 text-white px-4 py-2 rounded-lg text-sm font-semibold hover:bg-blue-700 disabled:opacity-50"
            >
              {defaultsSaving ? 'Сохраняю...' : 'Сохранить'}
            </button>
          </div>
        </div>

        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">{error}</div>
        )}

        {loading ? (
          <div className="flex items-center justify-center py-20">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600" />
          </div>
        ) : videos.length === 0 ? (
          <div className="bg-white rounded-lg shadow p-8 text-center text-gray-500">
            Нет видео в этом разделе
          </div>
        ) : (
          <div className="space-y-6">
            {videos.map(video => (
              <div key={video.id} className="bg-white rounded-lg shadow-lg overflow-hidden">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6 p-6">
                  {/* Медиа */}
                  <div className="bg-black rounded-lg flex items-center justify-center">
                    {renderMedia(video)}
                  </div>

                  {/* Инфо + действия */}
                  <div className="flex flex-col gap-3">
                    {edit?.id === video.id ? (
                      /* Режим редактирования */
                      <div className="flex flex-col gap-3">
                        <div>
                          <label className="block text-sm font-medium text-gray-700 mb-1">Заголовок</label>
                          <input
                            type="text"
                            value={edit.title}
                            onChange={e => setEdit({ ...edit, title: e.target.value })}
                            className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                          />
                        </div>
                        <div>
                          <label className="block text-sm font-medium text-gray-700 mb-1">Описание</label>
                          <textarea
                            value={edit.description}
                            onChange={e => setEdit({ ...edit, description: e.target.value })}
                            rows={3}
                            className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                          />
                        </div>
                        <div>
                          <label className="block text-sm font-medium text-gray-700 mb-1">
                            Теги <span className="text-gray-400 font-normal">(через запятую)</span>
                          </label>
                          <input
                            type="text"
                            value={edit.tags}
                            onChange={e => setEdit({ ...edit, tags: e.target.value })}
                            placeholder="юмор, смешное, приколы"
                            className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                          />
                        </div>
                        <div className="flex gap-2 mt-2">
                          <button
                            onClick={handleSaveEdit}
                            disabled={actioning === video.id}
                            className="flex-1 bg-blue-600 text-white px-4 py-2 rounded-lg font-semibold hover:bg-blue-700 disabled:opacity-50"
                          >
                            {actioning === video.id ? 'Сохраняю...' : 'Сохранить'}
                          </button>
                          <button
                            onClick={() => setEdit(null)}
                            className="flex-1 bg-gray-200 text-gray-700 px-4 py-2 rounded-lg font-semibold hover:bg-gray-300"
                          >
                            Отмена
                          </button>
                        </div>
                      </div>
                    ) : (
                      /* Режим просмотра */
                      <>
                        <div>
                          <h2 className="text-xl font-bold">{video.title}</h2>
                          {video.description && (
                            <p className="text-gray-600 text-sm mt-1">{video.description}</p>
                          )}
                        </div>

                        {video.tags && (
                          <div className="flex flex-wrap gap-1">
                            {video.tags.split(',').map((tag, i) => (
                              <span key={i} className="text-xs bg-gray-100 text-gray-600 px-2 py-1 rounded-full">
                                #{tag.trim()}
                              </span>
                            ))}
                          </div>
                        )}

                        <div className="text-sm text-gray-500 flex flex-wrap gap-3">
                          <span>@{video.user.username}</span>
                          <span>{video.views} просм.</span>
                          <span>{formatDistanceToNow(new Date(video.created_at), { addSuffix: true })}</span>
                        </div>

                        {/* Кнопка редактирования */}
                        <button
                          onClick={() => setEdit({
                            id: video.id,
                            title: video.title,
                            description: video.description || '',
                            tags: video.tags || '',
                          })}
                          className="w-full bg-gray-100 text-gray-700 px-4 py-2 rounded-lg font-semibold hover:bg-gray-200 transition-colors text-sm"
                        >
                          ✏️ Редактировать текст и теги
                        </button>

                        {/* Публикация в ВК */}
                        <button
                          onClick={() => handlePublishVK(video.id)}
                          disabled={publishingVK === video.id}
                          className="w-full bg-blue-500 text-white px-4 py-2 rounded-lg font-semibold hover:bg-blue-600 disabled:opacity-50 transition-colors text-sm"
                        >
                          {publishingVK === video.id ? 'Публикую...' : '📤 Опубликовать в ВК'}
                        </button>

                        {/* Публикация в Telegram */}
                        <button
                          onClick={() => handlePublishTelegram(video.id)}
                          disabled={publishingTelegram === video.id}
                          className="w-full bg-sky-500 text-white px-4 py-2 rounded-lg font-semibold hover:bg-sky-600 disabled:opacity-50 transition-colors text-sm"
                        >
                          {publishingTelegram === video.id ? 'Публикую...' : '📤 Опубликовать в Telegram'}
                        </button>

                        {/* Публикация в Max */}
                        <button
                          onClick={() => handlePublishMax(video.id)}
                          disabled={publishingMax === video.id}
                          className="w-full bg-violet-500 text-white px-4 py-2 rounded-lg font-semibold hover:bg-violet-600 disabled:opacity-50 transition-colors text-sm"
                        >
                          {publishingMax === video.id ? 'Публикую...' : '📤 Опубликовать в Max'}
                        </button>

                        {/* Кнопки модерации */}
                        <div className="flex flex-col gap-2 mt-auto">
                          {tab === 'pending' && (
                            <div className="flex gap-2">
                              <button
                                onClick={() => handleModerate(video.id, 'approved')}
                                disabled={actioning === video.id}
                                className="flex-1 bg-green-600 text-white px-4 py-2 rounded-lg font-semibold hover:bg-green-700 disabled:opacity-50 transition-colors"
                              >
                                {actioning === video.id ? '...' : '✓ Одобрить'}
                              </button>
                              <button
                                onClick={() => handleModerate(video.id, 'rejected')}
                                disabled={actioning === video.id}
                                className="flex-1 bg-red-600 text-white px-4 py-2 rounded-lg font-semibold hover:bg-red-700 disabled:opacity-50 transition-colors"
                              >
                                {actioning === video.id ? '...' : '✗ Отклонить'}
                              </button>
                            </div>
                          )}

                          {tab === 'approved' && (
                            <button
                              onClick={() => handleHide(video.id)}
                              disabled={actioning === video.id}
                              className="w-full bg-orange-500 text-white px-4 py-2 rounded-lg font-semibold hover:bg-orange-600 disabled:opacity-50 transition-colors"
                            >
                              {actioning === video.id ? '...' : '🚫 Убрать из раздачи'}
                            </button>
                          )}

                          {tab === 'hidden' && (
                            <button
                              onClick={() => handleUnhide(video.id)}
                              disabled={actioning === video.id}
                              className="w-full bg-green-600 text-white px-4 py-2 rounded-lg font-semibold hover:bg-green-700 disabled:opacity-50 transition-colors"
                            >
                              {actioning === video.id ? '...' : '✓ Вернуть в раздачу'}
                            </button>
                          )}
                        </div>
                      </>
                    )}
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
