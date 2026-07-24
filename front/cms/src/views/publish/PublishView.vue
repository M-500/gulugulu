<template>
  <section class="publish-page">
    <div v-if="!currentDraft" class="upload-stage">
      <div class="publish-tabs">
        <button
          v-for="item in publishTypes"
          :key="item.key"
          type="button"
          :class="{ 'is-active': activeType === item.key }"
          @click="activeType = item.key"
        >
          {{ item.label }}
        </button>

        <button class="draft-entry" type="button" @click="openDraftDrawer">
          草稿箱({{ draftStore.drafts.length }})
        </button>
      </div>

      <label
        class="upload-dropzone"
        :class="{ 'is-dragging': dragging }"
        @dragenter.prevent="dragging = true"
        @dragover.prevent
        @dragleave.prevent="dragging = false"
        @drop.prevent="handleDrop"
      >
        <input :accept="activeOption.accept" multiple type="file" @change="handleFileInput">
        <span class="upload-icon">↑</span>
        <strong>{{ activeOption.hint }}</strong>
        <button type="button">{{ activeOption.buttonText }}</button>
        <small>{{ activeOption.description }}</small>
      </label>

      <div class="upload-tips">
        <div>
          <strong>文件大小</strong>
          <span>单个素材建议小于 4GB，最多支持 20 个文件</span>
        </div>
        <div>
          <strong>素材格式</strong>
          <span>视频支持 mp4、mov，图片支持 jpg、png、webp</span>
        </div>
        <div>
          <strong>编辑体验</strong>
          <span>上传后自动保存到本地草稿，可继续编辑发布</span>
        </div>
      </div>
    </div>

    <div v-else class="editor-stage">
      <div class="editor-main">
        <section class="editor-card asset-card">
          <div class="editor-card__header">
            <h2>{{ currentDraft.type === 'imageText' ? '图文文件' : currentDraft.type === 'mixed' ? '混合作品素材' : '视频文件' }}</h2>
            <label class="reupload-button">
              重新上传
              <input :accept="activeDraftOption.accept" multiple type="file" @change="replaceFiles">
            </label>
          </div>

          <div class="asset-summary">
            <div v-for="asset in currentDraft.assets" :key="asset.assetId" class="asset-item">
              <span :class="['asset-kind', `asset-kind--${asset.kind}`]">{{ asset.kind === 'video' ? '视频' : '图片' }}</span>
              <span class="asset-name">{{ asset.name }}</span>
              <small>{{ formatSize(asset.size) }} · {{ asset.uploadStatus === 'uploaded' ? `已上传 #${asset.mediaId}` : '本地草稿' }}</small>
            </div>
          </div>
        </section>

        <section class="editor-card cover-card">
          <div class="editor-card__header">
            <div>
              <h2>设置封面</h2>
              <p>默认取第一份素材作为封面，可以切换推荐素材。</p>
            </div>
            <label class="switch-row">
              <span>PK 封面</span>
              <input v-model="pkCover" type="checkbox">
            </label>
          </div>

          <div class="cover-list">
            <button
              v-for="asset in coverAssets"
              :key="asset.assetId"
              type="button"
              :class="{ 'is-active': currentDraft.coverAssetId === asset.assetId }"
              @click="updateDraft({ coverAssetId: asset.assetId })"
            >
              <img v-if="asset.kind === 'image'" :src="asset.url" alt="">
              <video v-else :src="asset.url" muted playsinline />
              <span>推荐封面</span>
            </button>
          </div>
        </section>

        <section class="editor-card compose-card">
          <input
            v-model.trim="form.title"
            class="title-input"
            maxlength="50"
            placeholder="填写标题会有更多赞哦"
            @blur="persistForm"
          >
          <textarea
            v-model.trim="form.body"
            maxlength="1000"
            placeholder="输入正文描述，真诚有价值的分享予人温暖"
            @blur="persistForm"
          />

          <div class="topic-row">
            <button v-for="topic in defaultTopics" :key="topic" type="button" @click="toggleTopic(topic)">
              #{{ topic }}
            </button>
            <button type="button">更多⌄</button>
          </div>

          <div class="compose-actions">
            <button type="button"># 话题</button>
            <button type="button">@ 用户</button>
            <button type="button">☺ 表情</button>
            <span>{{ form.body.length }}/1000</span>
          </div>
        </section>

        <section class="editor-card">
          <div class="editor-card__header">
            <h2>活动话题</h2>
            <button class="text-button" type="button">更多 ›</button>
          </div>
          <div class="activity-grid">
            <div v-for="activity in activityTopics" :key="activity.title" class="activity-item">
              <span>{{ activity.image }}</span>
              <strong>{{ activity.title }}</strong>
              <small>{{ activity.action }}</small>
              <button type="button">活动详情 ›</button>
            </div>
          </div>
        </section>

        <section class="editor-card settings-card">
          <div class="editor-card__header">
            <h2>内容设置</h2>
            <button class="text-button" type="button">收起⌃</button>
          </div>

          <div class="setting-list">
            <button class="setting-line is-disabled" type="button">
              <span>☷</span>
              <strong>添加章节</strong>
              <small>视频时长不足15秒，不支持添加章节</small>
            </button>
            <button class="setting-line" type="button" @click="collectionOpen = true">
              <span>▱</span>
              <strong>{{ form.collection || '加入合集' }}</strong>
              <small>汇集系列视频，有利于连续观看</small>
            </button>
            <button class="setting-line" type="button">
              <span>☏</span>
              <strong>引用笔记</strong>
              <small>关联已有作品</small>
            </button>
            <label class="setting-line">
              <span>♢</span>
              <strong>原创声明</strong>
              <input v-model="form.original" type="checkbox" @change="persistForm">
            </label>
          </div>

          <h3>添加组件</h3>
          <div class="setting-grid">
            <button type="button">⌖ 添加地点</button>
            <button type="button">♙ 选择群聊</button>
            <button type="button">◇ 标记地点或标记朋友</button>
            <button type="button">↝ 添加路线</button>
          </div>
        </section>

        <section class="editor-card more-card">
          <div class="editor-card__header">
            <h2>更多设置</h2>
            <button class="text-button" type="button">收起⌃</button>
          </div>
          <label class="select-line">
            <span>公开可见</span>
            <select v-model="form.visibility" @change="persistForm">
              <option value="public">公开可见</option>
              <option value="fans">粉丝可见</option>
              <option value="private">仅自己可见</option>
            </select>
          </label>
          <label class="select-line">
            <span>定时发布</span>
            <input v-model="form.scheduled" type="checkbox" @change="persistForm">
          </label>
        </section>
      </div>

      <aside class="preview-panel">
        <div class="preview-tabs">
          <button class="is-active" type="button">笔记预览</button>
          <button type="button">封面预览</button>
        </div>

        <div class="phone-frame">
          <div class="phone-screen">
            <div class="phone-status">9:41 <span>▮▮ ◒ ▰</span></div>
            <div class="phone-media">
              <video v-if="previewAsset?.kind === 'video'" :src="previewAsset.url" controls muted playsinline />
              <img v-else-if="previewAsset" :src="previewAsset.url" alt="">
              <div v-else class="empty-preview">暂无素材</div>
            </div>
            <div class="phone-copy">
              <h3>{{ form.title || '项目工程' }}</h3>
              <p>{{ form.body || '提示词 提示词 提示词' }}</p>
              <div class="phone-user">
                <span>咕</span>
                <strong>咕噜创作者</strong>
                <button type="button">关注</button>
              </div>
            </div>
            <div class="phone-actions">
              <span>说点什么...</span>
              <strong>♡ 点赞</strong>
              <strong>☆ 收藏</strong>
              <strong>☏ 评论</strong>
            </div>
          </div>
        </div>
      </aside>

      <div class="publish-bar">
        <button class="ghost-button" type="button" @click="saveAndLeave">暂存离开</button>
        <button class="primary-button" type="button" @click="publishDraft">发布</button>
      </div>
    </div>

    <div v-if="uploading" class="upload-progress">
      <strong>正在上传素材</strong>
      <span>{{ uploadMessage }}</span>
    </div>

    <div v-if="draftDrawerOpen" class="drawer-mask" @click.self="draftDrawerOpen = false">
      <aside class="draft-drawer">
        <div class="drawer-header">
          <h2>草稿箱</h2>
          <button type="button" @click="draftDrawerOpen = false">×</button>
        </div>

        <div class="drawer-tabs">
          <button
            v-for="item in drawerTabs"
            :key="item.key"
            type="button"
            :class="{ 'is-active': drawerType === item.key }"
            @click="drawerType = item.key"
          >
            {{ item.label }}({{ item.count }})
          </button>
        </div>

        <p class="drawer-tip">草稿存储于当前浏览器本地，清除浏览器数据时会被删除；最多保存100篇草稿。</p>

        <div class="draft-list">
          <article v-for="draft in filteredDrafts" :key="draft.id" class="draft-item">
            <button class="draft-thumb" type="button" @click="editDraft(draft.id)">
              <span>{{ draft.type === 'video' ? '视频' : draft.type === 'mixed' ? '混合' : '图文' }}</span>
              <small>{{ draft.assetMetas?.length || 0 }} 个素材</small>
            </button>
            <button class="draft-info" type="button" @click="editDraft(draft.id)">
              <strong>{{ draft.title || '暂无笔记标题' }}</strong>
              <span>保存于{{ draft.savedText }}</span>
            </button>
            <button class="draft-action" type="button" @click="editDraft(draft.id)">编辑</button>
            <button class="draft-action" type="button" @click="deleteDraft(draft.id)">删除</button>
          </article>

          <p v-if="filteredDrafts.length === 0" class="empty-drafts">暂无草稿</p>
        </div>
      </aside>
    </div>

    <div v-if="collectionOpen" class="collection-popover" @click.self="collectionOpen = false">
      <div>
        <h3>添加视频合集</h3>
        <input v-model.trim="collectionName" placeholder="输入合集名称，例如：项目工程日志">
        <button type="button" @click="applyCollection">创建并加入合集</button>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'

