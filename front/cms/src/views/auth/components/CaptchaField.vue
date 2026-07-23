<template>
  <div class="form-field">
    <span>图片验证码</span>
    <div class="captcha-row">
      <input
        :value="modelValue"
        type="text"
        autocomplete="off"
        placeholder="请输入验证码"
        @input="$emit('update:modelValue', $event.target.value.trim())"
      >
      <button
        class="captcha-image"
        type="button"
        title="点击刷新验证码"
        :disabled="loading"
        @click="$emit('refresh')"
      >
        <img v-if="imagePath" :src="imagePath" alt="图片验证码">
        <span v-else>{{ loading ? '加载中...' : '刷新' }}</span>
      </button>
    </div>
  </div>
</template>

<script setup>
defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  imagePath: {
    type: String,
    default: ''
  },
  loading: {
    type: Boolean,
    default: false
  }
})

defineEmits(['update:modelValue', 'refresh'])
</script>

<style scoped>
.form-field {
  display: grid;
  gap: 8px;
}

.form-field > span {
  color: var(--color-text);
  font-size: 14px;
  font-weight: 700;
}

.captcha-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 118px;
  gap: 10px;
}

.captcha-row input {
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

.captcha-row input:focus {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(239, 77, 88, 0.14);
}

.captcha-image {
  height: 46px;
  overflow: hidden;
  border: 1px solid var(--color-border);
  border-radius: 6px;
  background: #fff;
  color: var(--color-text);
  cursor: pointer;
  font-weight: 700;
  padding: 0;
}

.captcha-image:disabled {
  cursor: wait;
  opacity: 0.7;
}

.captcha-image img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: contain;
}

@media (max-width: 480px) {
  .captcha-row {
    grid-template-columns: 1fr;
  }
}
</style>
