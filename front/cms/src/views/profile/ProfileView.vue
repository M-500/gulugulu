<template>
  <section class="profile-page">
    <header class="profile-page__header">
      <div>
        <p>ACCOUNT SETTINGS</p>
        <h1>个人中心</h1>
        <span>管理你在咕噜咕噜 CMS 中展示的个人信息。</span>
      </div>
      <el-tag type="success" effect="light" round>账号状态正常</el-tag>
    </header>

    <div class="profile-layout">
      <el-card class="profile-card profile-card--avatar" shadow="never">
        <template #header>
          <div class="profile-card__title">
            <div>
              <strong>个人头像</strong>
              <span>用于顶部导航和内容作者信息展示</span>
            </div>
            <el-icon><Picture /></el-icon>
          </div>
        </template>

        <div class="profile-avatar">
          <el-avatar :size="128" :src="previewUrl || form.avatarUrl">
            {{ avatarText }}
          </el-avatar>
          <div>
            <strong>上传新的头像</strong>
            <p>支持 JPG、PNG、WEBP，文件大小不超过 5MB。</p>
            <el-upload
              :auto-upload="false"
              :on-change="handleAvatarChange"
              :show-file-list="false"
              accept="image/jpeg,image/png,image/webp"
            >
              <el-button :icon="Upload">选择图片</el-button>
            </el-upload>
          </div>
        </div>
      </el-card>

      <el-card class="profile-card" shadow="never">
        <template #header>
          <div class="profile-card__title">
            <div>
              <strong>基本资料</strong>
              <span>修改后会同步更新顶部用户信息</span>
            </div>
            <el-icon><EditPen /></el-icon>
          </div>
        </template>

        <el-form label-position="top" @submit.prevent="saveProfile">
          <div class="profile-form-grid">
            <el-form-item label="登录邮箱">
              <el-input v-model="form.email" disabled>
                <template #prefix><el-icon><Message /></el-icon></template>
              </el-input>
              <p class="profile-form-tip">登录邮箱暂不支持在个人中心修改。</p>
            </el-form-item>
            <el-form-item label="用户昵称">
              <el-input
                v-model.trim="form.nickName"
                maxlength="32"
                minlength="2"
                placeholder="请输入2到32个字符"
                show-word-limit
              >
                <template #prefix><el-icon><User /></el-icon></template>
              </el-input>
            </el-form-item>
            <el-form-item label="性别">
              <el-select v-model="form.sex" placeholder="请选择性别">
                <el-option label="保密" :value="0" />
                <el-option label="男" :value="1" />
                <el-option label="女" :value="2" />
              </el-select>
            </el-form-item>
            <el-form-item label="生日">
              <el-date-picker
                v-model="form.bothDay"
                type="date"
                value-format="YYYY-MM-DD"
                placeholder="请选择生日"
                :disabled-date="disableFutureDate"
                clearable
              />
            </el-form-item>
          </div>

          <el-form-item label="个性签名">
            <el-input
              v-model="form.bio"
              type="textarea"
              :rows="4"
              maxlength="200"
              placeholder="介绍一下自己吧"
              resize="none"
              show-word-limit
            />
          </el-form-item>

          <div class="profile-form-actions">
            <span>资料更新后立即生效</span>
            <el-button type="primary" :loading="saving" @click="saveProfile">保存修改</el-button>
          </div>
        </el-form>
      </el-card>
    </div>
  </section>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import { EditPen, Message, Picture, Upload, User } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

import { getCurrentUser, updateCurrentUser, uploadCurrentUserAvatar } from '@/api/profile'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const saving = ref(false)
const selectedAvatar = ref(null)
const previewUrl = ref('')
const form = reactive({
  email: '',
  nickName: '',
  avatarUrl: '',
  bio: '',
  sex: 0,
  bothDay: ''
})
const avatarText = computed(() => (form.nickName || form.email || '咕').slice(0, 1))

onMounted(loadProfile)
onBeforeUnmount(revokePreview)

async function loadProfile() {
  try {
    applyProfile(await getCurrentUser())
  } catch (error) {
    ElMessage.error(error.message || '加载个人资料失败')
  }
}

function applyProfile(profile) {
  Object.assign(form, profile)
  authStore.setProfile(profile)
}

function handleAvatarChange(uploadFile) {
  const file = uploadFile.raw
  if (!file) return
  if (!['image/jpeg', 'image/png', 'image/webp'].includes(file.type)) {
    ElMessage.warning('头像仅支持 JPG、PNG 或 WEBP 格式')
    return
  }
  if (file.size > 5 * 1024 * 1024) {
    ElMessage.warning('头像大小不能超过 5MB')
    return
  }
  revokePreview()
  selectedAvatar.value = file
  previewUrl.value = URL.createObjectURL(file)
}

async function saveProfile() {
  if (form.nickName.length < 2 || form.nickName.length > 32) {
    ElMessage.warning('昵称长度必须是2到32个字符')
    return
  }
  saving.value = true
  try {
    let profile = await updateCurrentUser({
      nickName: form.nickName,
      bio: form.bio,
      sex: form.sex,
      bothDay: form.bothDay || ''
    })
    if (selectedAvatar.value) {
      profile = await uploadCurrentUserAvatar(selectedAvatar.value)
    }
    selectedAvatar.value = null
    revokePreview()
    applyProfile(profile)
    ElMessage.success('个人资料已更新')
  } catch (error) {
    ElMessage.error(error.message || '保存个人资料失败')
  } finally {
    saving.value = false
  }
}

function disableFutureDate(date) {
  return date.getTime() > Date.now()
}

function revokePreview() {
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value)
    previewUrl.value = ''
  }
}
</script>

<style src="./profile.css"></style>
