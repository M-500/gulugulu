import { defineStore } from 'pinia'

const TOKEN_KEY = 'gulugulu_cms_access_token'
const TOKEN_META_KEY = 'gulugulu_cms_token_meta'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    accessToken: localStorage.getItem(TOKEN_KEY) || '',
    tokenMeta: JSON.parse(localStorage.getItem(TOKEN_META_KEY) || '{}')
  }),
  getters: {
    isLoggedIn: (state) => Boolean(state.accessToken)
  },
  actions: {
    setSession(payload) {
      this.accessToken = payload.accessToken || ''
      this.tokenMeta = {
        accessExpire: payload.accessExpire,
        refreshAfter: payload.refreshAfter
      }

      localStorage.setItem(TOKEN_KEY, this.accessToken)
      localStorage.setItem(TOKEN_META_KEY, JSON.stringify(this.tokenMeta))
    },
    clearSession() {
      this.accessToken = ''
      this.tokenMeta = {}
      localStorage.removeItem(TOKEN_KEY)
      localStorage.removeItem(TOKEN_META_KEY)
    }
  }
})