import { completeUpload, createUploadPresign, uploadToObjectStorage } from '@/api/media'
import { usePublishDraftStore } from '@/stores/publishDrafts'

import { activityTopics, defaultTopics, publishTypes } from './publishOptions'

const draftStore = usePublishDraftStore()

const activeType = ref('video')
const drawerType = ref('video')
const dragging = ref(false)
const draftDrawerOpen = ref(false)
const collectionOpen = ref(false)
const collectionName = ref('')
const pkCover = ref(false)
const uploading = ref(false)
const uploadMessage = ref('')

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

async function createDraftFromFiles(fileList, draftType = activeType.value) {
  const files = Array.from(fileList).filter((file) => file.type.startsWith('video/') || file.type.startsWith('image/')).slice(0, 20)

  if (!files.length) {
    return
  }

  uploading.value = true
  uploadMessage.value = `准备上传 ${files.length} 个素材`

  try {
    const uploadedAssets = []

    for (const [index, file] of files.entries()) {
      uploadMessage.value = `正在上传 ${index + 1}/${files.length}：${file.name}`
      const presign = await createUploadPresign({
        resourceType: file.type.startsWith('video/') ? 'Video' : 'Image',
        fileName: file.name,
        contentType: 'application/octet-stream'
      })

      await uploadToObjectStorage(file, presign)

      const completed = await completeUpload({
        mediaId: presign.mediaId,
        objectKey: presign.objectKey
      })

      uploadedAssets.push({
        mediaId: completed.mediaId,
        bucket: completed.bucket,
        objectKey: completed.objectKey,
        previewUrl: completed.previewUrl || presign.previewUrl,
        status: completed.status
      })
    }

    await draftStore.createDraft(draftType, files, uploadedAssets)
  } finally {
    uploading.value = false
    uploadMessage.value = ''
  }
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

function formatSize(size) {
  if (size >= 1024 * 1024 * 1024) {
    return `${(size / 1024 / 1024 / 1024).toFixed(2)}GB`
  }

  if (size >= 1024 * 1024) {
    return `${(size / 1024 / 1024).toFixed(1)}MB`
  }

  return `${Math.max(1, Math.round(size / 1024))}KB`
}
</script>

<style scoped>
.publish-page {
  min-height: calc(100vh - 140px);
  color: #1f2937;
}

.upload-stage {
  display: grid;
  gap: 18px;
}

.publish-tabs {
  display: flex;
  align-items: center;
  gap: 18px;
  border-radius: 18px;
  background: #fff;
  padding: 12px 18px;
}

.publish-tabs button {
  min-width: 96px;
  height: 38px;
  border: 0;
  border-radius: 19px;
  background: transparent;
  color: #9ca3af;
  cursor: pointer;
  font-weight: 700;
}

.publish-tabs button.is-active {
  background: #f8f8f8;
  color: #333;
}

.publish-tabs .draft-entry {
  margin-left: auto;
  background: #f7f7f7;
  color: #333;
}

.upload-dropzone {
  display: grid;
  min-height: 560px;
  place-items: center;
  align-content: center;
  gap: 12px;
  border: 1px dashed transparent;
  border-radius: 16px;
  background: #fafafa;
  cursor: pointer;
  transition: border-color 0.2s ease, background 0.2s ease;
}

.upload-dropzone.is-dragging {
  border-color: var(--color-primary);
  background: #fff6f7;
}

.upload-dropzone input {
  position: absolute;
  width: 1px;
  height: 1px;
  opacity: 0;
}

.upload-icon {
  display: grid;
  width: 110px;
  height: 84px;
  place-items: center;
  border-radius: 42px;
  background: #e5e7eb;
  color: #fff;
  font-size: 54px;
  font-weight: 900;
  line-height: 1;
}

.upload-dropzone strong,
.upload-dropzone small {
  color: #8c939f;
  font-size: 14px;
}

.upload-dropzone button {
  height: 42px;
  border: 0;
  border-radius: 21px;
  background: #ff2442;
  color: #fff;
  cursor: pointer;
  font-weight: 800;
  padding: 0 28px;
}

.upload-tips {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 1px;
  border-radius: 16px;
  background: #fff;
  padding: 20px 0;
}

.upload-tips div {
  display: grid;
  gap: 8px;
  padding: 0 42px;
}

.upload-tips div + div {
  border-left: 1px solid #edf0f5;
}

.upload-tips strong {
  color: #6b7280;
  font-size: 13px;
}

.upload-tips span {
  color: #a0a7b2;
  font-size: 12px;
  line-height: 1.6;
}

.editor-stage {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 410px;
  gap: 28px;
  padding-bottom: 92px;
}

.editor-main {
  display: grid;
  gap: 14px;
}

.editor-card {
  border-radius: 18px;
  background: #fff;
  padding: 24px;
}

.editor-card__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 18px;
}

