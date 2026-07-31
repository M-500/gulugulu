<script setup>
import { computed } from 'vue'

import Avatar from '@/components/Avatar/avatar.vue'
import BaseIcon from '@/components/common/BaseIcon.vue'
import { primaryNavItems, secondaryNavItems } from '@/constants/navigation'
import { useAuthStore } from '@/stores/auth'

const emit = defineEmits(['login-request'])
const authStore = useAuthStore()
const visiblePrimaryNavItems = computed(() => primaryNavItems.filter((item) => item.key !== 'profile'))
</script>

<template>
  <aside
    class="app-sidebar"
    aria-label="主导航"
  >
    <RouterLink
      class="brand"
      to="/"
      aria-label="咕噜咕噜首页"
    >
      咕噜
    </RouterLink>

    <nav class="sidebar-nav">
      <RouterLink
        v-for="item in visiblePrimaryNavItems"
        :key="item.key"
        class="sidebar-link"
        :class="{ active: item.active }"
        to="/"
      >
        <span class="sidebar-icon">
          <BaseIcon
            :name="item.icon"
            size="20"
          />
        </span>
        <span>{{ item.label }}</span>
        <small
          v-if="item.badge"
          class="sidebar-badge"
        >{{ item.badge }}</small>
      </RouterLink>
      <button
        v-if="!authStore.isLoggedIn"
        type="button"
        class="sidebar-link sidebar-auth-entry sidebar-auth-entry--login"
        @click="emit('login-request')"
      >
        <span class="sidebar-icon">
          <BaseIcon
            name="user"
            size="20"
          />
        </span>
        <span>登录</span>
      </button>
      <RouterLink
        v-else
        class="sidebar-link sidebar-auth-entry sidebar-auth-entry--profile"
        :to="`/users/${authStore.user?.userId || 'me'}`"
      >
        <span class="sidebar-icon sidebar-profile-avatar">
          <Avatar
            :avatar-url="authStore.user?.avatarUrl"
            :username="authStore.user?.nickName || '我'"
            :size="24"
          />
        </span>
        <span class="sidebar-auth-entry__label">我</span>
      </RouterLink>
    </nav>

    <nav
      class="sidebar-nav sidebar-nav-bottom"
      aria-label="辅助导航"
    >
      <RouterLink
        v-for="item in secondaryNavItems"
        :key="item.key"
        class="sidebar-link"
        to="/"
      >
        <span class="sidebar-icon">
          <BaseIcon
            :name="item.icon"
            size="20"
          />
        </span>
        <span>{{ item.label }}</span>
      </RouterLink>
    </nav>
  </aside>
</template>

<style scoped>
.app-sidebar {
  position: sticky;
  top: 0;
  display: flex;
  width: var(--layout-sidebar);
  height: 100vh;
  flex: 0 0 var(--layout-sidebar);
  flex-direction: column;
  align-items: center;
  padding: 30px 28px 26px;
  background: #fff;
  border-right: 1px solid #f4f4f5;
}

.brand {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 76px;
  height: 36px;
  margin-bottom: 48px;
  border-radius: 999px;
  color: #fff;
  font-weight: 900;
  font-size: 18px;
  line-height: 1;
  letter-spacing: 0;
  background: var(--color-primary);
}

.sidebar-nav {
  display: grid;
  width: 100%;
  gap: 12px;
}

.sidebar-nav-bottom {
  margin-top: auto;
}

.sidebar-link {
  display: flex;
  min-height: 54px;
  align-items: center;
  gap: 14px;
  padding: 0 20px;
  border-radius: 999px;
  color: #202124;
  font-weight: 700;
  transition: background 0.2s ease, color 0.2s ease;
}

.sidebar-auth-entry {
  width: 100%;
}

.sidebar-auth-entry--login {
  border: 0;
  background: transparent;
  font: inherit;
  text-align: left;
  cursor: pointer;
}

.sidebar-auth-entry--profile {
  overflow: hidden;
}

.sidebar-auth-entry__label {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.sidebar-link:hover,
.sidebar-link.active {
  background: var(--color-fill);
}

.sidebar-icon {
  display: inline-flex;
  width: 24px;
  height: 24px;
  flex: 0 0 24px;
  align-items: center;
  justify-content: center;
  color: #383a40;
}

.sidebar-profile-avatar {
  overflow: hidden;
  border: 1px solid #eceef2;
  border-radius: 50%;
}

.sidebar-badge {
  padding: 1px 5px;
  border-radius: 6px;
  color: #12a675;
  font-size: 10px;
  font-weight: 800;
  background: #ddfaee;
}

@media (max-width: 900px) {
  .app-sidebar {
    display: none;
  }
}
</style>
