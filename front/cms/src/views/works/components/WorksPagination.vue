<template>
  <footer class="works-pagination">
    <span>共 {{ total }} 个作品</span>
    <div>
      <button type="button" :disabled="page <= 1" @click="$emit('change', page - 1)">上一页</button>
      <strong>{{ page }} / {{ pageCount }}</strong>
      <button type="button" :disabled="page >= pageCount" @click="$emit('change', page + 1)">下一页</button>
    </div>
  </footer>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  page: { type: Number, required: true },
  pageSize: { type: Number, required: true },
  total: { type: Number, required: true }
})

defineEmits(['change'])

const pageCount = computed(() => Math.max(1, Math.ceil(props.total / props.pageSize)))
</script>
