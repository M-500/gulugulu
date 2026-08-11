<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Button as VanButton } from 'vant'

import LoginDialog from '@/components/auth/LoginDialog.vue'
import AppSidebar from '@/components/layout/AppSidebar.vue'
import MobileTabBar from '@/components/layout/MobileTabBar.vue'
import SearchHeader from '@/components/layout/SearchHeader.vue'
import FeedState from '@/components/note/FeedState.vue'
import ImageWorkDialog from '@/components/note/ImageWorkDialog.vue'
import WaterfallFeed from '@/components/note/WaterfallFeed.vue'
import { getAppUserProfile, getUserPublishedWorks } from '@/services/profileService'
import { useAuthStore } from '@/stores/auth'
import { withProfileMock } from './profileMock'
import ProfileEmptyState from './components/ProfileEmptyState.vue'
import ProfileHeader from './components/ProfileHeader.vue'
import ProfileTabs from './components/ProfileTabs.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const activeTab = ref('notes')
const keyword = ref('')
const showLoginDialog = ref(false)
const selectedImageNote = ref(null)
const showImageDialog = ref(false)
const profileUser = ref(null)
const profileLoading = ref(false)
const profileError = ref('')
const profileNotes = ref([])
const worksLoading = ref(false)
const worksError = ref('')
const worksPage = ref(1)
const worksPageSize = ref(20)
const worksTotal = ref(0)
const worksHasMore = ref(false)
let requestVersion = 0

const isOwnProfile = computed(() => {
  const currentUserId = Number(authStore.user?.userId)
  const profileUserId = Number(profileUser.value?.userId)
  return currentUserId > 0 && profileUserId > 0 && currentUserId === profileUserId
})

const visibleNotes = computed(() => {
  const value = keyword.value.trim().toLowerCase()
  if (!value) return profileNotes.value
  return profileNotes.value.filter((note) => {
    return [note.title, note.author, note.quote].some((text) => String(text || '').toLowerCase().includes(value))
  })
})

watch(
  () => route.params.userId,
  (userId) => loadProfilePage(userId),
  { immediate: true }
)

function openImageWork(note) {
  selectedImageNote.value = note
  showImageDialog.value = true
}

function openAuthorProfile(note) {
  router.push({ name: 'user-profile', params: { userId: note.authorId || 'mock' } })
}

function handleLikeChange({ note, result }) {
	if (!note) return
  note.liked = result.liked
  note.likes = String(result.count)
}

async function resolveUserId(rawUserId) {
  if (rawUserId !== 'me') return Number(rawUserId)
  if (!authStore.user?.userId && authStore.isLoggedIn) {
    await authStore.fetchUserInfo()
  }
  return Number(authStore.user?.userId)
}

async function loadProfilePage(rawUserId) {
  const currentVersion = ++requestVersion
  profileLoading.value = true
  worksLoading.value = true
  profileError.value = ''
  worksError.value = ''
  profileUser.value = null
  profileNotes.value = []
  worksTotal.value = 0
  worksHasMore.value = false

  try {
	if (authStore.isLoggedIn && !authStore.user?.userId) {
	  await authStore.fetchUserInfo().catch(() => {})
	}
    const userId = await resolveUserId(rawUserId)
    if (!Number.isInteger(userId) || userId <= 0) {
      throw new Error('用户ID不正确')
    }
    const [userResult, worksResult] = await Promise.allSettled([
      getAppUserProfile(userId),
	  getUserPublishedWorks(userId, { page: 1, pageSize: worksPageSize.value, token: authStore.accessToken })
    ])
    if (currentVersion !== requestVersion) return

    if (userResult.status === 'rejected') {
      throw userResult.reason
    }
    profileUser.value = withProfileMock(userResult.value)

    if (worksResult.status === 'fulfilled') {
      const works = worksResult.value
      profileNotes.value = works.list
      worksTotal.value = works.total
      worksPage.value = works.page
      worksPageSize.value = works.pageSize
      worksHasMore.value = works.hasMore
    } else {
      worksError.value = worksResult.reason?.message || '作品列表加载失败'
    }
  } catch (error) {
    if (currentVersion === requestVersion) {
      profileError.value = error.message || '用户主页加载失败'
    }
  } finally {
    if (currentVersion === requestVersion) {
      profileLoading.value = false
      worksLoading.value = false
    }
  }
}

