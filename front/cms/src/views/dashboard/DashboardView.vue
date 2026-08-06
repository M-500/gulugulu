<template>
  <div class="dashboard-shell">
    <DashboardHeader
      :menu-open="menuOpen"
      :profile="authStore.profile"
      @profile="openProfile"
      @toggle-menu="menuOpen = !menuOpen"
      @logout="logout"
    />

    <div class="dashboard-shell__body">
      <DashboardSidebar
        :open="menuOpen"
        :active-item="activeMenu"
        :review-count="reviewCount"
        @select="handleMenuSelect"
      />

      <button
        v-if="menuOpen"
        class="dashboard-shell__overlay"
        type="button"
        aria-label="关闭菜单"
        @click="menuOpen = false"
      />

      <main class="dashboard-shell__content">
        <DashboardOverview v-if="activeMenu === 'home'" />
        <PublishView v-else-if="activeMenu === 'publish'" :edit-work-id="editingWorkId" @close-edit="closeWorkEdit" @edit-saved="closeWorkEdit" />
        <WorksManagementView v-else-if="activeMenu === 'works'" @edit-work="openWorkEdit" />
        <ReviewCenterView v-else-if="activeMenu === 'review'" @count-change="reviewCount = $event" />
        <ProfileView v-else-if="activeMenu === 'profile'" />
        <DashboardPlaceholder v-else :menu="currentMenu" />
      </main>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { useAuthStore } from '@/stores/auth'
import { getCurrentUser } from '@/api/profile'

import DashboardHeader from './components/DashboardHeader.vue'
import DashboardOverview from './components/DashboardOverview.vue'
import DashboardPlaceholder from './components/DashboardPlaceholder.vue'
import DashboardSidebar from './components/DashboardSidebar.vue'
import PublishView from '../publish/PublishView.vue'
import WorksManagementView from '../works/WorksManagementView.vue'
import ReviewCenterView from '../review/ReviewCenterView.vue'
import ProfileView from '../profile/ProfileView.vue'
import { dashboardMenus } from './dashboardMenus'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const availableMenuKeys = new Set([...dashboardMenus.map((item) => item.key), 'profile'])
const activeMenu = ref(resolveMenu(route.query.section))
const menuOpen = ref(false)
const reviewCount = ref(0)
const editingWorkId = ref(resolveWorkId(route.query.edit))
const currentMenu = computed(() => dashboardMenus.find((item) => item.key === activeMenu.value))

onMounted(loadProfile)

// URL 是当前栏目状态的唯一持久化来源，刷新和浏览器前进后退都会恢复原页面。
watch(
  () => [route.query.section, route.query.edit],
  ([section, edit]) => {
    activeMenu.value = resolveMenu(section)
    editingWorkId.value = activeMenu.value === 'publish' ? resolveWorkId(edit) : null
  }
)

async function loadProfile() {
  try {
    authStore.setProfile(await getCurrentUser())
  } catch {
    // 页面主体接口会处理登录失效，这里不阻断工作台加载。
  }
}

function handleMenuSelect(key) {
  menuOpen.value = false
  navigateToSection(key)
}

function openWorkEdit(workId) {
  navigateToSection('publish', workId)
}

function closeWorkEdit() {
  navigateToSection('works')
}

function logout() {
  authStore.clearSession()
  router.push({ name: 'Auth' })
}

function openProfile() {
  menuOpen.value = false
  navigateToSection('profile')
}

function navigateToSection(section, editWorkId = null) {
  const query = section === 'home' ? {} : { section }
  if (section === 'publish' && editWorkId) query.edit = String(editWorkId)
  router.push({ name: 'Dashboard', query })
}

function resolveMenu(section) {
  const key = String(section || '')
  return availableMenuKeys.has(key) ? key : 'home'
}

function resolveWorkId(value) {
  const workId = Number(value)
  return Number.isInteger(workId) && workId > 0 ? workId : null
}
</script>

<style scoped>
.dashboard-shell {
  min-height: 100vh;
  background: #f5f6f8;
}

.dashboard-shell__body {
  display: flex;
  min-height: calc(100vh - 68px);
}

.dashboard-shell__content {
  flex: 1;
  min-width: 0;
  padding: 28px 32px 44px;
}

.dashboard-shell__overlay {
  display: none;
}

@media (max-width: 900px) {
  .dashboard-shell__content {
    padding: 24px 20px 40px;
  }

  .dashboard-shell__overlay {
    position: fixed;
    z-index: 20;
    inset: 68px 0 0;
    display: block;
    border: 0;
    background: rgba(15, 23, 42, 0.32);
  }
}

@media (max-width: 560px) {
  .dashboard-shell__content {
    padding: 20px 16px 32px;
  }
}
</style>
