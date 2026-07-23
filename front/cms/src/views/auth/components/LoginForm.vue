<template>
  <form class="auth-form" @submit.prevent="handleSubmit">
    <AuthFormHeading
      title="欢迎回来"
      description="使用邮箱和密码进入内容管理后台"
    />

    <label class="form-field">
      <span>邮箱</span>
      <input v-model.trim="form.email" type="email" autocomplete="email" placeholder="name@example.com">
    </label>

    <label class="form-field">
      <span>密码</span>
      <input
        v-model="form.password"
        type="password"
        autocomplete="current-password"
        placeholder="至少 6 位密码"
      >
    </label>

    <CaptchaField
      v-model="form.captchaCode"
      :image-path="captcha.id ? captcha.path : ''"
      :loading="captchaLoading"
      @refresh="loadCaptcha"
    />

    <AuthFormMessage :error="errorMessage" :success="successMessage" />

    <button class="submit-button" type="submit" :disabled="submitting">
      {{ submitting ? '提交中...' : '登录后台' }}
    </button>
  </form>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

import { getCaptcha, login } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'
import { isEmail, isStrongEnoughPassword } from '@/utils/validators'

import AuthFormHeading from './AuthFormHeading.vue'
import AuthFormMessage from './AuthFormMessage.vue'
import CaptchaField from './CaptchaField.vue'

const router = useRouter()
const authStore = useAuthStore()

const submitting = ref(false)
const captchaLoading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const form = reactive({
  email: '',
  password: '',
  captchaCode: ''
})

const captcha = reactive({
  id: '',
  path: ''
})

async function loadCaptcha() {
  if (captchaLoading.value) {
    return
  }

  captchaLoading.value = true

  try {
    const data = await getCaptcha()

    if (!data?.captchaId || !data?.captchaPath) {
      throw new Error('验证码接口返回的数据格式不正确')
    }

    captcha.id = data.captchaId
    captcha.path = data.captchaPath
    form.captchaCode = ''
  } catch (error) {
    captcha.id = ''
    captcha.path = ''
    errorMessage.value = error.message
  } finally {
    captchaLoading.value = false
  }
}

function validate() {
  if (!isEmail(form.email)) {
    errorMessage.value = '请输入正确的邮箱地址'
    return false
  }

  if (!isStrongEnoughPassword(form.password)) {
    errorMessage.value = '密码至少需要 6 位'
    return false
  }

  if (!form.captchaCode || !captcha.id) {
    errorMessage.value = '请输入图片验证码'
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
    const payload = await login({ ...form, captchaId: captcha.id })
    authStore.setSession(payload)
    successMessage.value = '登录成功'
    await router.push({ name: 'Dashboard' })
  } catch (error) {
    errorMessage.value = error.message
    await loadCaptcha()
  } finally {
    submitting.value = false
  }
}

onMounted(loadCaptcha)
</script>

<style scoped>
@import './auth-form.css';

.auth-form {
  gap: 24px;
  padding-top: 38px;
  padding-bottom: 38px;
}
</style>