async function loadMoreWorks() {
  if (!profileUser.value || worksLoading.value || !worksHasMore.value) return

  const currentVersion = requestVersion
  const userId = profileUser.value.userId
  worksLoading.value = true
  worksError.value = ''
  try {
    const works = await getUserPublishedWorks(userId, {
      page: worksPage.value + 1,
	  pageSize: worksPageSize.value,
	  token: authStore.accessToken
    })
    if (currentVersion !== requestVersion) return

    profileNotes.value = [...profileNotes.value, ...works.list]
    worksPage.value = works.page
    worksHasMore.value = works.hasMore
    worksTotal.value = works.total
  } catch (error) {
    if (currentVersion === requestVersion) {
      worksError.value = error.message || '作品列表加载失败'
    }
  } finally {
    if (currentVersion === requestVersion) {
      worksLoading.value = false
    }
  }
}
</script>

<template>
  <div class="profile-shell">
    <AppSidebar @login-request="showLoginDialog = true" />

    <main class="profile-main">
      <SearchHeader v-model="keyword" />
      <div class="profile-content">
        <FeedState
          v-if="profileLoading && !profileUser"
          type="loading"
          message="正在加载用户主页"
        />
        <FeedState
          v-else-if="profileError"
          type="error"
          :message="profileError"
          @retry="loadProfilePage(route.params.userId)"
        />
        <template v-else-if="profileUser">
          <ProfileHeader
            :user="profileUser"
            :is-own-profile="isOwnProfile"
          />
          <ProfileTabs
            v-model="activeTab"
            :work-count="worksTotal"
          />

          <template v-if="activeTab === 'notes'">
            <FeedState
              v-if="!visibleNotes.length && worksLoading"
              type="loading"
              message="正在加载已发布作品"
            />
            <FeedState
              v-else-if="worksError && !profileNotes.length"
              type="error"
              :message="worksError"
              @retry="loadProfilePage(route.params.userId)"
            />
            <FeedState
              v-else-if="!visibleNotes.length"
              type="empty"
              :message="keyword.trim() ? '没有找到符合条件的作品' : '该用户暂时没有已发布作品'"
            />
            <template v-else>
              <WaterfallFeed
                class="profile-feed"
                :notes="visibleNotes"
                @open-author="openAuthorProfile"
                @open-note="openImageWork"
                @login-request="showLoginDialog = true"
                @like-change="handleLikeChange"
              />
              <div
                v-if="worksHasMore || worksLoading || worksError"
                class="profile-more"
              >
                <p v-if="worksError">
                  {{ worksError }}
                </p>
                <VanButton
                  color="#2e3036"
                  round
                  :loading="worksLoading"
                  @click="loadMoreWorks"
                >
                  {{ worksError ? '重新加载' : '加载更多' }}
                </VanButton>
              </div>
            </template>
          </template>
          <ProfileEmptyState v-else />
        </template>
      </div>
    </main>

    <MobileTabBar @login-request="showLoginDialog = true" />
    <LoginDialog
      v-model:show="showLoginDialog"
      @success="loadProfilePage(route.params.userId)"
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
.profile-shell {
  display: flex;
  min-height: 100vh;
  background: #fff;
}

.profile-main {
  min-width: 0;
  flex: 1;
}

.profile-content {
  width: min(100%, 1500px);
  margin: 0 auto;
  padding: 0 36px 80px;
}

.profile-feed {
  margin-top: 8px;
}

.profile-more {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 10px 0 32px;
}

.profile-more p {
  margin: 0;
  color: #9a5b64;
  font-size: 13px;
}

.profile-feed :deep(.waterfall-feed) {
  column-count: 5;
}

@media (max-width: 900px) {
  .profile-shell {
    display: block;
  }

  .profile-content {
    padding: 0 12px 86px;
  }
}
</style>
