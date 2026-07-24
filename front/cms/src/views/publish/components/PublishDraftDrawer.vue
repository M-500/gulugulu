<template>
  <div class="drawer-mask" @click.self="$emit('close')"><aside class="draft-drawer">
    <div class="drawer-header"><h2>草稿箱</h2><button type="button" @click="$emit('close')">×</button></div>
    <div class="drawer-tabs"><button v-for="item in tabs" :key="item.key" type="button" :class="{ 'is-active': modelValue === item.key }" @click="$emit('update:modelValue', item.key)">{{ item.label }}({{ item.count }})</button></div>
    <p class="drawer-tip">草稿存储于当前浏览器本地，清除浏览器数据时会被删除；最多保存100篇草稿。</p>
    <div class="draft-list"><article v-for="draft in drafts" :key="draft.id" class="draft-item">
      <button class="draft-thumb" type="button" @click="$emit('edit', draft.id)"><span>{{ draft.type === 'video' ? '视频' : draft.type === 'mixed' ? '混合' : '图文' }}</span><small>{{ draft.assetMetas?.length || 0 }} 个素材</small></button>
      <button class="draft-info" type="button" @click="$emit('edit', draft.id)"><strong>{{ draft.title || '暂无笔记标题' }}</strong><span>保存于{{ draft.savedText }}</span></button>
      <button class="draft-action" type="button" @click="$emit('edit', draft.id)">编辑</button><button class="draft-action" type="button" @click="$emit('delete', draft.id)">删除</button>
    </article><p v-if="drafts.length === 0" class="empty-drafts">暂无草稿</p></div>
  </aside></div>
</template>
<script setup>
defineProps({ modelValue: { type: String, required: true }, tabs: { type: Array, required: true }, drafts: { type: Array, required: true } })
defineEmits(['update:modelValue', 'close', 'edit', 'delete'])
</script>
