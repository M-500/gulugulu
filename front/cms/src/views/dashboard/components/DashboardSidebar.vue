<template>
  <aside class="dashboard-sidebar" :class="{ 'is-open': open }">
    <nav class="dashboard-sidebar__nav" aria-label="后台功能菜单">
      <button
        v-for="item in dashboardMenus"
        :key="item.key"
        type="button"
        :class="{ 'is-active': item.key === activeItem }"
        @click="$emit('select', item.key)"
      >
        <span class="menu-icon">{{ item.icon }}</span>
        <span>{{ item.label }}</span>
        <span v-if="item.key === 'review'" class="menu-count">6</span>
      </button>
    </nav>

    <div class="sidebar-help">
      <span class="sidebar-help__icon">?</span>
      <div>
        <strong>需要帮助？</strong>
        <p>查看创作者中心使用指南</p>
      </div>
    </div>

    <p class="dashboard-sidebar__version">GULUGULU CMS · V1.0</p>
  </aside>
</template>

<script setup>
import { dashboardMenus } from '../dashboardMenus'

defineProps({
  activeItem: {
    type: String,
    required: true
  },
  open: {
    type: Boolean,
    default: false
  }
})

defineEmits(['select'])
</script>

<style scoped>
.dashboard-sidebar {
  position: sticky;
  top: 68px;
  display: flex;
  width: 232px;
  height: calc(100vh - 68px);
  flex: 0 0 232px;
  flex-direction: column;
  border-right: 1px solid #e8eaf0;
  background: #fff;
  padding: 24px 14px 18px;
}

.dashboard-sidebar__nav {
  display: grid;
  gap: 6px;
}

.dashboard-sidebar__nav button {
  display: grid;
  min-height: 46px;
  grid-template-columns: 28px 1fr auto;
  align-items: center;
  gap: 8px;
  border: 0;
  border-radius: 10px;
  background: transparent;
  color: #667085;
  cursor: pointer;
  font-size: 14px;
  font-weight: 600;
  padding: 0 14px;
  text-align: left;
  transition: background 0.2s ease, color 0.2s ease;
}

.dashboard-sidebar__nav button:hover {
  background: #f8f9fb;
  color: #1f2937;
}

.dashboard-sidebar__nav button.is-active {
  background: #fff0f1;
  color: #dd3f4b;
}

.menu-icon {
  display: grid;
  width: 25px;
  height: 25px;
  place-items: center;
  border-radius: 7px;
  background: #f1f3f6;
  color: #6b7280;
  font-size: 16px;
  line-height: 1;
}

.is-active .menu-icon {
  background: #fff;
  color: var(--color-primary);
}

.menu-count {
  display: grid;
  min-width: 20px;
  height: 20px;
  place-items: center;
  border-radius: 10px;
  background: #f3f4f6;
  color: #667085;
  font-size: 11px;
}

.sidebar-help {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  margin-top: auto;
  border-radius: 12px;
  background: #f8f9fb;
  padding: 14px;
}

.sidebar-help__icon {
  display: grid;
  width: 24px;
  height: 24px;
  flex: 0 0 24px;
  place-items: center;
  border-radius: 8px;
  background: #fff;
  color: var(--color-primary);
  font-size: 12px;
  font-weight: 800;
}

.sidebar-help strong {
  color: #374151;
  font-size: 12px;
}

.sidebar-help p {
  margin: 4px 0 0;
  color: #9ca3af;
  font-size: 11px;
  line-height: 1.5;
}

.dashboard-sidebar__version {
  margin: 16px 0 0;
  color: #c0c4cc;
  font-size: 9px;
  letter-spacing: 0.08em;
  text-align: center;
}

@media (max-width: 900px) {
  .dashboard-sidebar {
    position: fixed;
    z-index: 25;
    top: 68px;
    bottom: 0;
    left: 0;
    height: auto;
    transform: translateX(-100%);
    transition: transform 0.25s ease;
  }

  .dashboard-sidebar.is-open {
    transform: translateX(0);
  }
}
</style>
