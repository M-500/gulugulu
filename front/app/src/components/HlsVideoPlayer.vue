<script setup>
import Hls from 'hls.js/dist/hls.light.mjs'
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'

const props = defineProps({
  playlist: {
    type: String,
    default: ''
  },
  poster: {
    type: String,
    default: ''
  }
})

const videoElement = ref(null)
const playbackRate = ref(1)
const showRateMenu = ref(false)
const errorMessage = ref('')
const playbackRates = [0.5, 0.75, 1, 1.25, 1.5, 2]
let hls = null

onMounted(mountPlayer)
onBeforeUnmount(destroyPlayer)

watch(
  () => props.playlist,
  () => nextTick(mountPlayer)
)

function mountPlayer() {
  destroyPlayer()
  errorMessage.value = ''
  if (!props.playlist || !videoElement.value) {
    errorMessage.value = '视频播放清单暂不可用'
    return
  }

  const playlistUrl = new URL(props.playlist, window.location.origin).toString()
  if (Hls.isSupported()) {
    hls = new Hls()
    hls.loadSource(playlistUrl)
    hls.attachMedia(videoElement.value)
    hls.on(Hls.Events.ERROR, (_, data) => {
      if (data.fatal) errorMessage.value = '视频加载失败，请稍后重试'
    })
    return
  }
  if (videoElement.value.canPlayType('application/vnd.apple.mpegurl')) {
    videoElement.value.src = playlistUrl
    return
  }
  errorMessage.value = '当前浏览器不支持 HLS 视频播放'
}

function setPlaybackRate(rate) {
  playbackRate.value = rate
  showRateMenu.value = false
  if (videoElement.value) {
    videoElement.value.playbackRate = rate
  }
}

function destroyPlayer() {
  if (hls) {
    hls.destroy()
    hls = null
  }
  if (videoElement.value) {
    videoElement.value.removeAttribute('src')
    videoElement.value.load()
  }
}
</script>

<template>
  <div class="hls-player">
    <video
      ref="videoElement"
      :poster="poster"
      controls
      playsinline
      preload="metadata"
    />

    <div class="hls-player__rate">
      <button
        type="button"
        :aria-expanded="showRateMenu"
        aria-label="选择播放倍速"
        @click="showRateMenu = !showRateMenu"
      >
        {{ playbackRate === 1 ? '倍速' : `${playbackRate}x` }}
      </button>
      <div
        v-if="showRateMenu"
        class="hls-player__rate-menu"
      >
        <button
          v-for="rate in playbackRates"
          :key="rate"
          type="button"
          :class="{ active: playbackRate === rate }"
          @click="setPlaybackRate(rate)"
        >
          {{ rate }}x
        </button>
      </div>
    </div>

    <p v-if="errorMessage">
      {{ errorMessage }}
    </p>
  </div>
</template>

<style scoped>
.hls-player {
  position: absolute;
  inset: 0;
  display: grid;
  place-items: center;
  overflow: hidden;
  background: #000;
}

.hls-player video {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.hls-player__rate {
  position: absolute;
  right: 14px;
  bottom: 54px;
  z-index: 3;
}

.hls-player__rate > button {
  min-width: 48px;
  height: 30px;
  border-radius: 6px;
  color: #fff;
  background: rgba(20, 20, 22, 0.72);
  font-size: 13px;
  padding: 0 10px;
}

.hls-player__rate-menu {
  position: absolute;
  right: 0;
  bottom: 36px;
  display: grid;
  width: 86px;
  overflow: hidden;
  border-radius: 8px;
  background: rgba(20, 20, 22, 0.92);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.3);
}

.hls-player__rate-menu button {
  height: 34px;
  color: #fff;
  font-size: 13px;
}

.hls-player__rate-menu button:hover,
.hls-player__rate-menu button.active {
  color: var(--color-primary);
  background: rgba(255, 255, 255, 0.08);
}

.hls-player > p {
  position: absolute;
  right: 20px;
  bottom: 92px;
  left: 20px;
  margin: 0;
  color: #fff;
  font-size: 13px;
  text-align: center;
}
</style>
