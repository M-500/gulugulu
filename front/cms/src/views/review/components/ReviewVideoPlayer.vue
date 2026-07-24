<template>
  <div class="review-video">
    <video ref="videoElement" :poster="poster" controls playsinline />
    <p v-if="errorMessage">{{ errorMessage }}</p>
  </div>
</template>

<script setup>
import Hls from 'hls.js/dist/hls.light.mjs'
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'

const props = defineProps({
  playlist: { type: String, default: '' },
  poster: { type: String, default: '' }
})

const videoElement = ref(null)
const errorMessage = ref('')
let hls = null
let playlistUrl = ''

onMounted(mountPlayer)
onBeforeUnmount(destroyPlayer)
watch(() => props.playlist, () => nextTick(mountPlayer))

function mountPlayer() {
  destroyPlayer()
  errorMessage.value = ''
  if (!props.playlist || !videoElement.value) {
    errorMessage.value = '视频播放清单暂不可用'
    return
  }

  playlistUrl = URL.createObjectURL(new Blob([props.playlist], {
    type: 'application/vnd.apple.mpegurl'
  }))
  if (Hls.isSupported()) {
    hls = new Hls()
    hls.loadSource(playlistUrl)
    hls.attachMedia(videoElement.value)
    hls.on(Hls.Events.ERROR, (_, data) => {
      if (data.fatal) errorMessage.value = '视频加载失败，请刷新后重试'
    })
    return
  }
  if (videoElement.value.canPlayType('application/vnd.apple.mpegurl')) {
    videoElement.value.src = playlistUrl
    return
  }
  errorMessage.value = '当前浏览器不支持 HLS 视频播放'
}

function destroyPlayer() {
  if (hls) {
    hls.destroy()
    hls = null
  }
  if (playlistUrl) {
    URL.revokeObjectURL(playlistUrl)
    playlistUrl = ''
  }
}
</script>
