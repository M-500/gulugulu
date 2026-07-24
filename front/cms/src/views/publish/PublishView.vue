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
          :assets="coverAssets"
          :cover-asset-id="currentDraft.coverAssetId"
          @select="updateDraft({ coverAssetId: $event })"
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

      <PublishPreview :asset="previewAsset" :assets="coverAssets" :form="form" />
      <div class="publish-bar">
        <button class="ghost-button" type="button" @click="saveAndLeave">暂存离开</button>
        <button class="primary-button" type="button" @click="publishDraft">发布</button>
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
const uploadMessage = ref('')
const bodyTopics = ref([])

const form = reactive({
  title: '',
  body: '',
  collection: '',
  visibility: 'public',
  original: false,
  scheduled: false,
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
  const maxFiles = draftType === 'imageText' ? 19 : 20
  const files = Array.from(fileList).filter((file) => file.type.startsWith('video/') || file.type.startsWith('image/')).slice(0, maxFiles)

  if (!files.length) {
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
      contentType: 'application/octet-stream'
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
  await persistForm()
  alert('发布接口接入后会提交当前作品，草稿已暂存到本地。')
}

</script>

<style src="./publish.css"></style>
