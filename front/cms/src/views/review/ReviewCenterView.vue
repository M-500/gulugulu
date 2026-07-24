<template>
  <section class="review-page">
    <ReviewToolbar
      v-model:keyword="keyword"
      v-model:type="workType"
      :active-status="activeStatus"
      :tabs="tabs"
      @search="handleSearch"
      @select-status="handleStatusSelect"
    />

    <div v-if="loading" class="review-state">
      <span class="review-state__spinner" />
      <strong>正在获取审核队列</strong>
      <p>正在同步最新待审核作品…</p>
    </div>
    <div v-else-if="errorMessage" class="review-state is-error">
      <span>!</span>
      <strong>审核队列加载失败</strong>
      <p>{{ errorMessage }}</p>
      <button type="button" @click="loadWorks">重新加载</button>
    </div>
    <div v-else-if="works.length" class="review-grid">
      <ReviewCard v-for="work in works" :key="work.workId" :work="work" @open="openWork" />
    </div>
    <div v-else class="review-state">
      <span>✓</span>
      <strong>{{ activeStatus === 'pending' ? '待审核队列已清空' : '当前筛选暂无记录' }}</strong>
      <p>{{ activeStatus === 'pending' ? '所有已处理完成的作品都已审核。' : '尝试切换状态或调整搜索条件。' }}</p>
    </div>

    <div v-if="total > pageSize" class="review-pagination">
      <button type="button" :disabled="page <= 1" @click="changePage(page - 1)">上一页</button>
      <span>第 {{ page }} / {{ totalPages }} 页</span>
      <button type="button" :disabled="page >= totalPages" @click="changePage(page + 1)">下一页</button>
    </div>

    <ReviewDetailDrawer
      v-if="selectedWork"
      :detail="detail"
      :error-message="detailError"
      :loading="detailLoading"
      :submitting="submitting"
      @approve="approveWork"
      @close="closeDetail"
      @reject="rejecting = true"
    />
    <RejectDialog
      v-if="rejecting"
      :submitting="submitting"
      @close="rejecting = false"
      @submit="rejectWork"
    />
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'

import { auditWork, getAuditWorkDetail, getAuditWorks } from '@/api/reviews'

import RejectDialog from './components/RejectDialog.vue'
import ReviewCard from './components/ReviewCard.vue'
import ReviewDetailDrawer from './components/ReviewDetailDrawer.vue'
import ReviewToolbar from './components/ReviewToolbar.vue'

const emit = defineEmits(['count-change'])
const activeStatus = ref('pending')
const workType = ref('')
const keyword = ref('')
const submittedKeyword = ref('')
const page = ref(1)
const pageSize = 12
const total = ref(0)
const works = ref([])
const counts = ref({ pending: 0, approved: 0, rejected: 0 })
const loading = ref(false)
const errorMessage = ref('')
const selectedWork = ref(null)
const detail = ref(null)
const detailLoading = ref(false)
const detailError = ref('')
const rejecting = ref(false)
const submitting = ref(false)

const tabs = computed(() => [
  { key: 'pending', label: '待审核', count: counts.value.pending },
  { key: 'approved', label: '已通过', count: counts.value.approved },
  { key: 'rejected', label: '已驳回', count: counts.value.rejected },
  { key: 'all', label: '全部记录', count: counts.value.pending + counts.value.approved + counts.value.rejected }
])
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)))

onMounted(loadWorks)

async function loadWorks() {
  loading.value = true
  errorMessage.value = ''
  try {
    const result = await getAuditWorks({
      page: page.value,
      pageSize,
      status: activeStatus.value,
      type: workType.value || undefined,
      keyword: submittedKeyword.value || undefined
    })
    works.value = result.list || []
    total.value = result.total || 0
    counts.value = {
      pending: result.pendingCount || 0,
      approved: result.approvedCount || 0,
      rejected: result.rejectedCount || 0
    }
    emit('count-change', counts.value.pending)
  } catch (error) {
    works.value = []
    total.value = 0
    errorMessage.value = error.message || '请求失败，请稍后重试'
  } finally {
    loading.value = false
  }
}

function handleStatusSelect(status) {
  activeStatus.value = status
  page.value = 1
  loadWorks()
}

function handleSearch() {
  submittedKeyword.value = keyword.value.trim()
  page.value = 1
  loadWorks()
}

function changePage(value) {
  page.value = value
  loadWorks()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

async function openWork(work) {
  selectedWork.value = work
  detail.value = null
  detailError.value = ''
  detailLoading.value = true
  try {
    detail.value = await getAuditWorkDetail(work.workId)
  } catch (error) {
    detailError.value = error.message || '作品详情加载失败'
  } finally {
    detailLoading.value = false
  }
}

function closeDetail() {
  if (submitting.value) return
  selectedWork.value = null
  detail.value = null
  rejecting.value = false
}

async function approveWork() {
  if (!window.confirm(`确认通过作品“${detail.value.title}”吗？审核通过后将按作品发布设置生效。`)) return
  await submitDecision('approve')
}

async function rejectWork(reason) {
  await submitDecision('reject', reason)
}

async function submitDecision(decision, reason = '') {
  submitting.value = true
  try {
    await auditWork(detail.value.workId, decision, reason)
    rejecting.value = false
    selectedWork.value = null
    detail.value = null
    await loadWorks()
  } catch (error) {
    window.alert(error.message || '提交审核结果失败')
  } finally {
    submitting.value = false
  }
}
</script>

<style src="./review.css"></style>
