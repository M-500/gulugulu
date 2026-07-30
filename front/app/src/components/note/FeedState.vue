<script setup>
import { Button as VanButton, Empty as VanEmpty, Loading as VanLoading } from 'vant'

defineProps({
  type: {
    type: String,
    required: true
  },
  message: {
    type: String,
    default: ''
  }
})

defineEmits(['retry'])
</script>

<template>
  <div
    class="feed-state"
    :class="`feed-state--${type}`"
  >
    <VanLoading
      v-if="type === 'loading'"
      color="#ff2442"
      size="28"
    />
    <VanEmpty
      v-else
      :image="type === 'error' ? 'error' : 'search'"
      :description="message"
    />
    <p v-if="type === 'loading'">
      {{ message }}
    </p>
    <VanButton
      v-if="type === 'error'"
      type="primary"
      round
      size="small"
      @click="$emit('retry')"
    >
      重新加载
    </VanButton>
  </div>
</template>

<style scoped>
.feed-state {
  display: flex;
  min-height: 220px;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 14px;
  color: #8a8f99;
  text-align: center;
}

.feed-state p {
  margin: 0;
  font-size: 14px;
}

.feed-state :deep(.van-empty) {
  padding: 0;
}
</style>