.editor-card h2,
.editor-card h3 {
  margin: 0;
  color: #1f2937;
  font-size: 18px;
}

.editor-card p {
  margin: 8px 0 0;
  color: #9ca3af;
  font-size: 13px;
}

.text-button,
.reupload-button {
  border: 0;
  background: transparent;
  color: #8b93a1;
  cursor: pointer;
  font-size: 13px;
}

.reupload-button input {
  display: none;
}

.asset-summary {
  display: grid;
  gap: 10px;
  border-radius: 12px;
  background: #fafafa;
  padding: 16px;
}

.asset-item {
  display: grid;
  grid-template-columns: 46px minmax(0, 1fr) auto;
  align-items: center;
  gap: 12px;
  color: #4b5563;
}

.asset-kind {
  display: grid;
  height: 26px;
  place-items: center;
  border-radius: 7px;
  background: #eef2ff;
  color: #4f46e5;
  font-size: 12px;
  font-weight: 800;
}

.asset-kind--image {
  background: #ecfdf3;
  color: #079455;
}

.asset-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.asset-item small {
  color: #a0a7b2;
}

.switch-row,
.setting-line,
.select-line {
  display: grid;
  align-items: center;
}

.switch-row {
  grid-template-columns: auto auto;
  gap: 8px;
  color: #6b7280;
  font-size: 13px;
}

