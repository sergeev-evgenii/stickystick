'use client'

import { useEffect } from 'react'
import { useSettingsStore } from '@/store/settingsStore'

/** Загружает публичные настройки при монтировании приложения. */
export default function SettingsLoader() {
  const fetchSettings = useSettingsStore((s) => s.fetchSettings)
  useEffect(() => {
    fetchSettings()
  }, [fetchSettings])
  return null
}
