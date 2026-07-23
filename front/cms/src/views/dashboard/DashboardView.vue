<template>
  <div class="dashboard-shell">
    <DashboardHeader
      :menu-open="menuOpen"
      @toggle-menu="menuOpen = !menuOpen"
      @logout="logout"
    />

    <div class="dashboard-shell__body">
      <DashboardSidebar
        :open="menuOpen"
        :active-item="activeMenu"
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
        <PublishView v-else-if="activeMenu === 'publish'" />
        <DashboardPlaceholder v-else :menu="currentMenu" />
      </main>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

import { useAuthStore } from '@/stores/auth'

import DashboardHeader from './components/DashboardHeader.vue'
import DashboardOverview from './components/DashboardOverview.vue'
import DashboardPlaceholder from './components/DashboardPlaceholder.vue'
import DashboardSidebar from './components/DashboardSidebar.vue'
import PublishView from '../publish/PublishView.vue'
import { dashboardMenus } from './dashboardMenus'

const router = useRouter()
const authStore = useAuthStore()

const activeMenu = ref('home')
const menuOpen = ref(false)
const currentMenu = computed(() => dashboardMenus.find((item) => item.key === activeMenu.value))

function handleMenuSelect(key) {
  activeMenu.value = key
  menuOpen.value = false
}

function logout() {
  authStore.clearSession()
  router.push({ name: 'Auth' })
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
