
<template>
  <button
    type="button"
    class="author-wrapper"
    :class="{ 'is-inline': inline }"
    @click="$emit('click')"
  >
    <span class="author-wrapper__name">{{ nickname || '咕噜用户' }}</span>
    <small v-if="isAuthor">作者</small>
  </button>
</template>

<script setup>
defineProps({
  nickname: {
    type: String,
    default: ''
  },
  isAuthor: {
    type: Boolean,
    default: false
  },
  inline: {
    type: Boolean,
    default: false
  }
})

defineEmits(['click'])
</script>

<style scoped>
.author-wrapper {
  display: inline-flex;
  min-width: 0;
  align-items: center;
  gap: 6px;
  border: 0;
  color: inherit;
  background: transparent;
  cursor: pointer;
  font: inherit;
  padding: 0;
  text-align: left;
}

.author-wrapper__name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.author-wrapper small {
  flex: 0 0 auto;
  border-radius: 999px;
  color: #777;
  background: #f2f3f5;
  font-size: 10px;
  font-weight: 400;
  padding: 2px 5px;
}

.author-wrapper.is-inline {
  vertical-align: middle;
}
</style>
