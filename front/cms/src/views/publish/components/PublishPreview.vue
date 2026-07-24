<template>
  <aside class="preview-panel">
    <div class="preview-tabs"><button class="is-active" type="button">笔记预览</button><button type="button">封面预览</button></div>
    <div class="phone-frame"><div class="phone-screen">
      <div class="phone-status">9:41 <span>▮▮ ◒ ▰</span></div>
      <div class="phone-media"><video v-if="asset?.kind === 'video'" :src="asset.url" controls muted playsinline /><img v-else-if="asset" :src="asset.url" alt=""><div v-else class="empty-preview">暂无素材</div></div>
      <div class="phone-copy"><h3>{{ form.title || '项目工程' }}</h3><p>{{ form.body || '提示词 提示词 提示词' }}</p><div class="phone-user"><span>咕</span><strong>咕噜创作者</strong><button type="button">关注</button></div></div>
      <div class="phone-actions"><span>说点什么...</span><strong>♡ 点赞</strong><strong>☆ 收藏</strong><strong>☏ 评论</strong></div>
    </div></div>
  </aside>
</template>
<script setup>
defineProps({ asset: { type: Object, default: null }, form: { type: Object, required: true } })
</script>