.cover-list {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
}

.cover-list button {
  position: relative;
  overflow: hidden;
  aspect-ratio: 16 / 9;
  border: 2px solid transparent;
  border-radius: 10px;
  background: #f3f4f6;
  cursor: pointer;
  padding: 0;
}

.cover-list button.is-active {
  border-color: #ff2442;
}

.cover-list img,
.cover-list video {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.cover-list span {
  position: absolute;
  right: 8px;
  bottom: 8px;
  border-radius: 9px;
  background: rgba(17, 24, 39, 0.72);
  color: #fff;
  font-size: 11px;
  padding: 3px 8px;
}

.compose-card {
  display: grid;
  gap: 18px;
}

.title-input,
.compose-card textarea {
  width: 100%;
  border: 0;
  outline: 0;
}

.title-input {
  color: #1f2937;
  font-size: 20px;
  font-weight: 800;
}

.compose-card textarea {
  min-height: 154px;
  resize: vertical;
  color: #4b5563;
  font-size: 15px;
  line-height: 1.8;
}

.topic-row,
.compose-actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
}

.topic-row button,
.compose-actions button {
  min-height: 32px;
  border: 0;
  border-radius: 16px;
  background: #f7f8fa;
  color: #8b93a1;
  cursor: pointer;
  padding: 0 14px;
}

