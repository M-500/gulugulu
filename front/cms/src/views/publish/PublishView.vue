<template>
  <section class="publish-page">
    <PublishUploadStage
      v-if="!currentDraft"
      v-model="activeType"
      v-model:dragging="dragging"
      :publish-types="publishTypes"
      :active-option="activeOption"
      :draft-count="draftStore.drafts.length"
      @open-drafts="openDraftDrawer"
      @select-files="handleFileInput"
      @drop="handleDrop"
    />

    <div v-else class="editor-stage">
      <div class="editor-main">
        <PublishAssetsCard
          :type="currentDraft.type"
          :accept="activeDraftOption.accept"
          :assets="currentDraft.assets"
          @replace="replaceFiles"
          @add-images="appendImageFiles"
          @delete-image="deleteImage"
        />
        <PublishCoverCard
          :type="currentDraft.type"
          :assets="coverAssets"
          :cover-asset-id="currentDraft.coverAssetId"
          :custom-cover-url="currentDraft.coverLocalUrl"
          @select="updateDraft({ coverAssetId: $event })"
          @upload-cover="handleCoverInput"
        />
        <PublishComposeCard
          :form="form"
          :topics="defaultTopics"
          :selected-topics="form.topics"
          @update-form="Object.assign(form, $event)"
          @sync-body-topics="syncBodyTopics"
          @persist="persistForm"
          @toggle-topic="toggleTopic"
        />
        <PublishSettings :form="form" @update-form="Object.assign(form, $event)" @persist="persistForm" @open-collection="collectionOpen = true" />
      </div>

      <PublishPreview :asset="previewAsset" :assets="coverAssets" :cover-url="currentDraft.coverLocalUrl" :form="form" />
      <div class="publish-bar">
        <button class="ghost-button" type="button" @click="saveAndLeave">暂存离开</button>
        <button class="primary-button" type="button" :disabled="publishing" @click="publishDraft">{{ publishing ? '提交中…' : '发布' }}</button>
      </div>
    </div>

    <div v-if="uploading" class="upload-progress"><strong>正在上传素材</strong><span>{{ uploadMessage }}</span></div>
    <PublishDraftDrawer
      v-if="draftDrawerOpen"
      v-model="drawerType"
      :tabs="drawerTabs"
      :drafts="filteredDrafts"
      @close="draftDrawerOpen = false"
      @edit="editDraft"
      @delete="deleteDraft"
    />
    <div v-if="collectionOpen" class="collection-popover" @click.self="collectionOpen = false">
      <div><h3>添加视频合集</h3><input v-model.trim="collectionName" placeholder="输入合集名称，例如：项目工程日志"><button type="button" @click="applyCollection">创建并加入合集</button></div>
    </div>
  </section>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'

import { completeUpload, createUploadPresign, uploadToObjectStorage } from '@/api/media'
import { createWork } from '@/api/works'
import { usePublishDraftStore } from '@/stores/publishDrafts'

import PublishAssetsCard from './components/PublishAssetsCard.vue'
import PublishComposeCard from './components/PublishComposeCard.vue'
import PublishCoverCard from './components/PublishCoverCard.vue'
import PublishDraftDrawer from './components/PublishDraftDrawer.vue'
import PublishPreview from './components/PublishPreview.vue'
import PublishSettings from './components/PublishSettings.vue'
import PublishUploadStage from './components/PublishUploadStage.vue'
import { defaultTopics, publishTypes } from './publishOptions'

const draftStore = usePublishDraftStore()

const activeType = ref('video')
const drawerType = ref('video')
const dragging = ref(false)
const draftDrawerOpen = ref(false)
const collectionOpen = ref(false)
const collectionName = ref('')
const uploading = ref(false)
const publishing = ref(false)
const uploadMessage = ref('')
const bodyTopics = ref([])

const form = reactive({
  title: '',
  body: '',
  collection: '',
  visibility: 'public',
  original: false,
  scheduled: false,
  scheduledAt: '',
  topics: []
})

const currentDraft = computed(() => draftStore.currentDraft)
const activeOption = computed(() => publishTypes.find((item) => item.key === activeType.value))
const activeDraftOption = computed(() => publishTypes.find((item) => item.key === currentDraft.value?.type) || activeOption.value)
const coverAssets = computed(() => currentDraft.value?.assets || [])
const previewAsset = computed(() => currentDraft.value?.assets?.find((asset) => asset.kind === 'video') || currentDraft.value?.assets?.[0])
const drawerTabs = computed(() => publishTypes.map((item) => ({
  key: item.key,
  label: `${item.shortLabel}笔记`,
  count: draftStore.draftCountByType(item.key)
})))
const filteredDrafts = computed(() => draftStore.sortedDrafts.filter((draft) => draft.type === drawerType.value))

