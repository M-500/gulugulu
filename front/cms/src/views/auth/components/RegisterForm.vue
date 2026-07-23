<template>
  <form class="auth-form" @submit.prevent="handleSubmit">
    <AuthFormHeading
      title="创建管理员账号"
      description="填写邮箱验证码完成账号注册"
    />

    <label class="form-field">
      <span>昵称</span>
      <input v-model.trim="form.nickName" type="text" autocomplete="nickname" placeholder="请输入昵称">
    </label>

    <label class="form-field">
      <span>邮箱</span>
      <input v-model.trim="form.email" type="email" autocomplete="email" placeholder="name@example.com">
    </label>

    <label class="form-field">
      <span>密码</span>
      <input v-model="form.password" type="password" autocomplete="new-password" placeholder="至少 6 位密码">
    </label>

    <div class="form-field">
      <span>邮箱验证码</span>
      <div class="code-row">
        <input
          v-model.trim="form.verificationCode"
          type="text"
          autocomplete="one-time-code"
          placeholder="请输入验证码"
        >
        <button type="button" :disabled="codeCountdown > 0" @click="startCodeCountdown">
          {{ codeCountdown > 0 ? `${codeCountdown}s` : '获取验证码' }}
        </button>
      </div>
    </div>

    <AuthFormMessage :error="errorMessage" :success="successMessage" />

    <button class="submit-button" type="submit" :disabled="submitting">
      {{ submitting ? '提交中...' : '注册账号' }}
    </button>
  </form>
</template>

<script setup>
import { onBeforeUnmount, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

import { register } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'
import { isEmail, isStrongEnoughPassword } from '@/utils/validators'

import AuthFormHeading from './AuthFormHeading.vue'
import AuthFormMessage from './AuthFormMessage.vue'

const router = useRouter()
const authStore = useAuthStore()

const submitting = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const codeCountdown = ref(0)
let codeTimer = null

const form = reactive({
  nickName: '',
  email: '',
  password: '',
  verificationCode: ''
})

function validate() {
  if (!form.nickName) {
    errorMessage.value = '请输入昵称'
    return false
  }

  if (!isEmail(form.email)) {
    errorMessage.value = '请输入正确的邮箱地址'
    return false
  }

  if (!isStrongEnoughPassword(form.password)) {
    errorMessage.value = '密码至少需要 6 位'
    return false
  }

  if (!form.verificationCode) {
    errorMessage.value = '请输入邮箱验证码'
    return false
  }

  return true
}

async function handleSubmit() {
  errorMessage.value = ''
  successMessage.value = ''

  if (!validate()) {
    return
  }

  submitting.value = true

  try {
    const payload = await register(form)
    authStore.setSession(payload)
    successMessage.value = '注册成功'
    await router.push({ name: 'Dashboard' })
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    submitting.value = false
  }
}

function startCodeCountdown() {
  if (!isEmail(form.email)) {
    errorMessage.value = '请先输入正确的邮箱地址'
    return
  }

  errorMessage.value = ''
  successMessage.value = '邮箱验证码发送接口接入后即可发送验证码'
  codeCountdown.value = 60
  clearInterval(codeTimer)
  codeTimer = setInterval(() => {
    codeCountdown.value -= 1

    if (codeCountdown.value <= 0) {
      clearInterval(codeTimer)
    }
  }, 1000)
}

onBeforeUnmount(() => clearInterval(codeTimer))
</script>

<style scoped>
@import './auth-form.css';

.code-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 118px;
  gap: 10px;
}

.code-row button {
  height: 46px;
  border: 1px solid var(--color-border);
  border-radius: 6px;
  background: #fff;
  color: var(--color-text);
  cursor: pointer;
  font-weight: 700;
}

.code-row button:disabled {
  cursor: not-allowed;
  opacity: 0.62;
}

@media (max-width: 480px) {
  .code-row {
    grid-template-columns: 1fr;
  }
}
</style>
