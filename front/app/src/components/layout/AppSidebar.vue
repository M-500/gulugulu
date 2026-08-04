<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { showToast } from 'vant'

import Avatar from '@/components/Avatar/avatar.vue'
import BaseIcon from '@/components/common/BaseIcon.vue'
import { primaryNavItems, secondaryNavItems } from '@/constants/navigation'
import { useAuthStore } from '@/stores/auth'

const emit = defineEmits(['login-request'])
const authStore = useAuthStore()
const moreMenuRef = ref(null)
const showMoreMenu = ref(false)
const visiblePrimaryNavItems = computed(() => primaryNavItems.filter((item) => item.key !== 'profile'))
const visibleSecondaryNavItems = computed(() => secondaryNavItems.filter((item) => item.key !== 'more'))

onMounted(() => {
  document.addEventListener('pointerdown', handleOutsidePointerDown)
  document.addEventListener('keydown', handleKeydown)
})

onBeforeUnmount(() => {
  document.removeEventListener('pointerdown', handleOutsidePointerDown)
  document.removeEventListener('keydown', handleKeydown)
})

function toggleMoreMenu() {
  showMoreMenu.value = !showMoreMenu.value
}

function handleLogout() {
  authStore.logout()
  showMoreMenu.value = false
  showToast('已退出登录')
}

function handleLogin() {
  showMoreMenu.value = false
  emit('login-request')
}

function handleOutsidePointerDown(event) {
  if (showMoreMenu.value && !moreMenuRef.value?.contains(event.target)) {
    showMoreMenu.value = false
  }
}

function handleKeydown(event) {
  if (event.key === 'Escape') {
    showMoreMenu.value = false
  }
}
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
      <img
        src="@/assets/brand/gulugulu-logo.png"
        alt="咕噜咕噜"
      >
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
      <div
        ref="moreMenuRef"
        class="sidebar-more"
      >
        <button
          type="button"
          class="sidebar-link sidebar-more__trigger"
          aria-haspopup="menu"
          :aria-expanded="showMoreMenu"
          @click="toggleMoreMenu"
        >
          <span class="sidebar-icon">
            <BaseIcon
              name="menu"
              size="20"
            />
          </span>
          <span>更多</span>
        </button>
        <div
          v-if="showMoreMenu"
          class="sidebar-more__menu"
          role="menu"
        >
          <button
            v-if="authStore.isLoggedIn"
            type="button"
            class="sidebar-more__action sidebar-more__action--danger"
            role="menuitem"
            @click="handleLogout"
          >
            <BaseIcon
              name="logout"
              size="18"
            />
            <span>退出登录</span>
          </button>
          <button
            v-else
            type="button"
            class="sidebar-more__action"
            role="menuitem"
            @click="handleLogin"
          >
            <BaseIcon
              name="user"
              size="18"
            />
            <span>登录</span>
          </button>
        </div>
      </div>
      <RouterLink
        v-for="item in visibleSecondaryNavItems"
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
  display: block;
  width: 158px;
  height: 53px;
  margin-bottom: 40px;
  overflow: hidden;
  border-radius: 13px;
  background: #f20b12;
  box-shadow: 0 8px 22px rgba(240, 11, 18, 0.18);
}

.brand img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.sidebar-nav {
  display: grid;
  width: 100%;
  gap: 12px;
}

.sidebar-nav-bottom {
  margin-top: auto;
}

.sidebar-more {
  position: relative;
}

.sidebar-more__trigger {
  width: 100%;
  border: 0;
  background: transparent;
  cursor: pointer;
  font: inherit;
  text-align: left;
}

.sidebar-more__menu {
  position: absolute;
  z-index: 20;
  right: 0;
  bottom: calc(100% + 10px);
  left: 0;
  border: 1px solid #eceef2;
  border-radius: 14px;
  background: #fff;
  box-shadow: 0 14px 38px rgba(30, 32, 38, 0.14);
  padding: 7px;
}

.sidebar-more__action {
  display: flex;
  width: 100%;
  min-height: 42px;
  align-items: center;
  gap: 10px;
  border: 0;
  border-radius: 10px;
  background: transparent;
  color: #30333a;
  cursor: pointer;
  font: inherit;
  font-size: 14px;
  font-weight: 700;
  padding: 0 12px;
}

.sidebar-more__action:hover {
  background: var(--color-fill);
}

.sidebar-more__action--danger {
  color: #e5484d;
}

.sidebar-more__action--danger:hover {
  background: #fff1f2;
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
