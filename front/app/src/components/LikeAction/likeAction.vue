<template>
  <button
    type="button"
    class="like-action"
    :class="{ 'is-liked': liked }"
    :aria-label="label"
    :aria-pressed="liked"
    :disabled="pending"
    @click="toggle"
  >
    <SvgIcon
      :name="liked ? 'like-filled' : 'like'"
      :size="iconSize"
    />
    <span>{{ countText }}</span>
  </button>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { showToast } from 'vant'

import SvgIcon from '@/components/common/SvgIcon.vue'
import { setResourceLike } from '@/services/interactiveService'
import { useAuthStore } from '@/stores/auth'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  count: {
    type: [Number, String],
    default: 0
  },
  label: {
    type: String,
    default: '点赞'
  },
  iconSize: {
    type: [Number, String],
    default: 16
  },
  resourceId: {
    type: [Number, String],
    default: 0
  },
  resourceType: {
    type: String,
    default: 'work'
  }
})

const emit = defineEmits(['update:modelValue', 'login-request', 'change'])

const liked = ref(props.modelValue)
const delta = ref(0)
const serverCount = ref(null)
const pending = ref(false)
const authStore = useAuthStore()

watch(() => props.modelValue, (value) => {
  liked.value = value
  delta.value = 0
  serverCount.value = null
})

const countText = computed(() => {
  const numericCount = serverCount.value ?? Number(props.count)
  if (Number.isNaN(numericCount)) return String(props.count || '赞')
  const nextCount = Math.max(0, numericCount + delta.value)
  return nextCount === 0 ? '赞' : String(nextCount)
})

async function toggle() {
  if (pending.value) return
  if (!authStore.isLoggedIn) {
    emit('login-request')
    showToast('请先登录后点赞')
    return
  }
  const resourceId = Number(props.resourceId)
  if (!Number.isInteger(resourceId) || resourceId <= 0) {
    showToast('该内容暂不支持点赞')
    return
  }

  const previousLiked = liked.value
  const nextLiked = !previousLiked
  liked.value = nextLiked
  delta.value += nextLiked ? 1 : -1
  pending.value = true
  emit('update:modelValue', nextLiked)
  try {
    const result = await setResourceLike({
      resourceType: props.resourceType,
      resourceId,
      liked: nextLiked,
      token: authStore.accessToken
    })
    liked.value = result.liked
    serverCount.value = Number(result.count || 0)
    delta.value = 0
    emit('update:modelValue', result.liked)
    emit('change', result)
  } catch (error) {
    liked.value = previousLiked
    delta.value += nextLiked ? -1 : 1
    emit('update:modelValue', previousLiked)
    showToast(error.message || '点赞失败，请稍后重试')
  } finally {
    pending.value = false
  }
}
</script>

<style scoped>
.like-action {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  border: 0;
  color: inherit;
  background: transparent;
  cursor: pointer;
  font: inherit;
  white-space: nowrap;
  padding: 0;
}

.like-action.is-liked {
  color: var(--color-primary);
}

.like-action:disabled {
  cursor: wait;
  opacity: 0.72;
}
</style>
