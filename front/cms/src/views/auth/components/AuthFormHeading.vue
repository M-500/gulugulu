<template>
  <div class="auth-form-heading">
    <h2>{{ title }}</h2>
    <p>{{ description }}</p>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true
  },
  description: {
    type: String,
    required: true
  }
})
</script>

<style scoped>
.auth-form-heading {
  display: grid;
  gap: 8px;
}

.auth-form-heading h2 {
  margin: 0;
  color: var(--color-text);
  font-size: 24px;
  line-height: 1.2;
}

.auth-form-heading p {
  margin: 0;
  color: var(--color-text-muted);
  font-size: 14px;
}
</style>
