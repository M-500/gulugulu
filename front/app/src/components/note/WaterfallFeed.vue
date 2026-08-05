<script setup>
import NoteCard from '@/components/note/NoteCard.vue'

defineProps({
  notes: {
    type: Array,
    default: () => []
  }
})

defineEmits(['open-note', 'open-author', 'login-request', 'like-change'])
</script>

<template>
  <section
    class="waterfall-feed"
    aria-label="推荐笔记"
  >
    <NoteCard
      v-for="note in notes"
      :key="note.id"
      :note="note"
      @open-author="$emit('open-author', $event)"
      @open="$emit('open-note', $event)"
      @login-request="$emit('login-request')"
      @like-change="$emit('like-change', $event)"
    />
  </section>
</template>

<style scoped>
.waterfall-feed {
  column-count: 5;
  column-gap: 28px;
}

@media (max-width: 1500px) {
  .waterfall-feed {
    column-count: 4;
    column-gap: 24px;
  }
}

@media (max-width: 1200px) {
  .waterfall-feed {
    column-count: 3;
  }
}

@media (max-width: 760px) {
  .waterfall-feed {
    column-count: 2;
    column-gap: 12px;
  }
}
</style>
