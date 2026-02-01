import api from '../api'

export type MediaType = 'video' | 'photo' | 'gif'

export interface Video {
  id: number
  user_id: number
  category_id?: number
  title: string
  description?: string
  tags?: string
  media_type: MediaType
  media_url: string
  thumbnail_url?: string
  duration: number
  views: number
  moderation_status?: 'pending' | 'approved' | 'rejected'
  created_at: string
  user: {
    id: number
    username: string
    avatar?: string
  }
  category?: {
    id: number
    name: string
    slug: string
  }
  likes?: Array<{
    id: number
    user_id: number
    created_at: string
  }>
  comments?: Array<{
    id: number
    content: string
    created_at: string
    user: {
      id: number
      username: string
      avatar?: string
    }
  }>
}

export interface UploadVideoData {
  title: string
  description?: string
  video_url: string
  thumbnail_url?: string
  duration: number
}

export interface UploadMediaData {
  file: File
  title: string
  description?: string
  tags?: string
  category_id?: number
  thumbnail?: File
}

export const videoApi = {
  getFeed: async (limit = 20, offset = 0, categoryId?: number, tag?: string): Promise<Video[]> => {
    const params: any = { limit, offset }
    if (categoryId) params.category_id = categoryId
    if (tag) params.tag = tag
    
    const response = await api.get('/api/v1/videos', { params })
    return response.data
  },

  getVideo: async (id: number): Promise<Video> => {
    const response = await api.get(`/api/v1/videos/${id}`)
    return response.data
  },

  uploadVideo: async (data: UploadVideoData): Promise<Video> => {
    const response = await api.post('/api/v1/videos', data)
    return response.data
  },

  uploadMedia: async (data: UploadMediaData, onProgress?: (progress: number) => void): Promise<Video> => {
    const formData = new FormData()
    formData.append('file', data.file)
    formData.append('title', data.title)
    if (data.description) {
      formData.append('description', data.description)
    }
    if (data.tags) {
      formData.append('tags', data.tags)
    }
    if (data.category_id) {
      formData.append('category_id', data.category_id.toString())
    }
    if (data.thumbnail) {
      formData.append('thumbnail', data.thumbnail)
    }

    const response = await api.post('/api/v1/videos/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      timeout: 300000, // 5 минут
      onUploadProgress: (progressEvent) => {
        if (progressEvent.total && onProgress) {
          const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total)
          onProgress(percentCompleted)
        }
      },
    })
    return response.data
  },

  deleteVideo: async (id: number): Promise<void> => {
    await api.delete(`/api/v1/videos/${id}`)
  },

  likeVideo: async (id: number): Promise<void> => {
    await api.post(`/api/v1/videos/${id}/like`)
  },

  unlikeVideo: async (id: number): Promise<void> => {
    await api.delete(`/api/v1/videos/${id}/like`)
  },

  addComment: async (id: number, content: string): Promise<any> => {
    const response = await api.post(`/api/v1/videos/${id}/comment`, { content })
    return response.data
  },

  // Модерация (только для админов)
  getPendingModeration: async (limit = 20, offset = 0): Promise<Video[]> => {
    const response = await api.get('/api/v1/videos/moderation/pending', {
      params: { limit, offset },
    })
    return response.data
  },

  moderateVideo: async (id: number, status: 'approved' | 'rejected'): Promise<void> => {
    await api.post(`/api/v1/videos/${id}/moderate`, { status })
  },
}
