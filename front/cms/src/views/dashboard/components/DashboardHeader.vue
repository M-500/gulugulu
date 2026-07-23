<template>
  <header class="dashboard-header">
    <div class="dashboard-header__left">
      <button
        class="dashboard-header__menu"
        type="button"
        :aria-expanded="menuOpen"
        aria-label="切换侧边菜单"
        @click="$emit('toggle-menu')"
      >
        <span />
        <span />
        <span />
      </button>
      <AppLogo />
    </div>

    <div class="dashboard-header__right">
      <button class="header-action" type="button" aria-label="查看通知">
        <span class="header-action__icon">•</span>
        <span class="header-action__badge" />
      </button>

      <div class="header-user">
        <span class="header-user__avatar">咕</span>
        <span class="header-user__info">
          <strong>咕噜创作者</strong>
          <small>内容管理员</small>
        </span>
      </div>

      <button class="logout-button" type="button" @click="$emit('logout')">
        退出
      </button>
    </div>
  </header>
</template>

<script setup>
import AppLogo from '@/components/AppLogo.vue'

defineProps({
  menuOpen: {
    type: Boolean,
    default: false
  }
})

defineEmits(['toggle-menu', 'logout'])
</script>

<style scoped>
.dashboard-header {
  position: sticky;
  z-index: 30;
  top: 0;
  display: flex;
  height: 68px;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  border-bottom: 1px solid #e8eaf0;
  background: rgba(255, 255, 255, 0.94);
  padding: 0 28px;
  backdrop-filter: blur(14px);
}

.dashboard-header__left,
.dashboard-header__right,
.header-user {
  display: flex;
  align-items: center;
}

.dashboard-header__left {
  gap: 14px;
}

.dashboard-header__right {
  gap: 18px;
}

.dashboard-header__menu {
  display: none;
  width: 38px;
  height: 38px;
  place-content: center;
  gap: 4px;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  background: #fff;
  cursor: pointer;
}

.dashboard-header__menu span {
  display: block;
  width: 16px;
  height: 2px;
  border-radius: 2px;
  background: #374151;
}

.header-action {
  position: relative;
  display: grid;
  width: 38px;
  height: 38px;
  place-items: center;
  border: 1px solid #e8eaf0;
  border-radius: 50%;
  background: #fff;
  color: #4b5563;
  cursor: pointer;
}

.header-action__icon {
  width: 12px;
  height: 14px;
  border: 2px solid currentColor;
  border-radius: 8px 8px 5px 5px;
  color: #4b5563;
  font-size: 0;
}

.header-action__badge {
  position: absolute;
  top: 3px;
  right: 3px;
  width: 7px;
  height: 7px;
  border: 2px solid #fff;
  border-radius: 50%;
  background: var(--color-primary);
}

.header-user {
  gap: 10px;
}

.header-user__avatar {
  display: grid;
  width: 38px;
  height: 38px;
  place-items: center;
  border-radius: 12px;
  background: linear-gradient(145deg, #ff7b82, #e63e4a);
  color: #fff;
  font-size: 15px;
  font-weight: 800;
}

.header-user__info {
  display: grid;
  gap: 2px;
}

.header-user__info strong {
  color: #1f2937;
  font-size: 13px;
}

.header-user__info small {
  color: #9ca3af;
  font-size: 11px;
}

.logout-button {
  border: 0;
  background: transparent;
  color: #6b7280;
  cursor: pointer;
  font-size: 13px;
  padding: 8px 4px;
}

.logout-button:hover {
  color: var(--color-primary);
}

@media (max-width: 900px) {
  .dashboard-header {
    padding: 0 20px;
  }

  .dashboard-header__menu {
    display: grid;
  }
}

@media (max-width: 560px) {
  .dashboard-header {
    padding: 0 16px;
  }

  .header-user__info,
  .header-action {
    display: none;
  }
}
</style>
