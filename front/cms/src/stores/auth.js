import { defineStore } from 'pinia'

import {
  login as loginRequest,
  register as registerRequest
} from '@/api/auth'

const TOKEN_KEY = 'gulugulu_cms_access_token'
const TOKEN_META_KEY = 'gulugulu_cms_token_meta'
const PROFILE_KEY = 'gulugulu_cms_profile'

function getStoredTokenMeta() {
  try {
    return JSON.parse(localStorage.getItem(TOKEN_META_KEY) || '{}')
  } catch {
    localStorage.removeItem(TOKEN_META_KEY)
    return {}
  }
}

function getStoredProfile() {
  try {
    return JSON.parse(localStorage.getItem(PROFILE_KEY) || '{}')
  } catch {
    localStorage.removeItem(PROFILE_KEY)
    return {}
  }
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    accessToken: localStorage.getItem(TOKEN_KEY) || '',
    tokenMeta: getStoredTokenMeta(),
    profile: getStoredProfile()
  }),
  getters: {
    isLoggedIn: (state) => {
      if (!state.accessToken) {
        return false
      }

      if (!state.tokenMeta.accessExpire) {
        return true
      }

      return state.tokenMeta.accessExpire > Math.floor(Date.now() / 1000)
    }
  },
  actions: {
    setSession(payload) {
      if (!payload?.accessToken) {
        throw new Error('接口未返回有效的登录凭证')
      }

      this.accessToken = payload.accessToken
      this.tokenMeta = {
        accessExpire: payload.accessExpire,
        refreshAfter: payload.refreshAfter
      }

      localStorage.setItem(TOKEN_KEY, this.accessToken)
      localStorage.setItem(TOKEN_META_KEY, JSON.stringify(this.tokenMeta))

      if (payload.userId || payload.nickName || payload.avatarUrl) {
        this.setProfile({
          userId: payload.userId || 0,
          nickName: payload.nickName || '咕噜创作者',
          avatarUrl: payload.avatarUrl || ''
        })
      }
    },
    setProfile(profile) {
      this.profile = profile || {}
      localStorage.setItem(PROFILE_KEY, JSON.stringify(this.profile))
    },
    async login(form) {
      const session = await loginRequest({
        email: form.email,
        password: form.password,
        captchaId: form.captchaId,
        captchaCode: form.captchaCode
      })

      this.setSession(session)
      return session
    },
    async register(form) {
      const session = await registerRequest({
        email: form.email,
        nickName: form.nickName,
        password: form.password,
        verificationCode: form.verificationCode
      })

      this.setSession(session)
      return session
    },
    clearSession() {
      this.accessToken = ''
      this.tokenMeta = {}
      this.profile = {}
      localStorage.removeItem(TOKEN_KEY)
      localStorage.removeItem(TOKEN_META_KEY)
      localStorage.removeItem(PROFILE_KEY)
    }
  }
})
