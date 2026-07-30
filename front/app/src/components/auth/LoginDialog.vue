<script setup>
import { onMounted, reactive, ref } from 'vue'
import { Button as VanButton, Checkbox as VanCheckbox, Field as VanField, Popup as VanPopup } from 'vant'

import { getCaptcha } from '@/services/authService'
import { useAuthStore } from '@/stores/auth'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:show', 'success'])
const authStore = useAuthStore()

const form = reactive({
  email: '',
  password: '',
  captchaCode: ''
})
const captcha = reactive({
  id: '',
  path: ''
})
const submitting = ref(false)
const captchaLoading = ref(false)
const agreed = ref(true)
const errorMessage = ref('')

async function loadCaptcha() {
  if (captchaLoading.value) return
  captchaLoading.value = true
  try {
    const data = await getCaptcha()
    captcha.id = data.captchaId || ''
    captcha.path = data.captchaPath || ''
    form.captchaCode = ''
  } catch (error) {
    errorMessage.value = error.message || '验证码加载失败'
  } finally {
    captchaLoading.value = false
  }
}

function closeDialog() {
  emit('update:show', false)
}

function validate() {
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) return '请输入正确的邮箱'
  if (!form.password || form.password.length < 6) return '密码至少 6 位'
  if (!captcha.id || !form.captchaCode) return '请输入图片验证码'
  if (!agreed.value) return '请先同意用户协议'
  return ''
}

async function submit() {
  const message = validate()
  if (message) {
    errorMessage.value = message
    return
  }

  submitting.value = true
  errorMessage.value = ''
  try {
    await authStore.login({
      email: form.email,
      password: form.password,
      captchaId: captcha.id,
      captchaCode: form.captchaCode
    })
    emit('success')
    closeDialog()
  } catch (error) {
    errorMessage.value = error.message || '登录失败'
    await loadCaptcha()
  } finally {
    submitting.value = false
  }
}

onMounted(loadCaptcha)
</script>

<template>
  <VanPopup
    :show="props.show"
    round
    closeable
    class="login-dialog"
    position="center"
    teleport="body"
    @update:show="$emit('update:show', $event)"
    @click-close-icon="closeDialog"
    @open="loadCaptcha"
  >
    <div class="login-dialog__inner">
      <section class="login-dialog__brand">
        <span>登录后推荐更懂你的笔记</span>
        <strong>咕噜</strong>
        <div class="login-dialog__qr">
          <span />
          <span />
          <span />
          <span />
        </div>
        <p>使用邮箱登录后，可以同步喜欢、收藏和创作记录。</p>
      </section>

      <form
        class="login-dialog__form"
        @submit.prevent="submit"
      >
        <h2>邮箱登录</h2>
        <VanField
          v-model.trim="form.email"
          type="email"
          autocomplete="email"
          placeholder="输入邮箱"
        />
        <VanField
          v-model="form.password"
          type="password"
          autocomplete="current-password"
          placeholder="输入密码"
        />
        <div class="login-dialog__captcha">
          <VanField
            v-model.trim="form.captchaCode"
            placeholder="输入验证码"
          />
          <button
            type="button"
            :disabled="captchaLoading"
            @click="loadCaptcha"
          >
            <img
              v-if="captcha.path"
              :src="captcha.path"
              alt="图片验证码"
            >
            <span v-else>{{ captchaLoading ? '加载中' : '刷新' }}</span>
          </button>
        </div>
        <p
          v-if="errorMessage"
          class="login-dialog__error"
        >
          {{ errorMessage }}
        </p>
        <VanButton
          block
          round
          type="primary"
          native-type="submit"
          :loading="submitting"
        >
          登录
        </VanButton>
        <VanCheckbox
          v-model="agreed"
          icon-size="15px"
          shape="square"
          class="login-dialog__agree"
        >
          我已阅读并同意《用户协议》《隐私政策》
        </VanCheckbox>
        <small>新用户登录成功后将自动完成初始化</small>
      </form>
    </div>
  </VanPopup>
</template>

<style scoped>
.login-dialog {
  width: min(760px, calc(100vw - 32px));
}

.login-dialog__inner {
  display: grid;
  grid-template-columns: 1fr 1fr;
  overflow: hidden;
  min-height: 440px;
}

.login-dialog__brand,
.login-dialog__form {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 34px 46px;
}

.login-dialog__brand {
  border-right: 1px solid #f0f1f4;
}

.login-dialog__brand > span {
  border-radius: 999px;
  background: #eef4ff;
  color: #3477e5;
  font-size: 13px;
  font-weight: 700;
  padding: 9px 18px;
}

.login-dialog__brand strong {
  margin-top: 28px;
  border-radius: 999px;
  background: var(--color-primary);
  color: #fff;
  font-size: 20px;
  padding: 7px 14px;
}

.login-dialog__qr {
  display: grid;
  width: 148px;
  height: 148px;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-top: 34px;
  border: 1px solid #f0f1f4;
  border-radius: 10px;
  padding: 18px;
}

.login-dialog__qr span {
  border-radius: 6px;
  background:
    linear-gradient(90deg, #111 35%, transparent 35% 65%, #111 65%),
    linear-gradient(#111 35%, transparent 35% 65%, #111 65%),
    #fff;
}

.login-dialog__brand p {
  margin: 22px 0 0;
  color: #717782;
  font-size: 13px;
  line-height: 1.7;
  text-align: center;
}

.login-dialog__form {
  align-items: stretch;
  gap: 14px;
}

.login-dialog__form h2 {
  margin: 0 0 22px;
  color: #20242c;
  font-size: 18px;
  text-align: center;
}

.login-dialog__form :deep(.van-cell) {
  border-radius: 999px;
  background: #f6f7f8;
  padding: 13px 18px;
}

.login-dialog__captcha {
  display: grid;
  grid-template-columns: 1fr 118px;
  gap: 10px;
}

.login-dialog__captcha button {
  overflow: hidden;
  border: 0;
  border-radius: 999px;
  background: #f6f7f8;
  color: var(--color-primary);
  font-weight: 700;
}

.login-dialog__captcha img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.login-dialog__error {
  margin: -2px 0 0;
  color: #e5485b;
  font-size: 12px;
}

.login-dialog__agree {
  margin-top: 2px;
  color: #7d8490;
  font-size: 12px;
}

.login-dialog__form small {
  color: #a5abb5;
  font-size: 12px;
  text-align: center;
}

@media (max-width: 720px) {
  .login-dialog__inner {
    grid-template-columns: 1fr;
  }

  .login-dialog__brand {
    display: none;
  }

  .login-dialog__form {
    padding: 44px 24px 30px;
  }
}
</style>
