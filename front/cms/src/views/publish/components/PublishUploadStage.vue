<template>
  <div class="upload-stage">
    <div class="publish-tabs">
      <button v-for="item in publishTypes" :key="item.key" type="button"
        :class="{ 'is-active': modelValue === item.key }" @click="$emit('update:modelValue', item.key)">
        {{ item.label }}
      </button>
      <button class="draft-entry" type="button" @click="$emit('open-drafts')">草稿箱({{ draftCount }})</button>
    </div>
    <label class="upload-dropzone" :class="{ 'is-dragging': dragging }"
      @dragenter.prevent="$emit('update:dragging', true)" @dragover.prevent
      @dragleave.prevent="$emit('update:dragging', false)" @drop.prevent="$emit('drop', $event)">
      <input :accept="activeOption.accept" multiple type="file" @change="$emit('select-files', $event)">
      <span class="upload-icon">↑</span><strong>{{ activeOption.hint }}</strong>
      <button type="button">{{ activeOption.buttonText }}</button><small>{{ activeOption.description }}</small>
    </label>
    <div class="upload-tips">
      <div><strong>文件大小</strong><span>单个素材建议小于 4GB，最多支持 20 个文件</span></div>
      <div><strong>素材格式</strong><span>视频支持 mp4、mov，图片支持 jpg、png、webp</span></div>
      <div><strong>编辑体验</strong><span>上传后自动保存到本地草稿，可继续编辑发布</span></div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  modelValue: { type: String, required: true },
  publishTypes: { type: Array, required: true },
  activeOption: { type: Object, required: true },
  draftCount: { type: Number, default: 0 },
  dragging: { type: Boolean, default: false }
})
defineEmits(['update:modelValue', 'update:dragging', 'open-drafts', 'select-files', 'drop'])
</script>
