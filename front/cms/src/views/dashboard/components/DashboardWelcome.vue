<template>
  <section class="dashboard-welcome">
    <AppLogo />
    <h1>内容管理后台骨架已就绪</h1>
    <p>后续可以在这里接入笔记、视频、长文、审核和用户管理模块。</p>
    <button type="button" @click="$emit('logout')">退出登录</button>
  </section>
</template>

<script setup>
import AppLogo from '@/components/AppLogo.vue'

defineEmits(['logout'])
</script>

<style scoped>
.dashboard-welcome {
  display: grid;
  width: min(620px, 100%);
  gap: 18px;
  border: 1px solid rgba(17, 24, 39, 0.08);
  border-radius: 8px;
  background: #fff;
  box-shadow: 0 18px 50px rgba(15, 23, 42, 0.08);
  padding: 36px;
}

.dashboard-welcome h1 {
  margin: 10px 0 0;
  color: var(--color-text);
  font-size: 26px;
}

.dashboard-welcome p {
  margin: 0;
  color: var(--color-text-muted);
  line-height: 1.8;
}

.dashboard-welcome button {
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
