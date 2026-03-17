import api from '../api'

export interface PublicSettings {
  show_view_count: boolean
  default_publish_vk: string
  default_publish_telegram: string
  default_publish_max: string
}

export const settingsApi = {
  getPublic: async (): Promise<PublicSettings> => {
    const { data } = await api.get<PublicSettings>('/api/v1/settings')
    return data
  },

  setShowViewCount: async (showViewCount: boolean): Promise<PublicSettings> => {
    const { data } = await api.patch<PublicSettings>('/api/v1/settings', { show_view_count: showViewCount })
    return data
  },

  updatePublishDefaults: async (payload: Partial<Pick<PublicSettings, 'default_publish_vk' | 'default_publish_telegram' | 'default_publish_max'>>): Promise<PublicSettings> => {
    const { data } = await api.patch<PublicSettings>('/api/v1/settings', payload)
    return data
  },
}
