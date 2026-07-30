<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

import LoginDialog from '@/components/auth/LoginDialog.vue'
import AppSidebar from '@/components/layout/AppSidebar.vue'
import MobileTabBar from '@/components/layout/MobileTabBar.vue'
import SearchHeader from '@/components/layout/SearchHeader.vue'
import ImageWorkDialog from '@/components/note/ImageWorkDialog.vue'
import WaterfallFeed from '@/components/note/WaterfallFeed.vue'
import { profileNotes, profileUser } from './profileMock'
import ProfileEmptyState from './components/ProfileEmptyState.vue'
import ProfileHeader from './components/ProfileHeader.vue'
import ProfileTabs from './components/ProfileTabs.vue'

const router = useRouter()
const activeTab = ref('notes')
const keyword = ref('')
const showLoginDialog = ref(false)
const selectedImageNote = ref(null)
const showImageDialog = ref(false)

function openImageWork(note) {
  selectedImageNote.value = note
  showImageDialog.value = true
}

function openAuthorProfile(note) {
  router.push({ name: 'user-profile', params: { userId: note.authorId || 'mock' } })
}
</script>

<template>
  <div class="profile-shell">
    <AppSidebar @login-request="showLoginDialog = true" />

    <main class="profile-main">
      <SearchHeader v-model="keyword" />
      <div class="profile-content">
        <ProfileHeader :user="profileUser" />
        <ProfileTabs v-model="activeTab" />

        <WaterfallFeed
          v-if="activeTab === 'notes'"
          class="profile-feed"
          :notes="profileNotes"
          @open-author="openAuthorProfile"
          @open-note="openImageWork"
        />
        <ProfileEmptyState v-else />
      </div>
    </main>

    <MobileTabBar @login-request="showLoginDialog = true" />
    <LoginDialog v-model:show="showLoginDialog" />
    <ImageWorkDialog
      v-model:show="showImageDialog"
      :note="selectedImageNote"
      @login-request="showLoginDialog = true"
      @open-author="openAuthorProfile"
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
