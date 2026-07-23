<template>
  <main class="dashboard">
    <DashboardWelcome @logout="logout" />
  </main>
</template>

<script setup>
import { useRouter } from 'vue-router'

import { useAuthStore } from '@/stores/auth'

import DashboardWelcome from './components/DashboardWelcome.vue'

const router = useRouter()
const authStore = useAuthStore()

function logout() {
  authStore.clearSession()
  router.push({ name: 'Auth' })
}
</script>

<style scoped>
.dashboard {
  display: grid;
  min-height: 100vh;
  place-items: center;
  background: #f6f7fb;
  padding: 24px;
}
</style>
