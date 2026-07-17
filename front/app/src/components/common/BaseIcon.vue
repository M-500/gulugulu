<script setup>
const props = defineProps({
  name: {
    type: String,
    required: true
  },
  size: {
    type: [Number, String],
    default: 20
  },
  strokeWidth: {
    type: [Number, String],
    default: 2
  },
  fill: {
    type: String,
    default: 'none'
  }
})

const paths = {
  home: ['M3 10.5 12 3l9 7.5', 'M5 9.5V21h14V9.5', 'M9 21v-6h6v6'],
  ball: ['M12 22a10 10 0 1 0 0-20 10 10 0 0 0 0 20Z', 'M8 5.2 12 8l4-2.8', 'M8.5 18.8 10 14l-4-3', 'M15.5 18.8 14 14l4-3', 'M10 14h4l1.5-4L12 8l-3.5 2Z'],
  message: ['M21 11.5a8.4 8.4 0 0 1-8.8 8.4 9 9 0 0 1-3.8-.8L3 21l1.8-5A8.2 8.2 0 0 1 3 11.5a8.5 8.5 0 0 1 18 0Z'],
  play: ['M8 5v14l11-7Z'],
  video: ['M4 7h10a3 3 0 0 1 3 3v4a3 3 0 0 1-3 3H4Z', 'M17 10l4-2v8l-4-2Z'],
  plus: ['M12 5v14', 'M5 12h14'],
  bell: ['M18 8a6 6 0 0 0-12 0c0 7-3 7-3 9h18c0-2-3-2-3-9Z', 'M10 21h4'],
  user: ['M20 21a8 8 0 0 0-16 0', 'M12 13a4 4 0 1 0 0-8 4 4 0 0 0 0 8Z'],
  menu: ['M4 7h16', 'M4 12h16', 'M4 17h16'],
  info: ['M12 22a10 10 0 1 0 0-20 10 10 0 0 0 0 20Z', 'M12 10v7', 'M12 7h.01'],
  search: ['M11 19a8 8 0 1 0 0-16 8 8 0 0 0 0 16Z', 'm21 21-4.35-4.35'],
  heart: ['M20.8 4.6a5.5 5.5 0 0 0-7.8 0L12 5.7l-1-1.1a5.5 5.5 0 0 0-7.8 7.8l1 1L12 21l7.8-7.6 1-1a5.5 5.5 0 0 0 0-7.8Z'],
  up: ['M18 15 12 9l-6 6'],
  refresh: ['M21 12a9 9 0 0 1-15.4 6.3L3 15', 'M3 21v-6h6', 'M3 12A9 9 0 0 1 18.4 5.7L21 9', 'M21 3v6h-6']
}
</script>

<template>
  <svg
    class="base-icon"
    :width="props.size"
    :height="props.size"
    viewBox="0 0 24 24"
    :fill="props.fill"
    stroke="currentColor"
    :stroke-width="props.strokeWidth"
    stroke-linecap="round"
    stroke-linejoin="round"
    aria-hidden="true"
  >
    <path
      v-for="path in paths[props.name]"
      :key="path"
      :d="path"
    />
  </svg>
</template>
