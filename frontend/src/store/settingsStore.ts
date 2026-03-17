import { create } from 'zustand'
import { settingsApi } from '@/lib/api/settings'

interface SettingsState {
  showViewCount: boolean
  defaultPublishVK: string
  defaultPublishTelegram: string
  defaultPublishMax: string
  loaded: boolean
  fetchSettings: () => Promise<void>
  setShowViewCount: (value: boolean) => Promise<void>
  setPublishDefaults: (payload: { vk?: string; telegram?: string; max?: string }) => Promise<void>
}

export const useSettingsStore = create<SettingsState>((set, get) => ({
  showViewCount: true,
  defaultPublishVK: '',
  defaultPublishTelegram: '',
  defaultPublishMax: '',
  loaded: false,

  fetchSettings: async () => {
    try {
      const data = await settingsApi.getPublic()
      set({
        showViewCount: data.show_view_count,
        defaultPublishVK: data.default_publish_vk ?? '',
        defaultPublishTelegram: data.default_publish_telegram ?? '',
        defaultPublishMax: data.default_publish_max ?? '',
        loaded: true,
      })
    } catch {
      set({ showViewCount: true, loaded: true })
    }
  },

  setShowViewCount: async (value: boolean) => {
    try {
      const data = await settingsApi.setShowViewCount(value)
      set({ showViewCount: data.show_view_count })
    } catch (e) {
      throw e
    }
  },

  setPublishDefaults: async (payload) => {
    const req: any = {}
    if (payload.vk !== undefined) req.default_publish_vk = payload.vk
    if (payload.telegram !== undefined) req.default_publish_telegram = payload.telegram
    if (payload.max !== undefined) req.default_publish_max = payload.max
    const data = await settingsApi.updatePublishDefaults(req)
    set({
      defaultPublishVK: data.default_publish_vk ?? '',
      defaultPublishTelegram: data.default_publish_telegram ?? '',
      defaultPublishMax: data.default_publish_max ?? '',
    })
  },
}))
