import api from '../api'

export const analyticsApi = {
  logGenerateVideoClick: (): Promise<void> =>
    api.post('/api/v1/analytics/generate-video-click').then(() => {}),
}
