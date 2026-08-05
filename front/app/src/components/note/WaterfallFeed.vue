<script setup>
import { nextTick, onBeforeUnmount, onBeforeUpdate, onMounted, onUpdated, ref, watch } from 'vue'

import NoteCard from '@/components/note/NoteCard.vue'

const props = defineProps({
  notes: {
    type: Array,
    default: () => []
  }
})

defineEmits(['open-note', 'open-author', 'login-request', 'like-change'])

const container = ref(null)
const itemElements = ref([])
let itemResizeObserver
let layoutFrame = 0

function setItemElement(element, index) {
  if (element) itemElements.value[index] = element
}

function getLayoutOptions() {
  const viewportWidth = window.innerWidth
  if (viewportWidth <= 760) return { columns: 2, columnGap: 12, rowGap: 18 }
  if (viewportWidth <= 1200) return { columns: 3, columnGap: 24, rowGap: 24 }
  if (viewportWidth <= 1500) return { columns: 4, columnGap: 24, rowGap: 24 }
  return { columns: 5, columnGap: 28, rowGap: 24 }
}

function scheduleLayout() {
  window.cancelAnimationFrame(layoutFrame)
  layoutFrame = window.requestAnimationFrame(layoutItems)
}

// 按数据顺序逐个放入当前最短列：第一排自然从左到右，后续内容自动向上补位。
function layoutItems() {
  const root = container.value
  const elements = itemElements.value.filter(Boolean)
  if (!root || !elements.length) {
    if (root) root.style.height = '0px'
    return
  }

  const { columns, columnGap, rowGap } = getLayoutOptions()
  const itemWidth = Math.max(0, (root.clientWidth - columnGap * (columns - 1)) / columns)
  const columnHeights = Array(columns).fill(0)

  for (const element of elements) element.style.width = `${itemWidth}px`

  for (const element of elements) {
    const shortestHeight = Math.min(...columnHeights)
    const columnIndex = columnHeights.indexOf(shortestHeight)
    const x = columnIndex * (itemWidth + columnGap)
    const y = columnHeights[columnIndex]
    element.style.transform = `translate3d(${x}px, ${y}px, 0)`
    columnHeights[columnIndex] = y + element.offsetHeight + rowGap
  }

  root.style.height = `${Math.max(...columnHeights) - rowGap}px`
}

function observeItems() {
  itemResizeObserver?.disconnect()
  itemResizeObserver = new ResizeObserver(scheduleLayout)
  for (const element of itemElements.value) {
    if (element) itemResizeObserver.observe(element)
  }
  scheduleLayout()
}

onBeforeUpdate(() => {
  itemElements.value = []
})

onUpdated(() => {
  nextTick(observeItems)
})

onMounted(() => {
  window.addEventListener('resize', scheduleLayout)
  nextTick(observeItems)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', scheduleLayout)
  window.cancelAnimationFrame(layoutFrame)
  itemResizeObserver?.disconnect()
})

watch(() => props.notes.length, () => nextTick(observeItems))
</script>

<template>
  <section
    ref="container"
    class="waterfall-feed"
    aria-label="推荐笔记"
  >
    <div
      v-for="(note, index) in notes"
      :key="note.id"
      :ref="(element) => setItemElement(element, index)"
      class="waterfall-feed__item"
    >
      <NoteCard
        :note="note"
        @open-author="$emit('open-author', $event)"
        @open="$emit('open-note', $event)"
        @login-request="$emit('login-request')"
        @like-change="$emit('like-change', $event)"
      />
    </div>
  </section>
</template>

<style scoped>
.waterfall-feed {
  position: relative;
  width: 100%;
  transition: height 0.2s ease;
}

.waterfall-feed__item {
  position: absolute;
  top: 0;
  left: 0;
  will-change: transform;
  transition: transform 0.22s ease, width 0.22s ease;
}

.waterfall-feed__item :deep(.note-card) {
  display: block;
  margin: 0;
}

@media (prefers-reduced-motion: reduce) {
  .waterfall-feed,
  .waterfall-feed__item {
    transition: none;
  }
}
</style>
