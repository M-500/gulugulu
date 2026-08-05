import { defineStore } from 'pinia'

import { getHomeFeed } from '@/services/feedService'
import { useAuthStore } from '@/stores/auth'

export const useFeedStore = defineStore('feed', {
  state: () => ({
    items: [],
    loading: false,
    error: '',
    page: 1,
    pageSize: 20,
    hasMore: false
  }),
  actions: {
    async fetchHomeFeed({ reset = true } = {}) {
	  const authStore = useAuthStore()
      this.loading = true
      this.error = ''
      try {
        const nextPage = reset ? 1 : this.page + 1
        const data = await getHomeFeed({
          page: nextPage,
          pageSize: this.pageSize
		  , userId: authStore.user?.userId || 0
        })
        this.items = reset ? data.list : [...this.items, ...data.list]
        this.page = data.page
        this.pageSize = data.pageSize
        this.hasMore = data.hasMore
      } catch (error) {
        this.error = error.message || '推荐作品加载失败'
      } finally {
        this.loading = false
      }
    }
  }
})
