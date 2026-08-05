<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Button as VanButton } from 'vant'

import LoginDialog from '@/components/auth/LoginDialog.vue'
import FloatingActions from '@/components/layout/FloatingActions.vue'
import AppSidebar from '@/components/layout/AppSidebar.vue'
import MobileTabBar from '@/components/layout/MobileTabBar.vue'
import SearchHeader from '@/components/layout/SearchHeader.vue'
import FeedState from '@/components/note/FeedState.vue'
import ImageWorkDialog from '@/components/note/ImageWorkDialog.vue'
import RecommendTabs from '@/components/note/RecommendTabs.vue'
import WaterfallFeed from '@/components/note/WaterfallFeed.vue'
import { useAuthStore } from '@/stores/auth'
import { useFeedStore } from '@/stores/feed'

const keyword = ref('')
const router = useRouter()
const activeChannel = ref('recommend')
const showLoginDialog = ref(false)
const selectedImageNote = ref(null)
const showImageDialog = ref(false)
const feedStore = useFeedStore()
const authStore = useAuthStore()
const recommendTabs = [
  { key: 'recommend', label: '推荐' },
  { key: 'worldcup', label: '世界杯' },
  { key: 'fashion', label: '穿搭' },
  { key: 'food', label: '美食' },
  { key: 'beauty', label: '彩妆' },
  { key: 'film', label: '影视' },
  { key: 'career', label: '职场' },
  { key: 'emotion', label: '情感' },
  { key: 'home', label: '家居' },
  { key: 'game', label: '游戏' },
  { key: 'travel', label: '旅行' },
  { key: 'fitness', label: '健身' },
  { key: 'video', label: '视频' }
]

const notes = computed(() => {
  const value = keyword.value.trim().toLowerCase()
  if (!value) return feedStore.items

  return feedStore.items.filter((note) => {
    return [note.title, note.author, note.quote].some((text) => String(text || '').toLowerCase().includes(value))
  })
})

onMounted(async () => {
  if (authStore.isLoggedIn) {
	await authStore.fetchUserInfo().catch(() => {})
  } else {
    window.setTimeout(() => {
      showLoginDialog.value = true
    }, 500)
  }
	feedStore.fetchHomeFeed()
})

function handleChannelChange() {
  feedStore.fetchHomeFeed()
}

function openImageWork(note) {
  selectedImageNote.value = note
  showImageDialog.value = true
}

function openAuthorProfile(note) {
  router.push({ name: 'user-profile', params: { userId: note.authorId || note.id || 'mock' } })
}

function handleLikeChange({ note, result }) {
	if (!note) return
  note.liked = result.liked
  note.likes = String(result.count)
}
</script>

<template>
  <div class="home-shell">
    <AppSidebar @login-request="showLoginDialog = true" />

    <main class="home-main">
      <SearchHeader v-model="keyword" />
      <div class="home-content">
        <RecommendTabs
          v-model="activeChannel"
          :tabs="recommendTabs"
          @change="handleChannelChange"
        />
        <FeedState
          v-if="feedStore.loading && !feedStore.items.length"
          type="loading"
          message="正在加载推荐作品"
        />
        <FeedState
          v-else-if="feedStore.error"
          type="error"
          :message="feedStore.error"
          @retry="feedStore.fetchHomeFeed()"
        />
        <FeedState
          v-else-if="!notes.length"
          type="empty"
          message="暂时没有符合条件的作品"
        />
        <template v-else>
          <WaterfallFeed
            :notes="notes"
            @open-author="openAuthorProfile"
            @open-note="openImageWork"
            @login-request="showLoginDialog = true"
            @like-change="handleLikeChange"
          />
          <div
            v-if="feedStore.hasMore || feedStore.loading"
            class="home-more"
          >
            <VanButton
              color="#2e3036"
              round
              :disabled="feedStore.loading"
              @click="feedStore.fetchHomeFeed({ reset: false })"
            >
              {{ feedStore.loading ? '加载中' : '加载更多' }}
            </VanButton>
          </div>
        </template>
      </div>
    </main>

    <FloatingActions />
    <MobileTabBar @login-request="showLoginDialog = true" />
    <LoginDialog
      v-model:show="showLoginDialog"
      @success="feedStore.fetchHomeFeed()"
    />
    <ImageWorkDialog
      v-model:show="showImageDialog"
      :note="selectedImageNote"
      @login-request="showLoginDialog = true"
      @open-author="openAuthorProfile"
      @like-change="handleLikeChange"
    />
  </div>
</template>

<style scoped>
.home-shell {
  display: flex;
  min-height: 100vh;
  background: #fff;
}

.home-main {
  min-width: 0;
  flex: 1;
}

.home-content {
  width: min(100%, var(--layout-max));
  margin: 0 auto;
  padding: 0 36px 64px;
}

.home-more {
  display: flex;
  justify-content: center;
  padding: 10px 0 32px;
}

@media (max-width: 900px) {
  .home-shell {
    display: block;
  }

  .home-content {
    padding: 0 12px 86px;
  }
}
</style>
