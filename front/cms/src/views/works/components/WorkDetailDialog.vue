<template>
  <div class="work-detail-dialog" role="dialog" aria-modal="true" aria-label="作品详情" @click.self="$emit('close')">
    <section class="work-detail-dialog__panel">
      <header class="work-detail-dialog__header">
        <div>
          <span class="work-detail-dialog__eyebrow">WORK PREVIEW</span>
          <strong>作品详情</strong>
          <small v-if="detail">#{{ detail.workId }}</small>
        </div>
        <button type="button" aria-label="关闭作品详情" @click="$emit('close')">×</button>
      </header>

      <div v-if="loading" class="work-detail-dialog__state">
        <span class="works-state-card__spinner" />
        <strong>正在加载作品详情</strong>
        <p>正在获取图片与视频预览资源…</p>
      </div>

      <div v-else-if="errorMessage" class="work-detail-dialog__state is-error">
        <span>!</span>
        <strong>作品详情加载失败</strong>
        <p>{{ errorMessage }}</p>
        <button type="button" @click="$emit('retry')">重新加载</button>
      </div>

      <template v-else-if="detail">
        <div class="work-detail-dialog__content">
          <section class="work-detail-media" :class="{ 'is-video': detail.type === 'video' }">
            <HlsVideoPlayer
              v-if="detail.type === 'video'"
              :playlist="detail.videoPlaylist"
              :poster="detail.coverUrl"
            />

            <template v-else>
              <div class="work-detail-media__stage">
                <img
                  v-if="activeImage"
                  :src="activeImage.url"
                  :alt="detail.title"
                >
                <div v-else class="work-detail-media__empty">暂无可预览图片</div>
                <span v-if="imageAssets.length">{{ activeImageIndex + 1 }} / {{ imageAssets.length }}</span>
              </div>
              <div v-if="imageAssets.length > 1" class="work-detail-media__thumbs">
                <button
                  v-for="(asset, index) in imageAssets"
                  :key="asset.mediaId"
                  type="button"
                  :class="{ 'is-active': index === activeImageIndex }"
                  @click="activeImageIndex = index"
                >
                  <img :src="asset.url" :alt="`第 ${index + 1} 张图片`">
                </button>
              </div>
            </template>
          </section>

          <aside class="work-detail-info">
            <div class="work-detail-info__badges">
              <span :class="['is-status', `is-${statusMeta.tone}`]">{{ statusMeta.label }}</span>
              <span>{{ detail.type === 'video' ? '视频作品' : `图片作品 · ${imageAssets.length} 张` }}</span>
              <span>{{ visibilityLabel }}</span>
              <span v-if="detail.original">原创</span>
            </div>

            <h2>{{ detail.title || '无标题作品' }}</h2>
            <p class="work-detail-info__content">{{ detail.content || '作者未填写作品正文。' }}</p>

            <div v-if="detail.topics?.length" class="work-detail-info__topics">
              <span v-for="topic in detail.topics" :key="topic.topicId"># {{ topic.name }}</span>
            </div>

            <dl class="work-detail-info__meta">
              <div>
                <dt>作品编号</dt>
                <dd>#{{ detail.workId }}</dd>
              </div>
              <div>
                <dt>创建时间</dt>
                <dd>{{ formatDate(detail.createdAt) }}</dd>
              </div>
              <div v-if="detail.type === 'video'">
                <dt>视频时长</dt>
                <dd>{{ formatDuration(detail.durationMs) }}</dd>
              </div>
              <div v-if="detail.scheduledAt">
                <dt>计划发布</dt>
                <dd>{{ formatDate(detail.scheduledAt) }}</dd>
              </div>
              <div v-if="detail.publishedAt">
                <dt>实际发布</dt>
                <dd>{{ formatDate(detail.publishedAt) }}</dd>
              </div>
            </dl>

            <section v-if="detail.reviewStatus === 'rejected'" class="work-detail-info__review">
              <span>审核反馈</span>
              <p>{{ detail.reviewReason || '作品未通过审核，请修改后重新提交。' }}</p>
            </section>

            <footer>
              <div>
                <span>媒体资源</span>
                <strong>{{ mediaSummary }}</strong>
              </div>
              <button type="button" @click="$emit('close')">关闭预览</button>
            </footer>
          </aside>
        </div>
      </template>
    </section>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'

import HlsVideoPlayer from '@/components/HlsVideoPlayer.vue'

const props = defineProps({
  detail: { type: Object, default: null },
  errorMessage: { type: String, default: '' },
  loading: { type: Boolean, default: false }
})

const emit = defineEmits(['close', 'retry'])
const activeImageIndex = ref(0)

const imageAssets = computed(() => (props.detail?.assets || [])
  .filter((item) => item.role === 'image' && item.url)
  .sort((left, right) => left.sort - right.sort))
const activeImage = computed(() => imageAssets.value[activeImageIndex.value] || null)
const visibilityLabel = computed(() => ({
  public: '公开可见',
  private: '仅自己可见',
  mutual: '仅互关可见',
  selected: '指定用户可见',
  excluded: '部分用户不可见'
}[props.detail?.visibility] || props.detail?.visibility))
const statusMeta = computed(() => {
  if (props.detail?.reviewStatus === 'rejected') return { label: '审核未通过', tone: 'danger' }
  if (props.detail?.reviewStatus === 'pending_review') return { label: '审核中', tone: 'warning' }
  if (props.detail?.publishStatus === 'scheduled') return { label: '等待发布', tone: 'info' }
  if (props.detail?.publishStatus === 'published') return { label: '已发布', tone: 'success' }
  return { label: '审核通过', tone: 'success' }
})
const mediaSummary = computed(() => props.detail?.type === 'video'
  ? `HLS 流媒体 · ${formatDuration(props.detail.durationMs)}`
  : `${imageAssets.value.length} 张图片`)

watch(() => props.detail?.workId, () => {
  activeImageIndex.value = 0
})

onMounted(() => window.addEventListener('keydown', handleKeydown))
onBeforeUnmount(() => window.removeEventListener('keydown', handleKeydown))

function handleKeydown(event) {
  if (event.key === 'Escape') emit('close')
  if (props.detail?.type !== 'image' || imageAssets.value.length < 2) return
  if (event.key === 'ArrowRight') {
    activeImageIndex.value = (activeImageIndex.value + 1) % imageAssets.value.length
  }
  if (event.key === 'ArrowLeft') {
    activeImageIndex.value = (activeImageIndex.value - 1 + imageAssets.value.length) % imageAssets.value.length
  }
}

function formatDate(value) {
  if (!value) return '—'
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  }).format(new Date(value)).replace(/\//g, '-')
}

function formatDuration(value) {
  const total = Math.max(0, Math.floor((value || 0) / 1000))
  const hours = Math.floor(total / 3600)
  const minutes = Math.floor((total % 3600) / 60)
  const seconds = total % 60
  const shortValue = `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
  return hours ? `${String(hours).padStart(2, '0')}:${shortValue}` : shortValue
}
</script>
