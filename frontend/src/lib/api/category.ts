import api from '../api'

export interface Category {
  id: number
  name: string
  slug: string
  created_at: string
}

export const categoryApi = {
  getAll: async (): Promise<Category[]> => {
    const response = await api.get('/api/v1/categories')
    return response.data
  },

  getById: async (id: number): Promise<Category> => {
    const response = await api.get(`/api/v1/categories/${id}`)
    return response.data
  },

  create: async (name: string, slug?: string): Promise<Category> => {
    const response = await api.post('/api/v1/categories', { name, slug })
    return response.data
  },

  update: async (id: number, name?: string, slug?: string): Promise<Category> => {
    const response = await api.put(`/api/v1/categories/${id}`, { name, slug })
    return response.data
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/api/v1/categories/${id}`)
  },
}
