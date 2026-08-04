<template>
  <article class="work-card" tabindex="0" @click="$emit('open', work)" @keydown.enter="$emit('open', work)">
    <div class="work-card__cover">
      <img v-if="work.coverUrl" :src="work.coverUrl" :alt="work.title" loading="lazy">
      <div v-else class="work-card__placeholder"><span>咕</span></div>
      <i v-if="work.type === 'video'" class="work-card__play">▶</i>
      <small v-if="work.type === 'video'" class="work-card__duration">{{ formatDuration(work.durationMs) }}</small>
      <small v-if="work.visibility === 'private'" class="work-card__private">
        <el-icon><Lock /></el-icon>仅自己可见
      </small>
    </div>

    <div class="work-card__body">
      <div class="work-card__title-row">
        <span v-if="statusMeta.badge" :class="['work-status', `is-${statusMeta.tone}`]">{{ statusMeta.badge }}</span>
        <h2>{{ work.title || '无笔记标题' }}</h2>
      </div>

      <div class="work-card__actions" @click.stop>
        <el-tooltip content="修改可见范围" placement="top">
          <button type="button" aria-label="修改可见范围" @click="$emit('visibility', work)"><User /></button>
        </el-tooltip>
        <el-tooltip content="编辑作品" placement="top">
          <button type="button" aria-label="编辑作品" @click="$emit('edit', work)"><EditPen /></button>
        </el-tooltip>
        <el-tooltip content="删除作品" placement="top">
          <button class="is-danger" type="button" aria-label="删除作品" @click="$emit('delete', work)"><Delete /></button>
        </el-tooltip>
      </div>

      <p v-if="work.reviewStatus === 'rejected'" class="work-card__reason">{{ work.reviewReason || '作品未通过审核，请修改后重新提交。' }}</p>
      <p v-else class="work-card__status-copy">{{ statusMeta.description }}</p>
      <time :datetime="work.createdAt">{{ formatDate(work.createdAt) }}</time>

      <div class="work-card__stats" aria-label="作品数据">
        <span v-for="stat in stats" :key="stat.key" :title="stat.label">
          <el-icon><component :is="stat.icon" /></el-icon>{{ stat.value }}
        </span>
      </div>
    </div>
  </article>
</template>

<script setup>
import { computed } from 'vue'
import { ChatDotRound, CollectionTag, Delete, EditPen, Promotion, Star, User, View, Lock } from '@element-plus/icons-vue'

const props = defineProps({ work: { type: Object, required: true } })
defineEmits(['delete', 'edit', 'open', 'visibility'])

const stats = computed(() => [
  { key: 'view', label: '浏览量', icon: View, value: props.work.viewCount || 0 },
  { key: 'comment', label: '评论数', icon: ChatDotRound, value: props.work.commentCount || 0 },
  { key: 'like', label: '点赞数', icon: Star, value: props.work.likeCount || 0 },
  { key: 'favorite', label: '收藏数', icon: CollectionTag, value: props.work.favoriteCount || 0 },
  { key: 'share', label: '分享数', icon: Promotion, value: props.work.shareCount || 0 }
])

const statusMeta = computed(() => {
  if (props.work.reviewStatus === 'rejected') return { badge: '未通过', tone: 'danger', description: '' }
  if (props.work.reviewStatus === 'pending_review') return { badge: '审核中', tone: 'warning', description: '作品已处理完成，正在等待平台审核。' }
  if (props.work.publishStatus === 'scheduled') return { badge: '待发布', tone: 'info', description: '审核已通过，将在设定时间自动发布。' }
  if (props.work.publishStatus === 'published') return { badge: '', tone: 'success', description: '作品已通过审核并成功发布。' }
  return { badge: '已通过', tone: 'success', description: '作品已经通过审核，等待发布。' }
})

function formatDate(value) {
  if (!value) return ''
  return new Intl.DateTimeFormat('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', hour12: false }).format(new Date(value)).replace(/\//g, '-')
}

function formatDuration(durationMs) {
  const total = Math.max(0, Math.floor((durationMs || 0) / 1000))
  const result = [Math.floor((total % 3600) / 60), total % 60].map((item) => String(item).padStart(2, '0')).join(':')
  return total >= 3600 ? `${String(Math.floor(total / 3600)).padStart(2, '0')}:${result}` : result
}
</script>
