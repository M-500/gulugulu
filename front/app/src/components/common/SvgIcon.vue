<!-- eslint-disable vue/no-v-html -->
<template>
  <span
    class="svg-icon"
    :class="{ 'svg-icon--missing': !iconContent }"
    :style="sizeStyle"
    :title="title || undefined"
    :role="title ? 'img' : undefined"
    :aria-label="title || undefined"
    :aria-hidden="title ? undefined : 'true'"
    v-html="iconContent"
  />
</template>

<script setup>
import { computed } from 'vue'

// 扫描用户指定的 SVG 目录，新增文件后只需通过文件名引用，无需逐个 import。
const iconModules = import.meta.glob('../../src/assets/icons/**/*.svg', {
  eager: true,
  query: '?raw',
  import: 'default'
})
const icons = Object.fromEntries(
  Object.entries(iconModules).map(([path, source]) => [
    path.split('/icons/')[1].replace(/\.svg$/i, ''),
    normalizeSvg(source)
  ])
)

const props = defineProps({
  name: {
    type: String,
    required: true
  },
  size: {
    type: [Number, String],
    default: '1em'
  },
  width: {
    type: [Number, String],
    default: undefined
  },
  height: {
    type: [Number, String],
    default: undefined
  },
  title: {
    type: String,
    default: ''
  }
})

const normalizedName = computed(() => String(props.name || '')
  .trim()
  .replace(/\\/g, '/')
  .replace(/^\//, '')
  .replace(/\.svg$/i, ''))
const iconContent = computed(() => icons[normalizedName.value] || '')
const sizeStyle = computed(() => ({
  width: formatSize(props.width ?? props.size),
  height: formatSize(props.height ?? props.size)
}))

// 图标均来自项目内固定目录，可以安全内联；同时移除固定尺寸并继承外层文字颜色。
function normalizeSvg(source) {
  return String(source || '')
    .replace(/<\?xml[^>]*>/gi, '')
    .replace(/<!DOCTYPE[^>]*(?:\[[\s\S]*?\]\s*)?>/gi, '')
    .replace(/\s(?:width|height)="[^"]*"/gi, '')
    .replace(/(fill|stroke)="(?!none|currentColor)[^"]*"/gi, '$1="currentColor"')
    .replace('<svg', '<svg aria-hidden="true" focusable="false"')
}

function formatSize(value) {
  if (typeof value === 'number') return `${value}px`
  const normalized = String(value || '').trim()
  return /^\d+(?:\.\d+)?$/.test(normalized) ? `${normalized}px` : normalized
}
</script>

<style scoped>
.svg-icon {
  display: inline-flex;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  vertical-align: -0.125em;
}

.svg-icon :deep(svg) {
  display: block;
  width: 100%;
  height: 100%;
  overflow: visible;
}
</style>
