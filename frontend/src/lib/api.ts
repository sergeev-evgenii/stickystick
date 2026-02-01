import axios from 'axios'

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || '',
  timeout: 300000, // 5 минут для больших файлов
})

// Request interceptor для добавления токена
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Response interceptor для обработки ошибок
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.code === 'ECONNABORTED' || error.message.includes('timeout')) {
      return Promise.reject(new Error('Превышено время ожидания. Файл слишком большой или соединение медленное.'))
    }
    if (error.code === 'ERR_CONNECTION_RESET' || error.message.includes('ERR_CONNECTION_RESET')) {
      return Promise.reject(new Error('Соединение разорвано. Возможно, файл слишком большой или сервер не может его обработать.'))
    }
    return Promise.reject(error)
  }
)

export default api
