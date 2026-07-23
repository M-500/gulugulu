<template>
  <p v-if="message" class="auth-form-message" :class="`is-${type}`">
    {{ message }}
  </p>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  error: {
    type: String,
    default: ''
  },
  success: {
    type: String,
    default: ''
  }
})

const message = computed(() => props.error || props.success)
const type = computed(() => (props.error ? 'error' : 'success'))
</script>

<style scoped>
.auth-form-message {
  min-height: 38px;
  margin: -4px 0 0;
  border-radius: 6px;
  font-size: 13px;
  line-height: 1.5;
  padding: 9px 12px;
}

.auth-form-message.is-error {
  background: #fef2f2;
  color: #b42318;
}

.auth-form-message.is-success {
  background: #ecfdf3;
  color: #067647;
}
</style>