.compose-actions span {
  margin-left: auto;
  color: #b8bec9;
  font-size: 13px;
}

.activity-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.activity-item {
  display: grid;
  grid-template-columns: 52px minmax(0, 1fr) auto;
  align-items: center;
  gap: 12px;
  border-radius: 12px;
  background: #fafafa;
  padding: 12px;
}

.activity-item span {
  display: grid;
  width: 52px;
  height: 52px;
  place-items: center;
  border-radius: 10px;
  background: #ff2442;
  color: #fff;
  font-weight: 900;
}

.activity-item strong {
  font-size: 15px;
}

.activity-item small {
  color: #9ca3af;
}

.activity-item button {
  border: 0;
  background: transparent;
  color: #8b93a1;
  cursor: pointer;
}

.setting-list,
.setting-grid,
.more-card {
  display: grid;
  gap: 12px;
}

.setting-list {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.setting-line {
  min-height: 58px;
  grid-template-columns: 28px minmax(0, 1fr) auto;
  gap: 10px;
  border: 0;
  border-radius: 10px;
  background: #f8f8f8;
  color: #4b5563;
  cursor: pointer;
  padding: 10px 14px;
  text-align: left;
}

.setting-line small {
  color: #9ca3af;
}

.setting-line.is-disabled {
  cursor: not-allowed;
  opacity: 0.55;
}

.settings-card h3 {
  margin: 22px 0 14px;
  font-size: 16px;
}

.setting-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.setting-grid button,
.select-line {
  min-height: 48px;
  border: 0;
  border-radius: 10px;
  background: #f8f8f8;
  color: #6b7280;
  padding: 0 14px;
}

.setting-grid button {
  cursor: pointer;
  text-align: left;
}

.select-line {
  grid-template-columns: minmax(0, 1fr) auto;
}

.select-line select {
  border: 0;
  background: transparent;
  color: #6b7280;
  outline: 0;
}

.preview-panel {
  position: sticky;
  top: 96px;
  height: max-content;
}

.preview-tabs {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 4px;
  margin-bottom: 18px;
  border-radius: 20px;
  background: #eceef2;
  padding: 4px;
}

.preview-tabs button {
  height: 34px;
  border: 0;
  border-radius: 17px;
  background: transparent;
  color: #8b93a1;
}

.preview-tabs button.is-active {
  background: #fff;
  color: #1f2937;
  font-weight: 800;
}

.phone-frame {
  width: 330px;
  height: 684px;
  margin: 0 auto;
  border: 8px solid #4b4b4b;
  border-radius: 42px;
  background: #111;
  padding: 8px;
}

.phone-screen {
  position: relative;
  overflow: hidden;
  width: 100%;
  height: 100%;
  border-radius: 32px;
  background: #111;
  color: #fff;
}

.phone-status {
  position: absolute;
  z-index: 2;
  top: 18px;
  right: 22px;
  left: 22px;
  display: flex;
  justify-content: space-between;
  font-weight: 800;
}

.phone-media {
  position: absolute;
  inset: 0;
  display: grid;
  place-items: center;
  background: #222;
}

.phone-media video,
.phone-media img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.empty-preview {
  color: #8b93a1;
}

.phone-copy {
  position: absolute;
  right: 0;
  bottom: 76px;
  left: 0;
  padding: 0 20px;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.55);
}

