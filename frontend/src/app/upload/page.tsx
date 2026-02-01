'use client'

import { useRouter } from 'next/navigation'
import { useEffect } from 'react'
import MediaUpload from '@/components/MediaUpload'
import { useState } from 'react'
import { useAuthStore } from '@/store/authStore'
import Link from 'next/link'

export default function UploadPage() {
  const router = useRouter()
  const { isAuthenticated } = useAuthStore()
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null)

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login')
    }
  }, [isAuthenticated, router])

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-center">
          <p className="text-lg mb-4">Необходима авторизация</p>
          <Link
            href="/login"
            className="bg-blue-600 text-white px-6 py-3 rounded-lg font-semibold hover:bg-blue-700 transition-colors"
          >
            Войти
          </Link>
        </div>
      </div>
    )
  }

  const handleSuccess = (video: any) => {
    setMessage({ type: 'success', text: 'Медиа успешно загружено!' })
    setTimeout(() => {
      router.push('/')
    }, 2000)
  }

  const handleError = (error: string) => {
    setMessage({ type: 'error', text: error })
  }

  return (
    <div className="min-h-screen bg-gray-100 py-8">
      <div className="container mx-auto px-4">
        {message && (
          <div
            className={`mb-4 p-4 rounded-lg ${
              message.type === 'success'
                ? 'bg-green-100 text-green-800'
                : 'bg-red-100 text-red-800'
            }`}
          >
            {message.text}
          </div>
        )}
        <MediaUpload onSuccess={handleSuccess} onError={handleError} />
      </div>
    </div>
  )
}
