import api from '../api'

export interface RegisterData {
  username: string
  email: string
  password: string
}

export interface LoginData {
  email: string
  password: string
}

export interface AuthResponse {
  user: {
    id: number
    username: string
    email: string
    avatar?: string
    bio?: string
    is_admin?: boolean
  }
  token: string
}

export const authApi = {
  register: async (data: RegisterData): Promise<AuthResponse> => {
    const response = await api.post('/api/v1/auth/register', data)
    return response.data
  },

  login: async (data: LoginData): Promise<AuthResponse> => {
    const response = await api.post('/api/v1/auth/login', data)
    return response.data
  },
}
