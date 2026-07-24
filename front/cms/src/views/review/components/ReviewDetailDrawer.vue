<template>
  <div class="review-drawer" role="dialog" aria-modal="true" aria-label="作品审核详情" @click.self="$emit('close')">
    <aside class="review-drawer__panel">
      <header class="review-drawer__header">
        <div>
          <span>作品审核</span>
          <strong v-if="detail">#{{ detail.workId }}</strong>
        </div>
        <button type="button" aria-label="关闭" @click="$emit('close')">×</button>
      </header>

      <div v-if="loading" class="review-drawer__state">正在加载完整作品…</div>
      <div v-else-if="errorMessage" class="review-drawer__state is-error">
        <strong>加载失败</strong>
        <p>{{ errorMessage }}</p>
      </div>

      <template v-else-if="detail">
        <div class="review-preview">
          <ReviewVideoPlayer
            v-if="videoAsset"
            :playlist="detail.videoPlaylist"
            :poster="detail.coverUrl"
          />
          <div v-else class="review-preview__images">
            <img
              v-for="asset in imageAssets"
              :key="asset.mediaId"
              :src="asset.url"
              :alt="detail.title"
            >
          </div>
        </div>

        <div class="review-detail">
          <div class="review-detail__meta">
            <span>{{ detail.type === 'video' ? '视频作品' : `图片作品 · ${imageAssets.length} 张` }}</span>
            <span>{{ visibilityLabel }}</span>
            <span v-if="detail.original">原创声明</span>
          </div>
          <h2>{{ detail.title || '无标题作品' }}</h2>
          <p class="review-detail__content">{{ detail.content || '作者未填写作品正文。' }}</p>
          <div v-if="detail.topics.length" class="review-detail__topics">
            <span v-for="topic in detail.topics" :key="topic.topicId"># {{ topic.name }}</span>
          </div>

          <section class="review-author">
            <span>{{ (detail.authorName || '?').slice(0, 1) }}</span>
            <div>
              <strong>{{ detail.authorName }}</strong>
              <p>{{ detail.authorEmail }} · 用户 ID {{ detail.authorId }}</p>
            </div>
            <time>{{ formatDate(detail.submittedAt) }} 提交</time>
          </section>

          <section v-if="detail.reviewStatus !== 'pending_review'" class="review-result">
            <strong>{{ detail.reviewStatus === 'approved' ? '审核已通过' : '审核已驳回' }}</strong>
            <p v-if="detail.reviewReason">{{ detail.reviewReason }}</p>
            <time v-if="detail.reviewedAt">{{ formatDate(detail.reviewedAt) }}</time>
          </section>
        </div>

        <footer v-if="detail.reviewStatus === 'pending_review'" class="review-drawer__actions">
          <button class="is-reject" type="button" :disabled="submitting" @click="$emit('reject', detail)">驳回作品</button>
          <button class="is-approve" type="button" :disabled="submitting" @click="$emit('approve', detail)">
            {{ submitting ? '正在提交…' : '通过并发布' }}
          </button>
        </footer>
      </template>
    </aside>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import ReviewVideoPlayer from './ReviewVideoPlayer.vue'

const props = defineProps({
  detail: { type: Object, default: null },
  errorMessage: { type: String, default: '' },
  loading: { type: Boolean, default: false },
  submitting: { type: Boolean, default: false }
})

defineEmits(['approve', 'close', 'reject'])

const imageAssets = computed(() => (props.detail?.assets || []).filter((item) => item.role === 'image'))
const videoAsset = computed(() => (props.detail?.assets || []).find((item) => item.role === 'video'))
const visibilityLabel = computed(() => ({
  public: '公开可见',
  private: '仅自己可见',
  mutual: '仅互关可见',
  selected: '指定用户可见',
  excluded: '部分用户不可见'
}[props.detail?.visibility] || props.detail?.visibility))

function formatDate(value) {
  if (!value) return ''
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  }).format(new Date(value)).replace(/\//g, '-')
}
</script>
