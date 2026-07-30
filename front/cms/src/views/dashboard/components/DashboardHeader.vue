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

      <el-dropdown
        trigger="click"
        placement="bottom-end"
        popper-class="header-user-popper"
        :show-arrow="false"
        @command="handleCommand"
      >
        <button class="header-user" type="button">
          <el-avatar :size="30" :src="profile.avatarUrl">
            {{ avatarText }}
          </el-avatar>
          <strong class="header-user__name">{{ profile.nickName || '咕噜创作者' }}</strong>
          <el-icon class="header-user__arrow"><ArrowDown /></el-icon>
        </button>
        <template #dropdown>
          <el-dropdown-menu class="header-user-menu">
            <div class="header-user-menu__profile">
              <el-avatar :size="42" :src="profile.avatarUrl">
                {{ avatarText }}
              </el-avatar>
              <div>
                <strong>{{ profile.nickName || '咕噜创作者' }}</strong>
                <span>内容管理员</span>
              </div>
            </div>
            <div class="header-user-menu__label">账号管理</div>
            <el-dropdown-item command="profile">
              <span class="header-user-menu__icon">
                <el-icon><User /></el-icon>
              </span>
              <span class="header-user-menu__content">
                <strong>个人中心</strong>
                <small>编辑头像和个人资料</small>
              </span>
              <el-icon class="header-user-menu__chevron"><ArrowRight /></el-icon>
            </el-dropdown-item>
            <div class="header-user-menu__divider" />
            <el-dropdown-item class="header-user-menu__logout" command="logout">
              <span class="header-user-menu__icon">
                <el-icon><SwitchButton /></el-icon>
              </span>
              <span class="header-user-menu__content">
                <strong>退出登录</strong>
                <small>安全退出当前账号</small>
              </span>
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </header>
</template>

<script setup>
import { computed } from 'vue'
import { ArrowDown, ArrowRight, SwitchButton, User } from '@element-plus/icons-vue'
import AppLogo from '@/components/AppLogo.vue'

const props = defineProps({
  menuOpen: {
    type: Boolean,
    default: false
  },
  profile: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['toggle-menu', 'logout', 'profile'])
const avatarText = computed(() => (props.profile.nickName || '咕').slice(0, 1))

function handleCommand(command) {
  emit(command)
}
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
  gap: 8px;
  border: 0;
  border-radius: 999px;
  background: transparent;
  cursor: pointer;
  padding: 7px 10px;
  transition: background 0.2s ease;
}

.header-user:hover,
.header-user:focus-visible {
  background: #f5f6f8;
  outline: 0;
}

.header-user :deep(.el-avatar) {
  flex: 0 0 30px;
  background: linear-gradient(145deg, #ff7b82, #e63e4a);
  color: #fff;
  font-size: 13px;
  font-weight: 800;
}

.header-user__name {
  max-width: 140px;
  overflow: hidden;
  color: #1f2937;
  font-size: 14px;
  font-weight: 650;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.header-user__arrow {
  color: #a1a7b0;
  font-size: 12px;
  transition: transform 0.2s ease;
}

.header-user[aria-expanded='true'] .header-user__arrow {
  transform: rotate(180deg);
}

:global(.header-user-popper.el-popper) {
  width: 272px;
  overflow: hidden;
  border: 1px solid rgba(229, 231, 235, 0.9);
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.98);
  box-shadow: 0 18px 50px rgba(31, 41, 55, 0.14);
  backdrop-filter: blur(18px);
}

:global(.header-user-popper .el-dropdown-menu) {
  padding: 8px;
  background: transparent;
}

:global(.header-user-popper .header-user-menu__profile) {
  display: flex;
  align-items: center;
  gap: 11px;
  margin: 0 0 4px;
  border-radius: 12px;
  background: linear-gradient(135deg, #fff5f6, #fff 72%);
  padding: 12px;
}

:global(.header-user-popper .header-user-menu__profile .el-avatar) {
  flex: 0 0 42px;
  background: linear-gradient(145deg, #ff7b82, #e63e4a);
  color: #fff;
  font-weight: 800;
}

:global(.header-user-popper .header-user-menu__profile div) {
  min-width: 0;
}

:global(.header-user-popper .header-user-menu__profile strong),
:global(.header-user-popper .header-user-menu__profile span) {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

:global(.header-user-popper .header-user-menu__profile strong) {
  color: #303846;
  font-size: 13px;
}

:global(.header-user-popper .header-user-menu__profile span) {
  margin-top: 4px;
  color: #9ca3af;
  font-size: 10px;
}

:global(.header-user-popper .header-user-menu__label) {
  padding: 9px 11px 6px;
  color: #b0b5bd;
  font-size: 9px;
  font-weight: 700;
  letter-spacing: 0.12em;
}

:global(.header-user-popper .el-dropdown-menu__item) {
  display: flex;
  min-height: 54px;
  gap: 10px;
  border-radius: 11px;
  color: #47505e;
  line-height: normal;
  padding: 8px 10px;
  transition: background 0.18s ease, color 0.18s ease;
}

:global(.header-user-popper .el-dropdown-menu__item:not(.is-disabled):focus),
:global(.header-user-popper .el-dropdown-menu__item:not(.is-disabled):hover) {
  background: #f7f8fa;
  color: var(--color-primary);
}

:global(.header-user-popper .header-user-menu__icon) {
  display: grid;
  width: 34px;
  height: 34px;
  flex: 0 0 34px;
  place-items: center;
  border-radius: 10px;
  background: #fff0f1;
  color: var(--color-primary);
  font-size: 16px;
}

:global(.header-user-popper .header-user-menu__content) {
  display: grid;
  flex: 1;
  gap: 3px;
}

:global(.header-user-popper .header-user-menu__content strong) {
  font-size: 12px;
  font-weight: 650;
}

:global(.header-user-popper .header-user-menu__content small) {
  color: #a2a8b1;
  font-size: 9px;
}

:global(.header-user-popper .header-user-menu__chevron) {
  color: #b8bdc5;
  font-size: 12px;
}

:global(.header-user-popper .header-user-menu__divider) {
  height: 1px;
  margin: 5px 10px;
  background: #eef0f3;
}

:global(.header-user-popper .header-user-menu__logout .header-user-menu__icon) {
  background: #f5f6f8;
  color: #89909b;
}

:global(.header-user-popper .header-user-menu__logout:hover),
:global(.header-user-popper .header-user-menu__logout:focus) {
  background: #fff5f5 !important;
  color: #dc4c58 !important;
}

:global(.header-user-popper .header-user-menu__logout:hover .header-user-menu__icon),
:global(.header-user-popper .header-user-menu__logout:focus .header-user-menu__icon) {
  background: #ffe9eb;
  color: #dc4c58;
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

  .header-user__name,
  .header-action {
    display: none;
  }
}
</style>
