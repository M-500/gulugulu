<template>
  <section class="editor-card cover-card">
    <div class="editor-card__header">
      <div><h2>设置封面</h2><p>默认取第一份素材作为封面，可以切换推荐素材。</p></div>
      <label class="switch-row"><span>PK 封面</span><input v-model="pkCover" type="checkbox"></label>
    </div>
    <div class="cover-list">
      <button v-for="asset in assets" :key="asset.assetId" type="button"
        :class="{ 'is-active': coverAssetId === asset.assetId }" @click="$emit('select', asset.assetId)">
        <img v-if="asset.kind === 'image'" :src="asset.url" alt=""><video v-else :src="asset.url" muted playsinline /><span>推荐封面</span>
      </button>
    </div>
  </section>
</template>
<script setup>
import { ref } from 'vue'
defineProps({ assets: { type: Array, default: () => [] }, coverAssetId: { type: [String, Number], default: null } })
defineEmits(['select'])
const pkCover = ref(false)
</script>
