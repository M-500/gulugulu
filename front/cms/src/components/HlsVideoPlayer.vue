<template>
  <div class="hls-player">
    <video ref="videoElement" :poster="poster" controls playsinline />
    <p v-if="errorMessage">{{ errorMessage }}</p>
  </div>
</template>

<script setup>
import Hls from 'hls.js/dist/hls.light.mjs'
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'

import { useAuthStore } from '@/stores/auth'

const props = defineProps({
  playlist: { type: String, default: '' },
  poster: { type: String, default: '' }
})

const authStore = useAuthStore()
const videoElement = ref(null)
const errorMessage = ref('')
let hls = null

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

  const playlistUrl = resolvePlaylistUrl(props.playlist)
  if (Hls.isSupported()) {
    hls = new Hls({
      xhrSetup(xhr, url) {
        if (shouldAttachAuth(url) && authStore.accessToken) {
          xhr.setRequestHeader('Authorization', `Bearer ${authStore.accessToken}`)
        }
      }
    })
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
}

function resolvePlaylistUrl(value) {
  return new URL(value, window.location.origin).toString()
}

function shouldAttachAuth(value) {
  const url = new URL(value, window.location.origin)
  return url.origin === window.location.origin && url.pathname.startsWith('/api/')
}
</script>

<style scoped>
.hls-player {
  position: relative;
  width: 100%;
}

.hls-player video {
  display: block;
  width: 100%;
  max-height: 100%;
}

.hls-player p {
  position: absolute;
  right: 16px;
  bottom: 12px;
  left: 16px;
  margin: 0;
  color: #fff;
  font-size: 11px;
  text-align: center;
}
</style>
