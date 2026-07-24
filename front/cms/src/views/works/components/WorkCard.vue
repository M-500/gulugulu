<template>
  <article class="work-card">
    <div class="work-card__cover">
      <img v-if="work.coverUrl" :src="work.coverUrl" :alt="work.title" loading="lazy">
      <div v-else class="work-card__placeholder"><span>咕</span></div>
      <i v-if="work.type === 'video'" class="work-card__play">▶</i>
      <small v-if="work.type === 'video'" class="work-card__duration">{{ formatDuration(work.durationMs) }}</small>
      <small v-if="work.visibility === 'private'" class="work-card__private">仅自己可见</small>
    </div>

    <div class="work-card__body">
      <div class="work-card__title-row">
        <span v-if="statusMeta.badge" :class="['work-status', `is-${statusMeta.tone}`]">{{ statusMeta.badge }}</span>
        <h2>{{ work.title || '无笔记标题' }}</h2>
      </div>
      <details class="work-card__actions">
        <summary title="作品操作" aria-label="作品操作">✎</summary>
        <div class="work-card__menu">
          <button type="button" @click="$emit('edit', work)">编辑标题</button>
          <span>设置权限</span>
          <button type="button" :class="{ 'is-selected': work.visibility === 'public' }" @click="$emit('visibility', { work, visibility: 'public' })">公开可见</button>
          <button type="button" :class="{ 'is-selected': work.visibility === 'private' }" @click="$emit('visibility', { work, visibility: 'private' })">仅自己可见</button>
          <button type="button" :class="{ 'is-selected': work.visibility === 'mutual' }" @click="$emit('visibility', { work, visibility: 'mutual' })">仅互关可见</button>
          <button class="is-danger" type="button" @click="$emit('delete', work)">删除作品</button>
        </div>
      </details>
      <p v-if="work.reviewStatus === 'rejected'" class="work-card__reason">
        {{ work.reviewReason || '作品未通过审核，请修改后重新提交。' }}
      </p>
      <p v-else class="work-card__status-copy">{{ statusMeta.description }}</p>
      <time :datetime="work.createdAt">{{ formatDate(work.createdAt) }}</time>

    </div>
  </article>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  work: { type: Object, required: true }
})

defineEmits(['delete', 'edit', 'visibility'])

const statusMeta = computed(() => {
  if (props.work.reviewStatus === 'rejected') {
    return { badge: '未通过', tone: 'danger', description: '', footer: '审核未通过' }
  }
  if (props.work.reviewStatus === 'pending_review') {
    return { badge: '审核中', tone: 'warning', description: '作品已处理完成，正在等待平台审核。', footer: '等待审核' }
  }
  if (props.work.publishStatus === 'scheduled') {
    return { badge: '待发布', tone: 'info', description: '审核已通过，将在设定时间自动发布。', footer: '定时发布' }
  }
  if (props.work.publishStatus === 'published') {
    return { badge: '', tone: 'success', description: '作品已通过审核并成功发布。', footer: '已发布' }
  }
  return { badge: '已通过', tone: 'success', description: '作品已经通过审核，等待发布。', footer: '审核通过' }
})

function formatDate(value) {
  if (!value) {
    return ''
  }
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  }).format(new Date(value)).replace(/\//g, '-')
}

function formatDuration(durationMs) {
  const totalSeconds = Math.max(0, Math.floor((durationMs || 0) / 1000))
  const hours = Math.floor(totalSeconds / 3600)
  const minutes = Math.floor((totalSeconds % 3600) / 60)
  const seconds = totalSeconds % 60
  const result = [minutes, seconds].map((item) => String(item).padStart(2, '0')).join(':')
  return hours ? `${String(hours).padStart(2, '0')}:${result}` : result
}

</script>
