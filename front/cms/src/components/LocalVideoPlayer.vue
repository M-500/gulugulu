<template>
  <div class="local-video-player">
    <video
      ref="videoElement"
      :controls="controls"
      :muted="muted"
      :autoplay="autoplay"
      :preload="preload"
      playsinline
      @error="handleMediaError"
    />
    <p v-if="errorMessage">{{ errorMessage }}</p>
  </div>
</template>

<script setup>
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'

import { mountVideoSource } from '@/utils/video'

const props = defineProps({
  source: { type: Object, default: null },
  controls: { type: Boolean, default: false },
  muted: { type: Boolean, default: true },
  autoplay: { type: Boolean, default: false },
  preload: { type: String, default: 'metadata' }
})

const videoElement = ref(null)
const errorMessage = ref('')
let playback = null

onMounted(mountPlayer)
onBeforeUnmount(destroyPlayer)
watch(
  () => [props.source?.url, props.source?.name, props.source?.type],
  () => nextTick(mountPlayer)
)

defineExpose({
  getCurrentTime: () => videoElement.value?.currentTime || 0
})

function mountPlayer() {
  destroyPlayer()
  errorMessage.value = ''

  if (!props.source?.url || !videoElement.value) {
    return
  }

  try {
    playback = mountVideoSource(videoElement.value, props.source, {
      onError: showError
    })
  } catch (error) {
    showError(error)
  }
}

function destroyPlayer() {
  playback?.destroy()
  playback = null
}

function handleMediaError() {
  if (props.source?.url) {
    errorMessage.value = '视频预览加载失败'
  }
}

function showError(error) {
  errorMessage.value = error?.message || '视频预览加载失败'
}
</script>

<style scoped>
.local-video-player {
  position: relative;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

.local-video-player video {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.local-video-player p {
  position: absolute;
  right: 8px;
  bottom: 8px;
  left: 8px;
  margin: 0;
  border-radius: 6px;
  background: rgba(17, 24, 39, 0.72);
  color: #fff;
  font-size: 11px;
  padding: 6px;
  text-align: center;
}
</style>
