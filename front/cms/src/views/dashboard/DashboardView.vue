<template>
  <main class="dashboard">
    <section class="dashboard__content">
      <AppLogo />
      <h1>内容管理后台骨架已就绪</h1>
      <p>后续可以在这里接入笔记、视频、长文、审核和用户管理模块。</p>
      <button type="button" @click="logout">退出登录</button>
    </section>
  </main>
</template>

<script setup>
import { useRouter } from 'vue-router'

import AppLogo from '@/components/AppLogo.vue'
import { useAuthStore } from '@/stores/auth'

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

.dashboard__content {
  display: grid;
  width: min(620px, 100%);
  gap: 18px;
  border: 1px solid rgba(17, 24, 39, 0.08);
  border-radius: 8px;
  background: #fff;
  box-shadow: 0 18px 50px rgba(15, 23, 42, 0.08);
  padding: 36px;
}

.dashboard h1 {
  margin: 10px 0 0;
  color: var(--color-text);
  font-size: 26px;
}

.dashboard p {
  margin: 0;
  color: var(--color-text-muted);
  line-height: 1.8;
}

.dashboard button {
  width: 120px;
  height: 42px;
  border: 0;
  border-radius: 6px;
  background: var(--color-primary);
  color: #fff;
  cursor: pointer;
  font-weight: 700;
}
</style>
