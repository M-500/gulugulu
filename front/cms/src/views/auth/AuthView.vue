<template>
  <AuthLayout>
    <div class="auth-card">
      <div class="auth-card__tabs" role="tablist" aria-label="账号操作">
        <button
          type="button"
          :class="{ 'is-active': mode === 'login' }"
          role="tab"
          :aria-selected="mode === 'login'"
          @click="switchMode('login')"
        >
          登录
        </button>
        <button
          type="button"
          :class="{ 'is-active': mode === 'register' }"
          role="tab"
          :aria-selected="mode === 'register'"
          @click="switchMode('register')"
        >
          注册
        </button>
      </div>

      <form class="auth-form" :class="`auth-form--${mode}`" @submit.prevent="handleSubmit">
        <div class="auth-form__heading">
          <h2>{{ mode === 'login' ? '欢迎回来' : '创建管理员账号' }}</h2>
          <p>{{ mode === 'login' ? '使用邮箱和密码进入内容管理后台' : '填写邮箱验证码完成账号注册' }}</p>
        </div>

        <label v-if="mode === 'register'" class="form-field">
          <span>昵称</span>
          <input v-model.trim="registerForm.nickName" type="text" autocomplete="nickname" placeholder="请输入昵称">
        </label>

        <label class="form-field">
          <span>邮箱</span>
          <input v-model.trim="activeForm.email" type="email" autocomplete="email" placeholder="name@example.com">
        </label>

        <label class="form-field">
          <span>密码</span>
          <input
            v-model="activeForm.password"
            type="password"
            :autocomplete="mode === 'login' ? 'current-password' : 'new-password'"
            placeholder="至少 6 位密码"
          >
        </label>

        <div v-if="mode === 'login'" class="form-field">
          <span>图片验证码</span>
          <div class="captcha-row">
            <input v-model.trim="loginForm.captchaCode" type="text" autocomplete="off" placeholder="请输入验证码">
            <button class="captcha-image" type="button" title="点击刷新验证码" @click="loadCaptcha">
              <img v-if="captcha.path" :src="captcha.path" alt="图片验证码">
              <span v-else>刷新</span>
            </button>
          </div>
        </div>

        <div v-else class="form-field">
          <span>邮箱验证码</span>
          <div class="code-row">
            <input v-model.trim="registerForm.verificationCode" type="text" autocomplete="one-time-code" placeholder="请输入验证码">
            <button type="button" :disabled="codeCountdown > 0" @click="startCodeCountdown">
              {{ codeCountdown > 0 ? `${codeCountdown}s` : '获取验证码' }}
            </button>
          </div>
        </div>

        <p v-if="errorMessage" class="auth-form__message is-error">{{ errorMessage }}</p>
        <p v-if="successMessage" class="auth-form__message is-success">{{ successMessage }}</p>

        <button class="submit-button" type="submit" :disabled="submitting">
          {{ submitting ? '提交中...' : mode === 'login' ? '登录后台' : '注册账号' }}
        </button>
      </form>
    </div>
  </AuthLayout>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

import { getCaptcha, login, register } from '@/api/auth'
import AuthLayout from '@/layouts/AuthLayout.vue'
import { useAuthStore } from '@/stores/auth'
import { isEmail, isStrongEnoughPassword } from '@/utils/validators'

const router = useRouter()
const authStore = useAuthStore()

const mode = ref('login')
const submitting = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const codeCountdown = ref(0)
let codeTimer = null

const loginForm = reactive({
  email: '',
  password: '',
  captchaCode: ''
})

const registerForm = reactive({
  nickName: '',
  email: '',
  password: '',
  verificationCode: ''
})

const captcha = reactive({
  id: '',
  path: ''
})

const activeForm = computed(() => (mode.value === 'login' ? loginForm : registerForm))

function switchMode(nextMode) {
  mode.value = nextMode
  errorMessage.value = ''
  successMessage.value = ''

  if (nextMode === 'login' && !captcha.id) {
    loadCaptcha()
  }
}

async function loadCaptcha() {
  try {
    const data = await getCaptcha()
    captcha.id = data.captchaId
    captcha.path = data.captchaPath
  } catch (error) {
    errorMessage.value = error.message
  }
}

