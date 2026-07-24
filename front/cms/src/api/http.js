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
  (response) => {
    const responseData = response.data

    if (!responseData || typeof responseData !== 'object' || !Object.prototype.hasOwnProperty.call(responseData, 'code')) {
      return Promise.reject(new Error('响应格式错误'))
    }

    if (responseData.code !== 0) {
      return Promise.reject(new Error(responseData.message || '请求失败'))
    }

    return responseData.data
  },
  (error) => {
    const responseData = error.response?.data
    const message = typeof responseData === 'string'
      ? responseData
      : responseData?.message || responseData?.error || error.message || '请求失败'

    return Promise.reject(new Error(message))
  }
)

export default http
