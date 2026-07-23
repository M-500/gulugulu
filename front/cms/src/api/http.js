import axios from 'axios'

import { useAuthStore } from '@/stores/auth'

const http = axios.create({
  baseURL: process.env.VUE_APP_API_BASE_URL || '',
  timeout: 10000
})

http.interceptors.request.use((config) => {
  const authStore = useAuthStore()

  if (authStore.accessToken) {
    config.headers.Authorization = `Bearer ${authStore.accessToken}`
  }

  return config
})

http.interceptors.response.use(
  (response) => response.data,
  (error) => {
    const responseData = error.response?.data
    const message = typeof responseData === 'string'
      ? responseData
      : responseData?.message || responseData?.error || error.message || '请求失败'

    return Promise.reject(new Error(message))
  }
)

export default http
