
<template>
  <span
    class="avatar"
    :style="avatarStyle"
    role="button"
    tabindex="0"
    @click="$emit('click')"
    @keydown.enter.prevent="$emit('click')"
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
  }
})

defineEmits(['click'])

const normalizedSize = computed(() => {
  return typeof props.size === 'number' ? `${props.size}px` : props.size
})

const avatarStyle = computed(() => ({
  width: normalizedSize.value,
  height: normalizedSize.value,
  flexBasis: normalizedSize.value
}))

const initial = computed(() => {
  return String(props.username || '咕').trim().slice(0, 1).toUpperCase()
})
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
  cursor: pointer;
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
