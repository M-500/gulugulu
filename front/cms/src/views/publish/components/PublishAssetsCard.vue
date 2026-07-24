<template>
  <section class="editor-card asset-card">
    <div class="editor-card__header">
      <h2>{{ type === 'imageText' ? '图文文件' : type === 'mixed' ? '混合作品素材' : '视频文件' }}</h2>
      <label class="reupload-button">重新上传<input :accept="accept" multiple type="file" @change="$emit('replace', $event)"></label>
    </div>
    <div class="asset-summary">
      <div v-for="asset in assets" :key="asset.assetId" class="asset-item">
        <span :class="['asset-kind', `asset-kind--${asset.kind}`]">{{ asset.kind === 'video' ? '视频' : '图片' }}</span>
        <span class="asset-name">{{ asset.name }}</span>
        <small>{{ formatSize(asset.size) }} · {{ asset.uploadStatus === 'uploaded' ? `已上传 #${asset.mediaId}` : '本地草稿' }}</small>
      </div>
    </div>
  </section>
</template>
<script setup>
defineProps({ type: { type: String, required: true }, accept: { type: String, required: true }, assets: { type: Array, default: () => [] } })
defineEmits(['replace'])
function formatSize(size) {
  if (size >= 1024 ** 3) return `${(size / 1024 ** 3).toFixed(2)}GB`
  if (size >= 1024 ** 2) return `${(size / 1024 ** 2).toFixed(1)}MB`
  return `${Math.max(1, Math.round(size / 1024))}KB`
}
</script>
