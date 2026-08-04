<template>
  <section class="works-page">
    <WorksToolbar
      v-model:keyword="keyword"
      :active-status="activeStatus"
      :tabs="tabs"
      @search="handleSearch"
      @select-status="handleStatusSelect"
    />

    <div v-if="loading" class="works-state-card">
      <span class="works-state-card__spinner" />
      <strong>正在加载作品</strong>
      <p>正在同步最新审核和发布状态…</p>
    </div>

    <div v-else-if="errorMessage" class="works-state-card works-state-card--error">
      <span>!</span>
      <strong>作品加载失败</strong>
      <p>{{ errorMessage }}</p>
      <button type="button" @click="loadWorks">重新加载</button>
    </div>

    <WorksGrid
      v-else-if="works.length"
      :works="works"
      @delete="handleDelete"
      @edit="handleEdit"
      @open="handleOpen"
      @visibility="handleVisibility"
    />

    <div v-else class="works-state-card">
      <span>▤</span>
      <strong>{{ keyword ? '没有找到相关作品' : '当前分类暂无作品' }}</strong>
      <p>{{ keyword ? '请尝试更换关键词重新搜索。' : '媒体处理完成后，作品会自动出现在这里。' }}</p>
    </div>

    <WorksPagination
      v-if="!loading && total > 0"
      :page="page"
      :page-size="pageSize"
      :total="total"
      @change="handlePageChange"
    />

    <WorkDetailDialog
      v-if="selectedWork"
      :detail="workDetail"
      :error-message="detailError"
      :loading="detailLoading"
      @close="closeDetail"
      @retry="loadWorkDetail"
    />
    <WorkVisibilityDialog
      v-model="visibilityDialogOpen"
      :work="visibilityWork"
      :saving="visibilitySaving"
      @confirm="saveVisibility"
    />
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'

import {
  deleteCreatorWork,
  getCreatorWorks,
  getWork,
  updateCreatorWorkVisibility
} from '@/api/works'

import WorksGrid from './components/WorksGrid.vue'
import WorkDetailDialog from './components/WorkDetailDialog.vue'
import WorkVisibilityDialog from './components/WorkVisibilityDialog.vue'
import WorksPagination from './components/WorksPagination.vue'
import WorksToolbar from './components/WorksToolbar.vue'

const emit = defineEmits(['edit-work'])

const activeStatus = ref('all')
const keyword = ref('')
const submittedKeyword = ref('')
const page = ref(1)
const pageSize = 12
const total = ref(0)
const works = ref([])
const loading = ref(false)
const errorMessage = ref('')
const selectedWork = ref(null)
const workDetail = ref(null)
const detailLoading = ref(false)
const detailError = ref('')
const visibilityDialogOpen = ref(false)
const visibilityWork = ref(null)
const visibilitySaving = ref(false)
const counts = ref({
  all: 0,
  published: 0,
  reviewing: 0,
  rejected: 0
})

const tabs = computed(() => [
  { key: 'all', label: '全部', count: counts.value.all },
  { key: 'published', label: '已发布', count: counts.value.published },
  { key: 'reviewing', label: '审核中', count: counts.value.reviewing },
  { key: 'rejected', label: '未通过', count: counts.value.rejected }
])

onMounted(loadWorks)

async function loadWorks() {
  loading.value = true
  errorMessage.value = ''

  try {
    const result = await getCreatorWorks({
      page: page.value,
      pageSize,
      status: activeStatus.value,
      keyword: submittedKeyword.value || undefined
    })
    works.value = result.list || []
    total.value = result.total || 0
    counts.value = {
      all: result.allCount || 0,
      published: result.publishedCount || 0,
      reviewing: result.reviewingCount || 0,
      rejected: result.rejectedCount || 0
    }
  } catch (error) {
    works.value = []
    total.value = 0
    errorMessage.value = error.message || '请求失败，请稍后重试'
  } finally {
    loading.value = false
  }
}

function handleStatusSelect(status) {
  if (activeStatus.value === status) {
    return
  }
  activeStatus.value = status
  page.value = 1
  loadWorks()
}

function handleSearch() {
  submittedKeyword.value = keyword.value.trim()
  page.value = 1
  loadWorks()
}

function handlePageChange(nextPage) {
  page.value = nextPage
  loadWorks()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function handleOpen(work) {
  selectedWork.value = work
  workDetail.value = null
  loadWorkDetail()
}

async function loadWorkDetail() {
  if (!selectedWork.value) return
  detailLoading.value = true
  detailError.value = ''
  try {
    workDetail.value = await getWork(selectedWork.value.workId)
  } catch (error) {
    detailError.value = error.message || '作品详情加载失败'
  } finally {
    detailLoading.value = false
  }
}

function closeDetail() {
  selectedWork.value = null
  workDetail.value = null
  detailError.value = ''
}

function handleEdit(work) {
  emit('edit-work', work.workId)
}

function handleVisibility(work) {
  visibilityWork.value = work
  visibilityDialogOpen.value = true
}

async function saveVisibility(visibility) {
  const work = visibilityWork.value
  if (!work || work.visibility === visibility) {
    visibilityDialogOpen.value = false
    return
  }
  visibilitySaving.value = true
  try {
    await updateCreatorWorkVisibility(work.workId, visibility)
    work.visibility = visibility
    visibilityDialogOpen.value = false
  } catch (error) {
    window.alert(error.message || '修改作品权限失败')
  } finally {
    visibilitySaving.value = false
  }
}

async function handleDelete(work) {
  if (!window.confirm(`确认删除作品“${work.title}”吗？删除后将不再展示。`)) {
    return
  }
  try {
    await deleteCreatorWork(work.workId)
    await loadWorks()
  } catch (error) {
    window.alert(error.message || '删除作品失败')
  }
}
</script>

<style src="./works.css"></style>
