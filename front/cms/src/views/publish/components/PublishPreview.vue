<template>
  <aside class="preview-panel">
    <div class="preview-tabs">
      <button type="button" :class="{ 'is-active': activeTab === 'note' }" @click="activeTab = 'note'">笔记预览</button>
      <button type="button" :class="{ 'is-active': activeTab === 'cover' }" @click="activeTab = 'cover'">封面预览</button>
    </div>

    <div class="phone-frame">
      <div class="phone-screen" :class="{ 'is-cover-preview': activeTab === 'cover' }">
        <template v-if="activeTab === 'note'">
          <div class="phone-status">9:41 <span>▮▮ ◒ ▰</span></div>
          <div class="phone-media">
            <video v-if="asset?.kind === 'video'" :src="asset.url" controls muted playsinline />
            <img v-else-if="asset" :src="asset.url" alt="">
            <div v-else class="empty-preview">暂无素材</div>
          </div>
          <div class="phone-copy">
            <h3>{{ form.title || '项目工程' }}</h3><p>{{ form.body || '提示词 提示词 提示词' }}</p>
            <div class="phone-user"><span>咕</span><strong>咕噜创作者</strong><button type="button">关注</button></div>
          </div>
          <div class="phone-actions"><span>说点什么...</span><strong>♡ 点赞</strong><strong>☆ 收藏</strong><strong>☏ 评论</strong></div>
        </template>

        <template v-else>
          <header class="discovery-header">
            <div class="discovery-status"><strong>9:41</strong><span>▮▮ ◉ ▰</span></div>
            <div class="discovery-nav"><span>☰</span><span>关注</span><strong>发现</strong><span>附近</span><span>⌕</span></div>
            <div class="discovery-categories"><strong>推荐</strong><span>直播</span><span>短剧</span><span>穿搭</span><span>旅行</span><span>动漫</span><span>⌄</span></div>
          </header>
          <div class="cover-feed">
            <article v-for="(card, index) in coverCards" :key="card.key" class="cover-note-card">
              <div class="cover-note-card__media" :class="`is-card-${index + 1}`">
                <video v-if="card.asset?.kind === 'video'" :src="card.asset.url" muted playsinline />
                <img v-else-if="card.asset" :src="card.asset.url" alt="">
                <div v-else class="cover-placeholder"><span>咕噜</span><small>记录生活灵感</small></div>
                <i v-if="card.asset?.kind === 'video'">▶</i>
              </div>
              <h4>{{ card.title }}</h4>
              <div class="cover-note-card__meta"><span><b>咕</b> 咕噜创作者</span><span>♡ 0</span></div>
            </article>
          </div>
          <nav class="discovery-footer"><strong>首页</strong><span>市集</span><button type="button">＋</button><span>消息</span><span>我</span></nav>
        </template>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { computed, ref } from 'vue'

const props = defineProps({
  asset: { type: Object, default: null },
  assets: { type: Array, default: () => [] },
  form: { type: Object, required: true }
})

const activeTab = ref('note')
const coverCards = computed(() => Array.from({ length: 4 }, (_, index) => ({
  key: props.assets[index]?.assetId || `placeholder-${index}`,
  asset: props.assets[index] || (index === 0 ? props.asset : null),
  title: index === 0 ? (props.form.title || '示例笔记标题') : `示例笔记标题${index + 1}`
})))
</script>
