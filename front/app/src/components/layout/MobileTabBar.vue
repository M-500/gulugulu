<script setup>
import { computed } from 'vue'

import BaseIcon from '@/components/common/BaseIcon.vue'
import { primaryNavItems } from '@/constants/navigation'
import { useAuthStore } from '@/stores/auth'

const emit = defineEmits(['login-request'])
const authStore = useAuthStore()
const items = computed(() => primaryNavItems.filter((item) => item.key !== 'profile').slice(0, 7))
</script>

<template>
  <nav
    class="mobile-tabbar"
    aria-label="移动端导航"
  >
    <RouterLink
      v-for="item in items"
      :key="item.key"
      class="mobile-tabbar__item"
      :class="{ active: item.active }"
      to="/"
      :aria-label="item.label"
    >
      <BaseIcon
        :name="item.icon"
        size="22"
      />
    </RouterLink>
    <button
      type="button"
      class="mobile-tabbar__item"
      :aria-label="authStore.isLoggedIn ? '我的' : '登录'"
      @click="authStore.isLoggedIn ? undefined : emit('login-request')"
    >
      <BaseIcon
        name="user"
        size="22"
      />
    </button>
  </nav>
</template>

<style scoped>
.mobile-tabbar {
  position: fixed;
  right: 0;
  bottom: 0;
  left: 0;
  z-index: 30;
  display: none;
  height: 58px;
  align-items: center;
  justify-content: space-around;
  padding: 0 max(10px, env(safe-area-inset-left)) env(safe-area-inset-bottom);
  background: rgba(255, 255, 255, 0.94);
  border-top: 1px solid var(--color-border);
  backdrop-filter: blur(18px);
}

.mobile-tabbar__item {
  display: inline-flex;
  width: 42px;
  height: 42px;
  align-items: center;
  justify-content: center;
  color: #666a72;
  border: 0;
  background: transparent;
}

.mobile-tabbar__item.active {
  color: #202124;
}

@media (max-width: 900px) {
  .mobile-tabbar {
    display: flex;
  }
}
</style>