function validateBaseForm(form) {
  if (!isEmail(form.email)) {
    errorMessage.value = '请输入正确的邮箱地址'
    return false
  }

  if (!isStrongEnoughPassword(form.password)) {
    errorMessage.value = '密码至少需要 6 位'
    return false
  }

  return true
}

function validateLoginForm() {
  if (!validateBaseForm(loginForm)) {
    return false
  }

  if (!loginForm.captchaCode || !captcha.id) {
    errorMessage.value = '请输入图片验证码'
    return false
  }

  return true
}

function validateRegisterForm() {
  if (!registerForm.nickName) {
    errorMessage.value = '请输入昵称'
    return false
  }

  if (!validateBaseForm(registerForm)) {
    return false
  }

  if (!registerForm.verificationCode) {
    errorMessage.value = '请输入邮箱验证码'
    return false
  }

  return true
}

async function handleSubmit() {
  errorMessage.value = ''
  successMessage.value = ''

  const isValid = mode.value === 'login' ? validateLoginForm() : validateRegisterForm()
  if (!isValid) {
    return
  }

  submitting.value = true

  try {
    const payload = mode.value === 'login'
      ? await login({ ...loginForm, captchaId: captcha.id })
      : await register(registerForm)

    authStore.setSession(payload)
    successMessage.value = mode.value === 'login' ? '登录成功' : '注册成功'
    await router.push({ name: 'Dashboard' })
  } catch (error) {
    errorMessage.value = error.message

    if (mode.value === 'login') {
      await loadCaptcha()
    }
  } finally {
    submitting.value = false
  }
}

function startCodeCountdown() {
  if (!isEmail(registerForm.email)) {
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

onMounted(loadCaptcha)
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

.auth-form {
  display: grid;
  min-height: 531px;
  align-content: start;
  gap: 18px;
  padding: 30px;
}

.auth-form--login {
  gap: 24px;
  padding-top: 38px;
  padding-bottom: 38px;
}

.auth-form--register {
  gap: 18px;
}

.auth-form__heading {
  display: grid;
  gap: 8px;
}

.auth-form--login .auth-form__heading {
  margin-bottom: 8px;
}

.auth-form__heading h2 {
  margin: 0;
  color: var(--color-text);
  font-size: 24px;
  line-height: 1.2;
}

.auth-form__heading p {
  margin: 0;
  color: var(--color-text-muted);
  font-size: 14px;
}

.form-field {
  display: grid;
  gap: 8px;
}

.form-field span {
  color: var(--color-text);
  font-size: 14px;
  font-weight: 700;
}

.form-field input {
  width: 100%;
  height: 46px;
  box-sizing: border-box;
  border: 1px solid var(--color-border);
  border-radius: 6px;
  background: #fff;
  color: var(--color-text);
  font-size: 15px;
  outline: none;
  padding: 0 14px;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.form-field input:focus {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(239, 77, 88, 0.14);
}

.captcha-row,
.code-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 118px;
  gap: 10px;
}

.captcha-image,
.code-row button {
  height: 46px;
  border: 1px solid var(--color-border);
  border-radius: 6px;
  background: #fff;
  color: var(--color-text);
  cursor: pointer;
  font-weight: 700;
}

.captcha-image {
  overflow: hidden;
  padding: 0;
}

.captcha-image img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.code-row button:disabled {
  cursor: not-allowed;
  opacity: 0.62;
}

.auth-form__message {
  min-height: 38px;
  margin: -4px 0 0;
  border-radius: 6px;
  font-size: 13px;
  line-height: 1.5;
  padding: 9px 12px;
}

.auth-form__message.is-error {
  background: #fef2f2;
  color: #b42318;
}

.auth-form__message.is-success {
  background: #ecfdf3;
  color: #067647;
}

.submit-button {
  height: 48px;
  border: 0;
  border-radius: 6px;
  background: var(--color-primary);
  color: #fff;
  cursor: pointer;
  font-size: 16px;
  font-weight: 800;
  transition: background 0.2s ease, transform 0.2s ease;
}

.submit-button:hover:not(:disabled) {
  background: #dc3545;
  transform: translateY(-1px);
}

.submit-button:disabled {
  cursor: not-allowed;
  opacity: 0.7;
}

@media (max-width: 480px) {
  .auth-card,
  .auth-form {
    min-height: auto;
  }

  .auth-form {
    padding: 24px 20px;
  }

  .captcha-row,
  .code-row {
    grid-template-columns: 1fr;
  }
}
</style>
