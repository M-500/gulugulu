<template>
  <span
    class="svg-icon"
    :class="{ 'svg-icon--missing': !iconUrl }"
    :style="containerStyle"
    :title="title || undefined"
    :role="accessibleLabel ? 'img' : undefined"
    :aria-label="accessibleLabel || undefined"
    :aria-hidden="accessibleLabel ? undefined : 'true'"
  >
    <span
      v-if="iconUrl && color"
      class="svg-icon__mask"
      :style="maskStyle"
    />
    <img
      v-else-if="iconUrl"
      class="svg-icon__image"
      :src="iconUrl"
      alt=""
    >
  </span>
</template>

<script setup>
import { computed } from 'vue'

// Webpack 会在构建时收集该目录下的 SVG，新增文件后可直接通过文件名引用。
const iconModules = require.context('@/assets/icons', true, /\.svg$/i)

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
  color: {
    type: String,
    default: ''
  },
  title: {
    type: String,
    default: ''
  }
})

const normalizedName = computed(() => normalizeName(props.name))
const accessibleLabel = computed(() => props.title.trim())
const iconUrl = computed(() => {
  if (!normalizedName.value) return ''
  const key = `./${normalizedName.value}.svg`
  if (!iconModules.keys().includes(key)) return ''
  const source = iconModules(key)
  return typeof source === 'string' ? source : source?.default || ''
})

const containerStyle = computed(() => ({
  width: formatSize(props.width ?? props.size),
  height: formatSize(props.height ?? props.size),
  color: props.color || undefined
}))

const maskStyle = computed(() => ({
  maskImage: `url("${iconUrl.value}")`,
  WebkitMaskImage: `url("${iconUrl.value}")`
}))

function normalizeName(value) {
  const normalized = String(value || '')
    .trim()
    .replace(/\\/g, '/')
    .replace(/^\.\//, '')
    .replace(/^\//, '')
    .replace(/\.svg$/i, '')
  if (!normalized || normalized.split('/').includes('..')) return ''
  return normalized
}

function formatSize(value) {
  return typeof value === 'number' ? `${value}px` : value
}
</script>

<style scoped>
.svg-icon {
  display: inline-flex;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  line-height: 1;
  vertical-align: -0.125em;
}

.svg-icon__image,
.svg-icon__mask {
  display: block;
  width: 100%;
  height: 100%;
}

.svg-icon__image {
  object-fit: contain;
}

.svg-icon__mask {
  background-color: currentColor;
  mask-position: center;
  mask-repeat: no-repeat;
  mask-size: contain;
  -webkit-mask-position: center;
  -webkit-mask-repeat: no-repeat;
  -webkit-mask-size: contain;
}
</style>
