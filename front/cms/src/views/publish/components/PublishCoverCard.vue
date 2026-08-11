<template>
  <section class="editor-card cover-card">
    <div class="editor-card__header">
      <div><h2>设置封面</h2><p>{{ type === 'video' ? '视频封面只能从视频帧中截取，默认使用第一帧。' : '图片封面只能从已上传图片中选择，默认使用第一张。' }}</p></div>
    </div>

    <div v-if="type === 'video'" class="video-cover-tool">
      <div class="video-cover-source">
        <LocalVideoPlayer v-if="videoAsset" ref="videoRef" :source="videoAsset" controls />
        <div v-else class="cover-empty">暂无视频素材</div>
      </div>
      <div class="video-cover-result">
        <img v-if="customCoverUrl" :src="customCoverUrl" alt="视频封面">
        <div v-else class="cover-empty">默认使用第一帧</div>
      </div>
      <div class="video-cover-actions">
        <button type="button" :disabled="!videoAsset" @click="captureCurrentFrame">截取当前帧</button>
        <button type="button" :disabled="!videoAsset" @click="$emit('capture-frame', 0)">使用第一帧</button>
      </div>
    </div>

    <div v-else class="cover-list">
      <button v-for="asset in selectableAssets" :key="asset.assetId" type="button"
        :class="{ 'is-active': coverAssetId === asset.assetId }" @click="$emit('select', asset.assetId)">
        <img :src="asset.url" :alt="asset.name"><span>{{ coverAssetId === asset.assetId ? '当前封面' : '选择封面' }}</span>
      </button>
      <div v-if="!selectableAssets.length" class="cover-empty">暂无可选图片</div>
    </div>
  </section>
</template>
<script setup>
import { computed, ref } from 'vue'

import LocalVideoPlayer from '@/components/LocalVideoPlayer.vue'

const props = defineProps({
  type: { type: String, required: true },
  assets: { type: Array, default: () => [] },
  coverAssetId: { type: [String, Number], default: null },
  customCoverUrl: { type: String, default: '' }
})
const emit = defineEmits(['select', 'capture-frame'])
const videoRef = ref(null)
const selectableAssets = computed(() => props.type === 'image'
  ? props.assets.filter((asset) => asset.kind === 'image')
  : [])
const videoAsset = computed(() => props.assets.find((asset) => asset.kind === 'video'))

function captureCurrentFrame() {
  emit('capture-frame', videoRef.value?.getCurrentTime() || 0)
}
</script>
