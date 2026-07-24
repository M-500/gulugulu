<template>
  <section class="editor-card asset-card">
    <template v-if="type === 'imageText'">
      <div class="editor-card__header image-editor__header">
        <div><h2>图片编辑</h2><span>{{ assets.length }}/19</span></div>
        <button class="cover-advice" type="button">◕ 获取封面建议</button>
      </div>
      <div class="image-editor__grid">
        <label v-if="assets.length < 19" class="image-editor__add">
          <input accept="image/*" multiple type="file" @change="$emit('add-images', $event)">
          <span>＋</span><small>添加图片</small>
        </label>
        <div v-for="(asset, index) in assets" :key="asset.assetId" class="image-editor__item">
          <img :src="asset.url" :alt="asset.name">
          <span>{{ index + 1 }}</span>
          <button type="button" aria-label="删除图片" @click="$emit('delete-image', asset.assetId)">×</button>
        </div>
      </div>
      <p v-if="assets.length >= 19" class="image-editor__limit">已达到 19 张图片上限</p>
    </template>

    <template v-else>
      <div class="editor-card__header">
        <h2>{{ type === 'mixed' ? '混合作品素材' : '视频文件' }}</h2>
        <label class="reupload-button">重新上传<input :accept="accept" multiple type="file" @change="$emit('replace', $event)"></label>
      </div>
      <div class="asset-summary">
        <div v-for="asset in assets" :key="asset.assetId" class="asset-item">
          <span :class="['asset-kind', `asset-kind--${asset.kind}`]">{{ asset.kind === 'video' ? '视频' : '图片' }}</span>
          <span class="asset-name">{{ asset.name }}</span>
          <small>{{ formatSize(asset.size) }} · {{ asset.uploadStatus === 'uploaded' ? `已上传 #${asset.mediaId}` : '本地草稿' }}</small>
        </div>
      </div>
    </template>
  </section>
</template>

<script setup>
defineProps({ type: { type: String, required: true }, accept: { type: String, required: true }, assets: { type: Array, default: () => [] } })
defineEmits(['replace', 'add-images', 'delete-image'])
function formatSize(size) {
  if (size >= 1024 ** 3) return `${(size / 1024 ** 3).toFixed(2)}GB`
  if (size >= 1024 ** 2) return `${(size / 1024 ** 2).toFixed(1)}MB`
  return `${Math.max(1, Math.round(size / 1024))}KB`
}
</script>