watch(currentDraft, (draft) => {
  if (!draft) {
    return
  }

  form.title = draft.title || ''
  form.body = draft.body || ''
  form.collection = draft.collection || ''
  form.visibility = draft.visibility || 'public'
  form.original = Boolean(draft.original)
  form.scheduled = Boolean(draft.scheduled)
  form.scheduledAt = draft.scheduledAt || ''
  form.topics = [...(draft.topics || [])]
  bodyTopics.value = extractTopics(draft.body || '')
}, { immediate: true })

onMounted(() => {
  draftStore.loadDrafts()
})

onBeforeUnmount(() => {
  draftStore.releaseCurrentUrls()
})

async function handleFileInput(event) {
  await createDraftFromFiles(event.target.files)
  event.target.value = ''
}

async function handleDrop(event) {
  dragging.value = false
  await createDraftFromFiles(event.dataTransfer.files)
}

async function replaceFiles(event) {
  const files = Array.from(event.target.files)
  event.target.value = ''

  if (!files.length) {
    return
  }

  await createDraftFromFiles(files, currentDraft.value.type)
}

async function appendImageFiles(event) {
  const remaining = 19 - currentDraft.value.assets.length
  const files = Array.from(event.target.files).filter((file) => file.type.startsWith('image/')).slice(0, remaining)
  event.target.value = ''

  if (!files.length) {
    return
  }

  uploading.value = true
  uploadMessage.value = `准备上传 ${files.length} 张图片`

  try {
    const uploadedAssets = await uploadFiles(files)
    await draftStore.appendImages(files, uploadedAssets)
  } finally {
    uploading.value = false
    uploadMessage.value = ''
  }
}

async function deleteImage(assetId) {
  await draftStore.deleteCurrentImage(assetId)
}

async function createDraftFromFiles(fileList, draftType = activeType.value) {
  const isImageWork = draftType === 'image'
  const files = Array.from(fileList)
    .filter((file) => isImageWork ? file.type.startsWith('image/') : file.type.startsWith('video/'))
    .slice(0, isImageWork ? 19 : 1)

  if (!files.length) {
    alert(isImageWork ? '图片作品只能上传图片。' : '视频作品只能上传一个视频。')
    return
  }

  uploading.value = true
  uploadMessage.value = `准备上传 ${files.length} 个素材`

  try {
    const uploadedAssets = await uploadFiles(files)
    await draftStore.createDraft(draftType, files, uploadedAssets)
  } finally {
    uploading.value = false
    uploadMessage.value = ''
  }
}

async function uploadFiles(files) {
  const uploadedAssets = []

  for (const [index, file] of files.entries()) {
    uploadMessage.value = `正在上传 ${index + 1}/${files.length}：${file.name}`
    const presign = await createUploadPresign({
      resourceType: file.type.startsWith('video/') ? 'Video' : 'Image',
      fileName: file.name,
      contentType: file.type || 'application/octet-stream'
    })

    await uploadToObjectStorage(file, presign)
    const completed = await completeUpload({ mediaId: presign.mediaId, objectKey: presign.objectKey })
    uploadedAssets.push({
      mediaId: completed.mediaId,
      bucket: completed.bucket,
      objectKey: completed.objectKey,
      previewUrl: completed.previewUrl || presign.previewUrl,
      status: completed.status
    })
  }

  return uploadedAssets
}

async function openDraftDrawer() {
  await draftStore.loadDrafts()
  draftDrawerOpen.value = true
}

async function editDraft(id) {
  await draftStore.openDraft(id)
  draftDrawerOpen.value = false
}

async function deleteDraft(id) {
  await draftStore.deleteDraft(id)
}

async function updateDraft(patch) {
  await draftStore.saveCurrentDraft(patch)
}

async function persistForm() {
  await updateDraft({
    title: form.title,
    body: form.body,
    collection: form.collection,
    visibility: form.visibility,
    original: form.original,
    scheduled: form.scheduled,
    scheduledAt: form.scheduledAt,
    topics: form.topics
  })
}

function toggleTopic(topic) {
  const existed = form.topics.includes(topic)

  form.topics = existed ? form.topics.filter((item) => item !== topic) : [...form.topics, topic]
  persistForm()
}

function syncBodyTopics(topics) {
  const manualTopics = form.topics.filter((topic) => !bodyTopics.value.includes(topic))
  bodyTopics.value = topics
  form.topics = [...new Set([...manualTopics, ...topics])]
}

