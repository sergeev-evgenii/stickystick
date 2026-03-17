'use client'

import { useState, useEffect } from 'react'
import { videoApi, UploadMediaData } from '@/lib/api/video'
import { categoryApi, Category } from '@/lib/api/category'

interface MediaUploadProps {
  onSuccess?: (video: any) => void
  onError?: (error: string) => void
}

export default function MediaUpload({ onSuccess, onError }: MediaUploadProps) {
  const [file, setFile] = useState<File | null>(null)
  const [thumbnail, setThumbnail] = useState<File | null>(null)
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [tags, setTags] = useState('')
  const [categoryId, setCategoryId] = useState<number | undefined>(undefined)
  const [categories, setCategories] = useState<Category[]>([])
  const [uploading, setUploading] = useState(false)
  const [uploadProgress, setUploadProgress] = useState(0)
  const [preview, setPreview] = useState<string | null>(null)

  useEffect(() => {
    // Загружаем список категорий
    categoryApi.getAll().then(setCategories).catch(console.error)
  }, [])

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0]
    if (!selectedFile) return

    setFile(selectedFile)

    // Создаем превью
    const reader = new FileReader()
    reader.onloadend = () => {
      setPreview(reader.result as string)
    }
    
    const fileType = selectedFile.type
    if (fileType.startsWith('image/')) {
      reader.readAsDataURL(selectedFile)
    } else if (fileType.startsWith('video/')) {
      reader.readAsDataURL(selectedFile)
    }
  }

  const handleThumbnailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0]
    if (selectedFile) {
      setThumbnail(selectedFile)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!file || !title) {
      onError?.('Пожалуйста, выберите файл и укажите заголовок')
      return
    }

    // Проверка размера файла (500MB)
    const maxSize = 500 * 1024 * 1024 // 500MB
    if (file.size > maxSize) {
      onError?.(`Файл слишком большой. Максимальный размер: ${maxSize / (1024 * 1024)} MB`)
      return
    }

    setUploading(true)

    try {
      const uploadData: UploadMediaData = {
        file,
        title,
        description: description || undefined,
        tags: tags || undefined,
        category_id: categoryId,
        thumbnail: thumbnail || undefined,
      }

      const video = await videoApi.uploadMedia(uploadData, (progress) => {
        setUploadProgress(progress)
      })
      onSuccess?.(video)
      setUploadProgress(0)

      // Сброс формы
      setFile(null)
      setThumbnail(null)
      setTitle('')
      setDescription('')
      setTags('')
      setCategoryId(undefined)
      setPreview(null)
      const fileInput = document.getElementById('file-input') as HTMLInputElement
      if (fileInput) fileInput.value = ''
    } catch (error: any) {
      const errorMessage = error.message || error.response?.data?.error || 'Ошибка при загрузке файла'
      onError?.(errorMessage)
      console.error('Upload error:', error)
    } finally {
      setUploading(false)
    }
  }

  const getFileType = (file: File): string => {
    if (file.type.startsWith('video/')) return 'video'
    if (file.type === 'image/gif') return 'gif'
    if (file.type.startsWith('image/')) return 'photo'
    return 'unknown'
  }

  return (
    <div className="max-w-2xl mx-auto p-6 bg-white rounded-lg shadow-lg">
      <h2 className="text-2xl font-bold mb-6">Загрузить медиа</h2>

      <form onSubmit={handleSubmit} className="space-y-4">
        {/* Превью файла */}
        {preview && (
          <div className="mb-4">
            {file?.type.startsWith('image/') ? (
              <img
                src={preview}
                alt="Preview"
                className="max-w-full h-auto rounded-lg max-h-64 object-contain"
              />
            ) : file?.type.startsWith('video/') ? (
              <video
                src={preview}
                controls
                className="max-w-full h-auto rounded-lg max-h-64"
              />
            ) : null}
          </div>
        )}

        {/* Выбор файла */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Выберите файл (видео, фото или GIF)
          </label>
          <input
            id="file-input"
            type="file"
            accept="video/*,image/*"
            onChange={handleFileChange}
            className="block w-full text-sm text-gray-500
              file:mr-4 file:py-2 file:px-4
              file:rounded-full file:border-0
              file:text-sm file:font-semibold
              file:bg-blue-50 file:text-blue-700
              hover:file:bg-blue-100"
            required
          />
          {file && (
            <p className="mt-2 text-sm text-gray-600">
              Выбран: {file.name} ({(file.size / 1024 / 1024).toFixed(2)} MB)
            </p>
          )}
        </div>

        {/* Заголовок */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Заголовок *
          </label>
          <input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            placeholder="Введите заголовок"
            required
          />
        </div>

        {/* Описание */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Описание
          </label>
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            placeholder="Введите описание (необязательно)"
            rows={3}
          />
        </div>

        {/* Теги */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Теги
          </label>
          <input
            type="text"
            value={tags}
            onChange={(e) => setTags(e.target.value)}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            placeholder="мем, смешно, прикол (через запятую)"
          />
        </div>

        {/* Превью для видео (опционально) */}
        {file && file.type.startsWith('video/') && (
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Превью (опционально)
            </label>
            <input
              type="file"
              accept="image/*"
              onChange={handleThumbnailChange}
              className="block w-full text-sm text-gray-500
                file:mr-4 file:py-2 file:px-4
                file:rounded-full file:border-0
                file:text-sm file:font-semibold
                file:bg-green-50 file:text-green-700
                hover:file:bg-green-100"
            />
          </div>
        )}

        {/* Индикатор прогресса */}
        {uploading && uploadProgress > 0 && (
          <div className="w-full">
            <div className="flex justify-between text-sm text-gray-600 mb-1">
              <span>Загрузка...</span>
              <span>{uploadProgress}%</span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-2.5">
              <div
                className="bg-blue-600 h-2.5 rounded-full transition-all duration-300"
                style={{ width: `${uploadProgress}%` }}
              />
            </div>
          </div>
        )}

        {/* Кнопка отправки */}
        <button
          type="submit"
          disabled={uploading || !file || !title}
          className="w-full bg-blue-600 text-white py-3 px-4 rounded-lg font-semibold
            hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed
            transition-colors"
        >
          {uploading ? 'Загрузка...' : 'Загрузить'}
        </button>
      </form>
    </div>
  )
}
