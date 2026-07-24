<template>
  <section class="editor-card cover-card">
    <div class="editor-card__header">
      <div><h2>设置封面</h2><p>{{ type === 'video' ? '视频作品需要上传一张 JPG 或 PNG 封面。' : '默认使用第一张图片，也可以切换或上传单独封面。' }}</p></div>
      <label class="cover-upload-button">上传封面<input accept="image/jpeg,image/png" type="file" @change="$emit('upload-cover', $event)"></label>
    </div>
    <div class="cover-list">
      <button v-if="customCoverUrl" type="button" class="is-active">
        <img :src="customCoverUrl" alt="自定义封面"><span>自定义封面</span>
      </button>
      <button v-for="asset in selectableAssets" :key="asset.assetId" type="button"
        :class="{ 'is-active': coverAssetId === asset.assetId }" @click="$emit('select', asset.assetId)">
        <img :src="asset.url" :alt="asset.name"><span>素材封面</span>
      </button>
    </div>
  </section>
</template>
<script setup>
import { computed } from 'vue'
const props = defineProps({
  type: { type: String, required: true },
  assets: { type: Array, default: () => [] },
  coverAssetId: { type: [String, Number], default: null },
  customCoverUrl: { type: String, default: '' }
})
defineEmits(['select', 'upload-cover'])
const selectableAssets = computed(() => props.type === 'image'
  ? props.assets.filter((asset) => asset.kind === 'image')
  : [])
</script>
