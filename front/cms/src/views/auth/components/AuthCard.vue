<template>
  <div class="auth-card">
    <div class="auth-card__tabs" role="tablist" aria-label="账号操作">
      <button
        v-for="tab in tabs"
        :key="tab.value"
        type="button"
        :class="{ 'is-active': mode === tab.value }"
        role="tab"
        :aria-selected="mode === tab.value"
        @click="mode = tab.value"
      >
        {{ tab.label }}
      </button>
    </div>

    <LoginForm v-if="mode === 'login'" />
    <RegisterForm v-else />
  </div>
</template>

<script setup>
import { ref } from 'vue'

import LoginForm from './LoginForm.vue'
import RegisterForm from './RegisterForm.vue'

const tabs = [
  { label: '登录', value: 'login' },
  { label: '注册', value: 'register' }
]

const mode = ref('login')
</script>

<style scoped>
.auth-card {
  width: 100%;
  min-height: 590px;
  border: 1px solid rgba(17, 24, 39, 0.08);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 24px 80px rgba(15, 23, 42, 0.12);
  backdrop-filter: blur(16px);
}

.auth-card__tabs {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 6px;
  padding: 8px;
  border-bottom: 1px solid rgba(17, 24, 39, 0.08);
}

.auth-card__tabs button {
  height: 42px;
  border: 0;
  border-radius: 6px;
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
  font-size: 15px;
  font-weight: 700;
}

.auth-card__tabs button.is-active {
  background: #fff1f2;
  color: var(--color-primary);
}

@media (max-width: 480px) {
  .auth-card {
    min-height: auto;
  }
}
</style>
