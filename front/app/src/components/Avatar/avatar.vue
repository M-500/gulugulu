
<template>
  <span
    class="avatar"
    :class="{ 'avatar--interactive': interactive }"
    :style="avatarStyle"
    :role="interactive ? 'button' : undefined"
    :tabindex="interactive ? 0 : undefined"
    @click="handleClick"
    @keydown.enter.prevent="handleClick"
    @keydown.space.prevent="handleClick"
  >
    <img
      v-if="avatarUrl"
      :src="avatarUrl"
      :alt="username"
    >
    <span v-else>{{ initial }}</span>
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  avatarUrl: {
    type: String,
    default: ''
  },
  username: {
    type: String,
    default: ''
  },
  size: {
    type: [Number, String],
    default: 20
  },
  interactive: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['click'])

const normalizedSize = computed(() => {
  if (typeof props.size === 'number') return `${props.size}px`

  const size = String(props.size).trim()
  return /^\d+(\.\d+)?$/.test(size) ? `${size}px` : size
})

const avatarStyle = computed(() => ({
  width: normalizedSize.value,
  height: normalizedSize.value,
  flexBasis: normalizedSize.value
}))

const initial = computed(() => {
  return String(props.username || '咕').trim().slice(0, 1).toUpperCase()
})

function handleClick () {
  if (props.interactive) emit('click')
}
</script>

<style scoped>
.avatar {
  display: inline-flex;
  overflow: hidden;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: #fff;
  font-size: 11px;
  font-weight: 800;
  line-height: 1;
  background: #2e3036;
  user-select: none;
}

.avatar--interactive {
  cursor: pointer;
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
