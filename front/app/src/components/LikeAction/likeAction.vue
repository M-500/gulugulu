<template>
  <button
    type="button"
    class="like-action"
    :class="{ 'is-liked': liked }"
    :aria-label="label"
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

import SvgIcon from '@/components/common/SvgIcon.vue'

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
  }
})

const emit = defineEmits(['update:modelValue'])

const liked = ref(props.modelValue)
const delta = ref(0)

watch(() => props.modelValue, (value) => {
  liked.value = value
  delta.value = 0
})

const countText = computed(() => {
  const numericCount = Number(props.count)
  if (Number.isNaN(numericCount)) return String(props.count || '赞')
  const nextCount = Math.max(0, numericCount + delta.value)
  return nextCount === 0 ? '赞' : String(nextCount)
})

function toggle() {
  liked.value = !liked.value
  delta.value += liked.value ? 1 : -1
  emit('update:modelValue', liked.value)
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
</style>