.phone-copy h3 {
  margin: 0 0 12px;
  color: #fff100;
  font-size: 22px;
}

.phone-copy p {
  display: -webkit-box;
  overflow: hidden;
  margin: 0 0 18px;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 4;
  font-size: 16px;
  line-height: 1.55;
}

.phone-user {
  display: flex;
  align-items: center;
  gap: 8px;
}

.phone-user span {
  display: grid;
  width: 32px;
  height: 32px;
  place-items: center;
  border-radius: 50%;
  background: #ff2442;
}

.phone-user button {
  height: 28px;
  border: 0;
  border-radius: 14px;
  background: #ff2442;
  color: #fff;
  font-weight: 800;
  padding: 0 12px;
}

.phone-actions {
  position: absolute;
  right: 12px;
  bottom: 14px;
  left: 12px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.phone-actions span {
  flex: 1;
  min-width: 0;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.12);
  color: #cfd3dc;
  padding: 9px 12px;
}

.phone-actions strong {
  font-size: 12px;
}

.publish-bar {
  position: fixed;
  z-index: 18;
  right: 0;
  bottom: 0;
  left: 232px;
  display: flex;
  justify-content: center;
  gap: 20px;
  background: linear-gradient(180deg, rgba(245, 246, 248, 0), rgba(245, 246, 248, 0.96) 38%, #f5f6f8);
  padding: 28px 24px 24px;
}

.ghost-button,
.primary-button {
  width: 132px;
  height: 44px;
  border: 0;
  border-radius: 22px;
  cursor: pointer;
  font-weight: 800;
}

.ghost-button {
  background: #fff;
  color: #6b7280;
}

.primary-button {
  background: #ff2442;
  color: #fff;
}

.drawer-mask,
.collection-popover {
  position: fixed;
  z-index: 60;
  inset: 0;
  background: rgba(17, 24, 39, 0.42);
}

.upload-progress {
  position: fixed;
  z-index: 70;
  right: 28px;
  bottom: 28px;
  display: grid;
  gap: 6px;
  min-width: 260px;
  border-radius: 12px;
  background: #111827;
  color: #fff;
  box-shadow: 0 18px 50px rgba(15, 23, 42, 0.26);
  padding: 16px 18px;
}

.upload-progress strong {
  font-size: 14px;
}

.upload-progress span {
  overflow: hidden;
  color: #d1d5db;
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.draft-drawer {
  position: absolute;
  top: 16px;
  right: 16px;
  bottom: 16px;
  width: min(620px, calc(100vw - 32px));
  overflow: auto;
  border-radius: 18px;
  background: #fff;
  padding: 24px;
}

.drawer-header,
.drawer-tabs {
  display: flex;
  align-items: center;
}

.drawer-header {
  justify-content: space-between;
}

.drawer-header h2 {
  margin: 0;
  font-size: 20px;
}

.drawer-header button {
  border: 0;
  background: transparent;
  color: #6b7280;
  cursor: pointer;
  font-size: 28px;
}

.drawer-tabs {
  gap: 22px;
  margin-top: 20px;
  border-bottom: 1px solid #edf0f5;
}

.drawer-tabs button {
  height: 40px;
  border: 0;
  border-bottom: 2px solid transparent;
  background: transparent;
  color: #9ca3af;
  cursor: pointer;
  font-weight: 700;
}

.drawer-tabs button.is-active {
  border-color: #ff2442;
  color: #1f2937;
}

.drawer-tip {
  margin: 18px 0;
  border-radius: 0;
  background: #eef3ff;
  color: #475467;
  font-size: 13px;
  line-height: 1.7;
  padding: 12px 16px;
}

.draft-list {
  display: grid;
  gap: 14px;
}

.draft-item {
  display: grid;
  grid-template-columns: 82px minmax(0, 1fr) auto auto;
  align-items: center;
  gap: 14px;
  border-radius: 10px;
  background: #fafafa;
  padding: 12px;
}

.draft-thumb {
  display: grid;
  width: 82px;
  height: 82px;
  place-items: center;
  border: 0;
  border-radius: 8px;
  background: linear-gradient(145deg, #2f3137, #111);
  color: #fff;
  cursor: pointer;
}

.draft-thumb span {
  color: #fff100;
  font-weight: 900;
}

.draft-info {
  display: grid;
  gap: 8px;
  border: 0;
  background: transparent;
  cursor: pointer;
  text-align: left;
}

.draft-info strong {
  color: #4b5563;
  font-size: 15px;
}

.draft-info span {
  color: #9ca3af;
  font-size: 13px;
}

.draft-action {
  border: 0;
  background: transparent;
  color: #6b7280;
  cursor: pointer;
}

.empty-drafts {
  margin: 80px 0;
  color: #9ca3af;
  text-align: center;
}

.collection-popover {
  display: grid;
  place-items: center;
}

.collection-popover > div {
  display: grid;
  width: min(420px, calc(100vw - 32px));
  gap: 16px;
  border-radius: 16px;
  background: #fff;
  padding: 24px;
}

.collection-popover h3 {
  margin: 0;
}

.collection-popover input {
  height: 42px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  outline: 0;
  padding: 0 12px;
}

.collection-popover button {
  height: 42px;
  border: 0;
  border-radius: 21px;
  background: #ff2442;
  color: #fff;
  cursor: pointer;
  font-weight: 800;
}

@media (max-width: 1180px) {
  .editor-stage {
    grid-template-columns: 1fr;
  }

  .preview-panel {
    position: static;
  }
}

@media (max-width: 900px) {
  .publish-bar {
    left: 0;
  }
}

@media (max-width: 720px) {
  .publish-tabs,
  .upload-tips,
  .setting-list,
  .setting-grid,
  .activity-grid,
  .cover-list {
    grid-template-columns: 1fr;
  }

  .publish-tabs {
    display: grid;
    gap: 8px;
  }

  .publish-tabs .draft-entry {
    margin-left: 0;
  }

  .upload-dropzone {
    min-height: 360px;
  }

  .upload-tips div + div {
    border-left: 0;
    border-top: 1px solid #edf0f5;
    padding-top: 18px;
  }

  .editor-card {
    padding: 18px;
  }

  .draft-item {
    grid-template-columns: 68px 1fr;
  }

  .draft-action {
    text-align: left;
  }

  .phone-frame {
    width: 290px;
    height: 600px;
  }
}
</style>