function extractTopics(value) {
  const matches = value.matchAll(/#([^\s#，。！？、,.!?；;：:]+)/g)
  return [...new Set(Array.from(matches, (match) => match[1]).filter(Boolean))]
}

function applyCollection() {
  if (!collectionName.value) {
    return
  }

  form.collection = collectionName.value
  collectionName.value = ''
  collectionOpen.value = false
  persistForm()
}

async function saveAndLeave() {
  await persistForm()
  draftStore.releaseCurrentUrls()
  draftStore.currentDraft = null
}

async function publishDraft() {
  if (publishing.value) {
    return
  }

  publishing.value = true
  uploading.value = true
  uploadMessage.value = '正在校验发布信息'

  try {
    await persistForm()
    validatePublishForm()
    const cover = await resolveCoverFile()
    const payload = buildPublishPayload()

    uploadMessage.value = '正在创建作品'
    const created = await createWork(payload, cover, currentDraft.value.idempotencyKey)

    const draftId = currentDraft.value.id
    await draftStore.deleteDraft(draftId)
    alert(`作品 #${created.workId} 已提交，后台将自动处理素材并进入审核。`)
  } catch (error) {
    alert(error.message || '发布失败，请稍后重试')
  } finally {
    publishing.value = false
    uploading.value = false
    uploadMessage.value = ''
  }
}

async function handleCoverInput(event) {
  const file = event.target.files?.[0]
  event.target.value = ''
  if (!file) {
    return
  }
  if (!['image/jpeg', 'image/png'].includes(file.type)) {
    alert('封面只支持 JPG、JPEG 或 PNG 格式。')
    return
  }
  if (file.size > 100 * 1024 * 1024) {
    alert('封面大小不能超过 100MB。')
    return
  }
  await updateDraft({
    coverBlob: file,
    coverName: file.name,
    coverType: file.type
  })
}

function validatePublishForm() {
  const draft = currentDraft.value
  if (!form.title.trim()) {
    throw new Error('请填写作品标题')
  }
  if (Array.from(form.title.trim()).length > 50) {
    throw new Error('作品标题不能超过50个字符')
  }
  if (!draft.assets.length || draft.assets.some((asset) => !asset.mediaId || asset.uploadStatus !== 'uploaded')) {
    throw new Error('存在尚未上传完成的作品素材')
  }
  if (draft.type === 'video' && draft.assets.length !== 1) {
    throw new Error('视频作品必须且只能包含一个视频')
  }
  if (draft.type === 'image' && (draft.assets.length < 1 || draft.assets.length > 19)) {
    throw new Error('图片作品必须包含1到19张图片')
  }
  if (form.scheduled) {
    const scheduledTime = new Date(form.scheduledAt)
    if (!form.scheduledAt || Number.isNaN(scheduledTime.getTime()) || scheduledTime.getTime() < Date.now() + 5 * 60 * 1000) {
      throw new Error('定时发布时间必须至少晚于当前时间5分钟')
    }
  }
}

function buildPublishPayload() {
  return {
    type: currentDraft.value.type,
    title: form.title.trim(),
    content: form.body,
    visibility: {
      type: form.visibility,
      userIds: []
    },
    topics: form.topics.map((name) => ({ id: 0, name })),
    assets: currentDraft.value.assets.map((asset, index) => ({
      mediaId: asset.mediaId,
      sort: index
    })),
    collectionId: 0,
    original: form.original,
    scheduledAt: form.scheduled ? new Date(form.scheduledAt).toISOString() : null
  }
}

async function resolveCoverFile() {
  const draft = currentDraft.value
  if (draft.coverBlob) {
    return new File([draft.coverBlob], draft.coverName || 'cover.jpg', {
      type: draft.coverType || draft.coverBlob.type || 'image/jpeg'
    })
  }
  if (draft.type === 'video') {
    throw new Error('请为视频作品上传一张封面图')
  }

  const selected = draft.assets.find((asset) => asset.assetId === draft.coverAssetId) || draft.assets[0]
  if (!selected?.blob) {
    throw new Error('无法读取所选封面，请重新上传封面图')
  }
  if (['image/jpeg', 'image/png'].includes(selected.blob.type)) {
    const extension = selected.blob.type === 'image/png' ? 'png' : 'jpg'
    return new File([selected.blob], `cover.${extension}`, { type: selected.blob.type })
  }
  return convertImageToPng(selected.blob)
}

async function convertImageToPng(blob) {
  const bitmap = await createImageBitmap(blob)
  const canvas = document.createElement('canvas')
  canvas.width = bitmap.width
  canvas.height = bitmap.height
  canvas.getContext('2d').drawImage(bitmap, 0, 0)
  bitmap.close()
  const pngBlob = await new Promise((resolve, reject) => {
    canvas.toBlob((result) => result ? resolve(result) : reject(new Error('封面格式转换失败')), 'image/png')
  })
  return new File([pngBlob], 'cover.png', { type: 'image/png' })
}

</script>

<style src="./publish.css"></style>
