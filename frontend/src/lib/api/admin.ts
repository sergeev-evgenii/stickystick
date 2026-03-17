import api from '../api'

export interface AnalyticsStats {
  unique_visitors: number
  total_video_views: number
  total_logins: number
  total_registrations: number
  total_likes: number
  total_uploads: number
  total_generate_video_clicks?: number
  videos_total_views?: number
}

export interface ActivityLogItem {
  id: number
  ip: string
  user_id: number | null
  action: string
  resource_type: string
  resource_id: number
  user_agent: string
  created_at: string
  user?: {
    id: number
    username: string
    email?: string
  }
}

export interface AnalyticsResponse {
  since: string
  period: string
  stats: AnalyticsStats
  activity: ActivityLogItem[]
}

export const adminApi = {
  getAnalytics: async (
    since: '24h' | '7d' | '30d' = '24h',
    limit = 100,
    offset = 0
  ): Promise<AnalyticsResponse> => {
    const response = await api.get('/api/v1/admin/analytics', {
      params: { since, limit, offset },
    })
    return response.data
  },
}
