<template>
  <article class="review-card" tabindex="0" @click="$emit('open', work)" @keydown.enter="$emit('open', work)">
    <div class="review-card__cover">
      <img v-if="work.coverUrl" :src="work.coverUrl" :alt="work.title" loading="lazy">
      <span v-else>咕</span>
      <i>{{ work.type === 'video' ? '视频' : '图文' }}</i>
      <small v-if="work.type === 'video'">{{ formatDuration(work.durationMs) }}</small>
    </div>
    <div class="review-card__body">
      <div class="review-card__heading">
        <span :class="['review-badge', `is-${statusMeta.tone}`]">{{ statusMeta.label }}</span>
        <time>{{ formatDate(work.submittedAt) }}</time>
      </div>
      <h2>{{ work.title || '无标题作品' }}</h2>
      <p>{{ work.contentExcerpt || '作者未填写作品正文。' }}</p>
      <footer>
        <span class="review-card__avatar">{{ (work.authorName || '?').slice(0, 1) }}</span>
        <div>
          <strong>{{ work.authorName }}</strong>
          <small>ID {{ work.authorId }} · {{ visibilityLabel }}</small>
        </div>
        <button type="button">开始审核 →</button>
      </footer>
    </div>
  </article>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  work: { type: Object, required: true }
})

defineEmits(['open'])

const statusMeta = computed(() => {
  if (props.work.reviewStatus === 'approved') return { label: '已通过', tone: 'success' }
  if (props.work.reviewStatus === 'rejected') return { label: '已驳回', tone: 'danger' }
  return { label: '待审核', tone: 'warning' }
})

const visibilityLabel = computed(() => ({
  public: '公开可见',
  private: '仅自己可见',
  mutual: '仅互关可见',
  selected: '指定用户可见',
  excluded: '部分用户不可见'
}[props.work.visibility] || props.work.visibility))

function formatDate(value) {
  if (!value) return ''
  return new Intl.DateTimeFormat('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  }).format(new Date(value))
}

function formatDuration(value) {
  const total = Math.max(0, Math.floor((value || 0) / 1000))
  const minutes = Math.floor(total / 60)
  return `${String(minutes).padStart(2, '0')}:${String(total % 60).padStart(2, '0')}`
}
</script>
