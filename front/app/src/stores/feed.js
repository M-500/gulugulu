import { defineStore } from 'pinia'

import { getHomeFeed } from '@/services/feedService'

export const useFeedStore = defineStore('feed', {
  state: () => ({
    items: [],
    loading: false
  }),
  actions: {
    async fetchHomeFeed() {
      this.loading = true
      try {
        this.items = await getHomeFeed()
      } finally {
        this.loading = false
      }
    }
  }
})
