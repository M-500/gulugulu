<template>
  <div class="works-grid">
    <WorkCard
      v-for="work in works"
      :key="work.workId"
      :work="work"
      @open="$emit('open', $event)"
      @delete="$emit('delete', $event)"
      @edit="$emit('edit', $event)"
      @visibility="$emit('visibility', $event)"
    />
  </div>
</template>

<script setup>
import WorkCard from './WorkCard.vue'

defineProps({
  works: { type: Array, default: () => [] }
})

defineEmits(['delete', 'edit', 'open', 'visibility'])
</script>
