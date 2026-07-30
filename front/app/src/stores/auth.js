import { defineStore } from 'pinia'

import { getUserInfo, loginWithEmail } from '@/services/authService'

const STORAGE_KEY = 'gulugulu-app-auth'

function readStoredAuth() {
  try {
    return JSON.parse(localStorage.getItem(STORAGE_KEY) || '{}')
  } catch {
    return {}
  }
}

export const useAuthStore = defineStore('app-auth', {
  state: () => {
    const stored = readStoredAuth()
    return {
      accessToken: stored.accessToken || '',
      user: stored.user || null
    }
  },
  getters: {
    isLoggedIn: (state) => Boolean(state.accessToken)
  },
  actions: {
    persist() {
      localStorage.setItem(STORAGE_KEY, JSON.stringify({
        accessToken: this.accessToken,
        user: this.user
      }))
    },
    setUser(profile = {}) {
      this.user = {
        userId: profile.userId || this.user?.userId || 0,
        nickName: profile.nickName || '我',
        avatarUrl: profile.avatarUrl || ''
      }
    },
    async login(payload) {
      const data = await loginWithEmail(payload)
      this.accessToken = data.accessToken || ''
      this.setUser(data)
      await this.fetchUserInfo().catch(() => {
        this.persist()
      })
      return data
    },
    async fetchUserInfo() {
      if (!this.accessToken) {
        return null
      }
      const profile = await getUserInfo(this.accessToken)
      this.setUser(profile)
      this.persist()
      return this.user
    },
    logout() {
      this.accessToken = ''
      this.user = null
      localStorage.removeItem(STORAGE_KEY)
    }
  }
})
