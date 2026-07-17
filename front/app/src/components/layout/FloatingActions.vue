<script setup>
import BaseIcon from '@/components/common/BaseIcon.vue'

function scrollToTop() {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}
</script>

<template>
  <div
    class="floating-actions"
    aria-label="页面操作"
  >
    <button
      type="button"
      class="floating-actions__button"
      aria-label="返回顶部"
      @click="scrollToTop"
    >
      <BaseIcon
        name="up"
        size="20"
      />
    </button>
    <button
      type="button"
      class="floating-actions__button"
      aria-label="刷新推荐"
    >
      <BaseIcon
        name="refresh"
        size="18"
      />
    </button>
  </div>
</template>

<style scoped>
.floating-actions {
  position: fixed;
  right: 28px;
  bottom: 34px;
  z-index: 20;
  display: grid;
  gap: 12px;
}

.floating-actions__button {
  display: inline-flex;
  width: 42px;
  height: 42px;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--color-border);
  border-radius: 50%;
  color: #575b64;
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 8px 24px rgba(16, 24, 40, 0.08);
  backdrop-filter: blur(16px);
}

@media (max-width: 900px) {
  .floating-actions {
    right: 12px;
    bottom: 74px;
  }
}
</style>
