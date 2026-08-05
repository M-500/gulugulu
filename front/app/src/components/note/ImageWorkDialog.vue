<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { Button as VanButton, Popup as VanPopup, showToast } from 'vant'

import Avatar from '@/components/Avatar/avatar.vue'
import AuthorWrapper from '@/components/AuthorWrapper/authWrapper.vue'
import CommentItem from '@/components/CommentItem/commentItem.vue'
import SvgIcon from '@/components/common/SvgIcon.vue'
import HlsVideoPlayer from '@/components/HlsVideoPlayer.vue'
import LikeAction from '@/components/LikeAction/likeAction.vue'
import { createComment, getWorkComments, uploadCommentImage } from '@/services/commentService'
import { getAppWorkDetail } from '@/services/workService'
import { useAuthStore } from '@/stores/auth'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  note: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['update:show', 'login-request', 'open-author', 'like-change'])
const activeIndex = ref(0)
const fullscreen = ref(false)
const imageRatio = ref(0.75)
const detail = ref(null)
const loading = ref(false)
const loadError = ref('')
const commentInput = ref(null)
const composerOpen = ref(false)
const commentDraft = ref('')
const replyContext = ref(null)
const comments = ref([])
const commentsPage = ref(1)
const commentsHasMore = ref(false)
const commentsLoading = ref(false)
const commentsError = ref('')
const localCommentCount = ref(0)
const commentImageInput = ref(null)
const commentImageFile = ref(null)
const commentImagePreview = ref('')
const commentSubmitting = ref(false)
const favorited = ref(false)
const favoriteDelta = ref(0)
const authStore = useAuthStore()
let requestVersion = 0
let commentRequestVersion = 0

const work = computed(() => {
  if (detail.value) return detail.value
  return {
    id: props.note?.id,
    type: props.note?.isVideo ? 'video' : 'image',
    title: props.note?.title || '',
    content: '',
    cover: props.note?.image || '',
    videoPlaylist: '',
    publishedAt: '',
    authorId: props.note?.authorId,
    author: props.note?.author || '咕噜用户',
    avatar: props.note?.avatar || '',
    likes: props.note?.likes || 0,
    favoriteCount: 306,
    commentCount: 0,
    shareCount: 26,
    assets: [],
    topics: []
  }
})

const isVideo = computed(() => work.value.type === 'video')
const images = computed(() => {
  const values = work.value.assets
    ?.filter((asset) => asset.role === 'image' && asset.url)
    .sort((a, b) => a.sort - b.sort)
    .map((asset) => asset.url) || []
  if (values.length) return values
  return work.value.cover ? [work.value.cover] : []
})
const displayComments = computed(() => comments.value)
const displayCommentCount = computed(() => Number(work.value.commentCount || 0) + localCommentCount.value)
const displayFavoriteCount = computed(() => Math.max(
  0,
  Number(work.value.favoriteCount || 0) + favoriteDelta.value
))
const canSubmitComment = computed(() => (
  Boolean(commentDraft.value.trim() || commentImageFile.value) && !commentSubmitting.value
))
const replyAuthorName = computed(() => (
  replyContext.value?.target?.author?.nickname
  || replyContext.value?.target?.author?.username
  || '用户'
))
const commentPlaceholder = computed(() => (
  replyContext.value ? `回复 ${replyAuthorName.value}` : '说点什么...'
))

watch([() => props.show, () => props.note?.id], ([show, workId], previous) => {
  const previousWorkId = previous?.[1]
  if (workId !== previousWorkId) {
    detail.value = null
    loadError.value = ''
    resetComments()
    favorited.value = false
    favoriteDelta.value = 0
  }
  cancelCommentComposer()
  activeIndex.value = 0
  fullscreen.value = false
  imageRatio.value = props.note?.isVideo ? 16 / 9 : 0.75
  if (show && workId) {
    loadDetail(workId)
    loadComments(workId, true)
  }
})

watch(activeIndex, () => {
  updateRatioFromAsset()
})

watch(detail, () => {
  updateRatioFromAsset()
})

const viewerStyle = computed(() => ({
  '--image-ratio': imageRatio.value
}))

function close () {
  fullscreen.value = false
  emit('update:show', false)
}

function prevImage () {
  if (!images.value.length) return
  activeIndex.value = (activeIndex.value + images.value.length - 1) % images.value.length
}

function nextImage () {
  if (!images.value.length) return
  activeIndex.value = (activeIndex.value + 1) % images.value.length
}

function openCommentComposer() {
  replyContext.value = null
  composerOpen.value = true
  focusCommentInput()
}

function toggleFavorite() {
  favorited.value = !favorited.value
  favoriteDelta.value += favorited.value ? 1 : -1
  showToast(favorited.value ? '已收藏' : '已取消收藏')
}

function handleWorkLikeChange(result) {
  if (detail.value) {
    detail.value.liked = result.liked
    detail.value.likes = String(result.count)
  }
  emit('like-change', { note: props.note, result })
}

function openReplyComposer(context) {
  replyContext.value = context
  composerOpen.value = true
  focusCommentInput()
}

function cancelCommentComposer() {
  composerOpen.value = false
  commentDraft.value = ''
  replyContext.value = null
  clearCommentImage()
}

function focusCommentInput() {
  nextTick(() => commentInput.value?.focus())
}

function insertCommentText(value) {
  commentDraft.value += value
  focusCommentInput()
}

async function submitComment() {
  const content = commentDraft.value.trim()
  if (!content && !commentImageFile.value) return
  if (!authStore.isLoggedIn) {
    emit('login-request')
    return
  }

  commentSubmitting.value = true
  try {
    let imageMediaId = 0
    if (commentImageFile.value) {
      const uploaded = await uploadCommentImage(commentImageFile.value, authStore.accessToken)
      imageMediaId = uploaded.mediaId
    }
    const created = await createComment(work.value.id, {
      content,
      imageMediaId,
      parentCommentId: replyContext.value?.target?.id || 0
    }, authStore.accessToken)
    created.author.isAuthor = Number(created.author.id) === Number(work.value.authorId)
    if (replyContext.value) {
	  created.replyToName = replyAuthorName.value
      const rootId = replyContext.value.root.id
      const root = comments.value.find((item) => String(item.id) === String(rootId))
      if (root) {
        root.replies = [...(root.replies || []), created]
        root.replyCount = Number(root.replyCount || 0) + 1
      }
    } else {
      comments.value.unshift(created)
    }
    localCommentCount.value += 1
    showToast(replyContext.value ? '回复已发送' : '评论已发送')
    cancelCommentComposer()
  } catch (error) {
    showToast(error.message || '评论发布失败')
  } finally {
    commentSubmitting.value = false
  }
}

function resetComments() {
	commentRequestVersion += 1
  comments.value = []
  commentsPage.value = 1
  commentsHasMore.value = false
  commentsError.value = ''
  localCommentCount.value = 0
}

function selectCommentImage(event) {
  const file = event.target.files?.[0]
  event.target.value = ''
  if (!file) return
  if (!file.type.startsWith('image/')) {
    showToast('请选择图片文件')
    return
  }
  if (file.size > 10 * 1024 * 1024) {
    showToast('评论图片不能超过10MB')
    return
  }
  clearCommentImage()
  commentImageFile.value = file
  commentImagePreview.value = URL.createObjectURL(file)
}

function clearCommentImage() {
  if (commentImagePreview.value) URL.revokeObjectURL(commentImagePreview.value)
  commentImageFile.value = null
  commentImagePreview.value = ''
}

async function loadComments(workId, reset = false) {
  if ((!reset && commentsLoading.value) || !workId) return
	const currentVersion = reset ? ++commentRequestVersion : commentRequestVersion
  commentsLoading.value = true
  commentsError.value = ''
  try {
    const nextPage = reset ? 1 : commentsPage.value + 1
    const data = await getWorkComments(workId, { page: nextPage, pageSize: 10 })
	if (currentVersion !== commentRequestVersion) return
    for (const comment of data.list) {
      comment.author.isAuthor = Number(comment.author.id) === Number(work.value.authorId)
      for (const reply of comment.replies || []) {
        reply.author.isAuthor = Number(reply.author.id) === Number(work.value.authorId)
      }
    }
    comments.value = reset ? data.list : [...comments.value, ...data.list]
    commentsPage.value = data.page
    commentsHasMore.value = data.hasMore
  } catch (error) {
	if (currentVersion !== commentRequestVersion) return
    commentsError.value = error.message || '评论加载失败'
  } finally {
	if (currentVersion === commentRequestVersion) commentsLoading.value = false
  }
}

function isEditableTarget(target) {
  if (!(target instanceof HTMLElement)) return false
  return Boolean(target.closest('input, textarea, select')) || target.isContentEditable
}

// 普通详情和全屏预览共用同一套图片索引，保证键盘与按钮切换结果一致。
function handleImageKeydown(event) {
  if (!props.show || isVideo.value || images.value.length <= 1) return
  if (event.altKey || event.ctrlKey || event.metaKey || event.shiftKey) return
  if (!fullscreen.value && isEditableTarget(event.target)) return

  if (event.key === 'ArrowLeft') {
    event.preventDefault()
    event.stopPropagation()
    prevImage()
  } else if (event.key === 'ArrowRight') {
    event.preventDefault()
    event.stopPropagation()
    nextImage()
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleImageKeydown, true)
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleImageKeydown, true)
  clearCommentImage()
})

function updateImageRatio(event) {
  const image = event.target
  if (!image?.naturalWidth || !image?.naturalHeight) return
  imageRatio.value = image.naturalWidth / image.naturalHeight
}

function openAuthor(author) {
  emit('open-author', {
	authorId: author?.id || author?.userId || work.value.authorId || 'mock',
    author: author?.nickname || work.value.author || '咕噜用户',
    avatar: author?.avatar || work.value.avatar || ''
  })
}

async function loadDetail(workId) {
  const currentVersion = ++requestVersion
  loading.value = true
  loadError.value = ''
  try {
    const data = await getAppWorkDetail(workId)
    if (currentVersion === requestVersion) {
      detail.value = data
    }
  } catch (error) {
    if (currentVersion === requestVersion) {
      loadError.value = error.message || '作品详情加载失败'
    }
  } finally {
    if (currentVersion === requestVersion) {
      loading.value = false
    }
  }
}

function updateRatioFromAsset() {
  if (isVideo.value) {
    const video = work.value.assets?.find((asset) => asset.role === 'video')
    imageRatio.value = video?.width && video?.height ? video.width / video.height : 16 / 9
    return
  }
  const image = work.value.assets
    ?.filter((asset) => asset.role === 'image')
    .sort((a, b) => a.sort - b.sort)[activeIndex.value]
  imageRatio.value = image?.width && image?.height ? image.width / image.height : 0.75
}

function formatPublishedAt(value) {
  if (!value) return '刚刚发布'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}
</script>

<template>
  <VanPopup
    :show="props.show"
    class="image-work-popup"
    :class="{ 'is-video': isVideo }"
    :style="viewerStyle"
    overlay-class="image-work-overlay"
    teleport="body"
    @update:show="$emit('update:show', $event)"
    @click-overlay="close"
  >
    <article
      v-if="note"
      class="image-work"
    >
      <section class="image-work__viewer">
        <button
          class="image-work__close"
          type="button"
          aria-label="关闭"
          @click="close"
        >
          ×
        </button>
        <div
          v-if="loading"
          class="image-work__media-state"
        >
          正在加载作品详情...
        </div>
        <div
          v-else-if="loadError"
          class="image-work__media-state"
        >
          <span>{{ loadError }}</span>
          <button
            type="button"
            @click="loadDetail(note.id)"
          >
            重新加载
          </button>
        </div>
        <HlsVideoPlayer
          v-else-if="props.show && isVideo"
          :playlist="work.videoPlaylist"
          :poster="work.cover"
        />
        <template v-else>
          <span
            v-if="images.length > 1"
            class="image-work__count"
          >{{ activeIndex + 1 }}/{{ images.length }}</span>
          <button
            v-if="images.length > 1"
            class="image-work__arrow is-left"
            type="button"
            aria-label="上一张"
            @click="prevImage"
          >
            ‹
          </button>
          <img
            v-if="images.length"
            :src="images[activeIndex]"
            :alt="work.title"
            @load="updateImageRatio"
            @dblclick="fullscreen = true"
          >
          <div
            v-else
            class="image-work__media-state"
          >
            暂无可展示的图片
          </div>
          <button
            v-if="images.length > 1"
            class="image-work__arrow is-right"
            type="button"
            aria-label="下一张"
            @click="nextImage"
          >
            ›
          </button>
          <div
            v-if="images.length > 1"
            class="image-work__dots"
          >
            <button
              v-for="(_, index) in images"
              :key="index"
              type="button"
              :class="{ active: index === activeIndex }"
              :aria-label="`切换到第${index + 1}张`"
              @click="activeIndex = index"
            />
          </div>
        </template>
      </section>

      <aside class="image-work__side">
        <header class="image-work__author">
          <Avatar
            :avatar-url="work.avatar"
            :username="work.author"
            :size="40"
            interactive
            @click="openAuthor({ nickname: work.author, avatar: work.avatar })"
          />
          <AuthorWrapper
            class="image-work__author-name"
            :nickname="work.author"
            @click="openAuthor({ nickname: work.author, avatar: work.avatar })"
          />
          <VanButton
            class="image-work__follow"
            type="primary"
            round
            size="small"
            @click="$emit('login-request')"
          >
            关注
          </VanButton>
        </header>

        <section class="image-work__content">
          <h2>{{ work.title }}</h2>
          <p v-if="work.content">
            {{ work.content }}
          </p>
          <div
            v-if="work.topics?.length"
            class="image-work__topics"
          >
            <span
              v-for="topic in work.topics"
              :key="topic.topicId"
            >#{{ topic.name }}</span>
          </div>
          <div class="image-work__meta">
            发布于 {{ formatPublishedAt(work.publishedAt) }}
            <button type="button">
              ···
            </button>
          </div>
        </section>

        <section class="image-work__comments">
          <p class="image-work__comment-count">
            共 {{ displayCommentCount }} 条评论
          </p>
          <CommentItem
            v-for="comment in displayComments"
            :key="comment.id"
            :comment="comment"
            :avatar-size="34"
            @open-author="openAuthor($event.author)"
            @reply="openReplyComposer"
            @login-request="$emit('login-request')"
          />
          <p
            v-if="commentsLoading && !displayComments.length"
            class="image-work__comment-state"
          >
            正在加载评论...
          </p>
          <p
            v-else-if="commentsError && !displayComments.length"
            class="image-work__comment-state is-error"
          >
            {{ commentsError }}
          </p>
          <p
            v-else-if="!displayComments.length"
            class="image-work__comment-state"
          >
            还没有评论，来说点什么吧
          </p>
          <button
            v-if="commentsHasMore || (commentsError && displayComments.length)"
            class="image-work__load-comments"
            type="button"
            :disabled="commentsLoading"
            @click="loadComments(work.id)"
          >
            {{ commentsLoading ? '加载中...' : commentsError ? '重新加载' : '加载更多评论' }}
          </button>
        </section>

        <footer
          class="image-work__actions"
          :class="{ 'is-composing': composerOpen }"
        >
          <template v-if="!composerOpen">
            <button
              class="image-work__comment-trigger"
              type="button"
              @click="openCommentComposer"
            >
              <Avatar
                :avatar-url="authStore.user?.avatarUrl"
                :username="authStore.user?.nickName || '我'"
                :size="26"
              />
              <span>说点什么...</span>
            </button>
            <span class="image-work__stat">
              <LikeAction
                :count="work.likes"
                label="点赞作品"
                icon-size="21"
                :resource-id="work.id"
                resource-type="work"
                :model-value="Boolean(work.liked)"
                @login-request="$emit('login-request')"
                @change="handleWorkLikeChange"
              />
            </span>
            <button
              class="image-work__stat"
              :class="{ 'is-active': favorited }"
              type="button"
              :aria-label="favorited ? '取消收藏' : '收藏作品'"
              :aria-pressed="favorited"
              @click="toggleFavorite"
            >
              <SvgIcon
                name="collect"
                size="21"
              />
              {{ displayFavoriteCount }}
            </button>
            <button
              class="image-work__stat"
              type="button"
              aria-label="发表评论"
              @click="openCommentComposer"
            >
              <SvgIcon
                name="comment"
                size="21"
              />
              {{ displayCommentCount }}
            </button>
            <button
              class="image-work__stat image-work__stat--share"
              type="button"
              aria-label="分享作品"
            >
              <SvgIcon
                name="share"
                size="21"
              />
            </button>
          </template>

          <form
            v-else
            class="comment-composer"
            @submit.prevent="submitComment"
          >
            <div
              v-if="replyContext"
              class="comment-composer__reference"
            >
              <span>回复 {{ replyAuthorName }}</span>
              <p>{{ replyContext.target.content }}</p>
            </div>
            <div class="comment-composer__input-row">
              <input
                ref="commentInput"
                v-model="commentDraft"
                :placeholder="commentPlaceholder"
                maxlength="500"
                autocomplete="off"
                @keydown.esc.prevent="cancelCommentComposer"
              >
              <div class="comment-composer__quick-emoji">
                <button
                  v-for="emoji in ['🫠', '🥵', '😭']"
                  :key="emoji"
                  type="button"
                  :aria-label="`插入${emoji}`"
                  @click="insertCommentText(emoji)"
                >
                  {{ emoji }}
                </button>
              </div>
            </div>
            <div
              v-if="commentImagePreview"
              class="comment-composer__image-preview"
            >
              <img
                :src="commentImagePreview"
                alt="待发送的评论图片"
              >
              <button
                type="button"
                aria-label="移除评论图片"
                @click="clearCommentImage"
              >
                ×
              </button>
            </div>
            <div class="comment-composer__toolbar">
              <button
                type="button"
                aria-label="提醒用户"
                @click="insertCommentText('@')"
              >
                @
              </button>
              <button
                type="button"
                aria-label="添加评论图片"
                @click="commentImageInput?.click()"
              >
                ▧
              </button>
              <input
                ref="commentImageInput"
                class="comment-composer__file-input"
                type="file"
                accept="image/*"
                @change="selectCommentImage"
              >
              <span />
              <button
                class="comment-composer__send"
                type="submit"
                :disabled="!canSubmitComment"
              >
                {{ commentSubmitting ? '发送中' : '发送' }}
              </button>
              <button
                class="comment-composer__cancel"
                type="button"
                @click="cancelCommentComposer"
              >
                取消
              </button>
            </div>
          </form>
        </footer>
      </aside>
    </article>
  </VanPopup>

  <div
    v-if="fullscreen"
    class="image-fullscreen"
    @dblclick="fullscreen = false"
  >
    <button
      class="image-fullscreen__close"
      type="button"
      aria-label="退出全屏"
      @click="fullscreen = false"
    >
      ×
    </button>
    <button
      class="image-fullscreen__arrow is-left"
      type="button"
      aria-label="上一张"
      @click.stop="prevImage"
    >
      ‹
    </button>
    <img
      :src="images[activeIndex]"
      :alt="work.title"
    >
    <button
      class="image-fullscreen__arrow is-right"
      type="button"
      aria-label="下一张"
      @click.stop="nextImage"
    >
      ›
    </button>
  </div>
</template>

<style scoped>
.image-work-popup {
  --detail-side-width: 440px;
  --detail-viewer-height: 80vh;
  --detail-viewer-width: clamp(360px, calc(var(--detail-viewer-height) * var(--image-ratio)), calc(100vw - var(--detail-side-width) - 96px));
  width: min(calc(var(--detail-viewer-width) + var(--detail-side-width)), calc(100vw - 96px));
  max-width: calc(100vw - 96px);
  height: 80vh;
  overflow: hidden;
  border-radius: 18px;
  background: #fff;
}

.image-work {
  display: grid;
  height: 100%;
  width: 100%;
  grid-template-columns: minmax(0, var(--detail-viewer-width)) minmax(0, var(--detail-side-width));
}

.image-work__viewer {
  position: relative;
  display: grid;
  width: var(--detail-viewer-width);
  min-width: 0;
  place-items: center;
  overflow: hidden;
  background: #efece3;
}

.image-work-popup.is-video .image-work__viewer {
  background: #000;
}

.image-work__viewer img {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: contain;
  cursor: zoom-in;
}

.image-work__media-state {
  position: relative;
  z-index: 3;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: #6f747d;
  font-size: 14px;
  text-align: center;
}

.image-work__media-state button {
  height: 34px;
  border-radius: 999px;
  color: #fff;
  background: #3c4048;
  padding: 0 16px;
}

.image-work-popup.is-video .image-work__media-state {
  color: #fff;
}

.image-work-popup.is-video .image-work__media-state button {
  background: rgba(255, 255, 255, 0.18);
}

.image-work__close,
.image-work__count,
.image-work__arrow,
.image-work__dots {
  position: absolute;
  z-index: 2;
}

.image-work__close {
  top: 14px;
  left: 14px;
  display: none;
  width: 34px;
  height: 34px;
  border-radius: 50%;
  color: #fff;
  background: rgba(25, 28, 34, 0.45);
  font-size: 24px;
}

.image-work__count {
  top: 18px;
  right: 18px;
  border-radius: 999px;
  color: #fff;
  background: rgba(24, 25, 28, 0.42);
  font-size: 13px;
  padding: 5px 9px;
}

.image-work__arrow,
.image-fullscreen__arrow {
  top: 50%;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  color: #fff;
  background: rgba(37, 39, 45, 0.25);
  font-size: 30px;
  line-height: 1;
  transform: translateY(-50%);
}

.image-work__arrow.is-left,
.image-fullscreen__arrow.is-left {
  left: 18px;
}

.image-work__arrow.is-right,
.image-fullscreen__arrow.is-right {
  right: 18px;
}

.image-work__dots {
  right: 0;
  bottom: 16px;
  left: 0;
  display: flex;
  justify-content: center;
  gap: 7px;
}

.image-work__dots button {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.55);
}

.image-work__dots button.active {
  background: #fff;
}

.image-work__side {
  display: flex;
  min-width: 0;
  overflow: hidden;
  flex-direction: column;
  background: #fff;
}

.image-work__author {
  display: grid;
  grid-template-columns: 38px minmax(0, 1fr) 96px;
  align-items: center;
  gap: 12px;
  padding: 24px 22px 16px;
}

.image-work__follow {
  width: 100%;
}

.image-work__author-name {
  overflow: hidden;
  color: #30333a;
  font-size: 16px;
}

.image-work__author-name :deep(.author-wrapper__name) {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.image-work__content {
  border-bottom: 1px solid #f0f1f4;
  padding: 6px 22px 18px;
}

.image-work__content h2 {
  margin: 0 0 10px;
  color: #20242c;
  font-size: 17px;
  line-height: 1.5;
}

.image-work__content p {
  margin: 0;
  color: #343941;
  font-size: 14px;
  line-height: 1.8;
}

.image-work__topics {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 16px;
  color: #1f5595;
  font-size: 14px;
  line-height: 1.6;
}

.image-work__meta {
  display: flex;
  justify-content: space-between;
  margin-top: 13px;
  color: #a0a5ad;
  font-size: 12px;
}

.image-work__meta button {
  color: inherit;
}

.image-work__comments {
  min-height: 0;
  flex: 1;
  overflow-y: auto;
  overscroll-behavior: contain;
  scrollbar-width: none;
  -ms-overflow-style: none;
  padding: 14px 22px 84px;
}

.image-work__comments::-webkit-scrollbar {
  display: none;
  width: 0;
  height: 0;
}

.image-work__comment-count {
  margin: 0 0 16px;
  color: #8a8f99;
  font-size: 13px;
}

.image-work__comment-state {
  margin: 26px 0;
  color: #9ba1aa;
  font-size: 13px;
  text-align: center;
}

.image-work__comment-state.is-error {
  color: #d0525d;
}

.image-work__load-comments {
  display: block;
  margin: 8px auto 22px;
  color: #315b91;
  font-size: 13px;
}

.image-work__load-comments:disabled {
  opacity: 0.55;
}

.image-work__actions {
  position: sticky;
  right: 0;
  bottom: 0;
  left: 0;
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto auto auto auto;
  align-items: center;
  gap: 16px;
  border-top: 1px solid #eff1f4;
  background: #fff;
  padding: 12px 22px;
}

.image-work__comment-trigger {
  display: flex;
  min-width: 0;
  height: 36px;
  align-items: center;
  gap: 9px;
  border-radius: 999px;
  color: #9da3ad;
  background: #f6f7f8;
  text-align: left;
  padding: 0 12px 0 6px;
}

.image-work__comment-trigger > span:last-child {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.image-work__stat {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  border: 0;
  background: transparent;
  color: #424751;
  cursor: pointer;
  font: inherit;
  white-space: nowrap;
  padding: 0;
  transition: color 0.18s ease, transform 0.18s ease;
}

.image-work__stat:hover {
  color: var(--color-primary);
}

.image-work__stat:active {
  transform: scale(0.92);
}

.image-work__stat.is-active {
  color: var(--color-primary);
}

.image-work__stat--share {
  width: 22px;
  justify-content: center;
}

.image-work__actions.is-composing {
  display: block;
  padding: 13px 18px 12px;
}

.comment-composer__reference {
  margin: -1px 10px 10px;
  color: #7f858f;
  font-size: 12px;
  line-height: 1.5;
}

.comment-composer__reference span {
  display: block;
  margin-bottom: 2px;
}

.comment-composer__reference p {
  display: -webkit-box;
  overflow: hidden;
  margin: 0;
  color: #4d535c;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

.comment-composer__input-row {
  position: relative;
}

.comment-composer__input-row input {
  width: 100%;
  height: 40px;
  border: 1px solid transparent;
  border-radius: 999px;
  outline: none;
  background: #f6f7f8;
  color: #30333a;
  font-size: 14px;
  padding: 0 118px 0 14px;
  transition: border-color 0.2s, background 0.2s;
}

.comment-composer__input-row input:focus {
  border-color: #ffd6da;
  background: #fafafa;
}

.comment-composer__quick-emoji {
  position: absolute;
  top: 0;
  right: 11px;
  bottom: 0;
  display: flex;
  align-items: center;
  gap: 7px;
}

.comment-composer__quick-emoji button {
  display: inline-grid;
  width: 24px;
  height: 28px;
  place-items: center;
  border-radius: 7px;
  font-size: 18px;
}

.comment-composer__quick-emoji button:hover {
  background: #eceef1;
}

.comment-composer__toolbar {
  display: grid;
  grid-template-columns: 32px 32px minmax(0, 1fr) auto auto;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
}

.comment-composer__file-input {
  display: none;
}

.comment-composer__image-preview {
  position: relative;
  width: 72px;
  height: 72px;
  margin: 9px 0 0 8px;
  overflow: hidden;
  border-radius: 9px;
  background: #f3f4f5;
}

.comment-composer__image-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.comment-composer__image-preview button {
  position: absolute;
  top: 3px;
  right: 3px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  color: #fff;
  background: rgba(22, 24, 28, 0.65);
  line-height: 1;
}

.comment-composer__toolbar > button {
  height: 34px;
  border-radius: 999px;
  color: #575d66;
  font-size: 19px;
}

.comment-composer__toolbar .comment-composer__send,
.comment-composer__toolbar .comment-composer__cancel {
  min-width: 64px;
  padding: 0 17px;
  font-size: 14px;
  font-weight: 700;
}

.comment-composer__toolbar .comment-composer__send {
  color: #fff;
  background: #ff8993;
}

.comment-composer__toolbar .comment-composer__send:disabled {
  background: #ffc6cb;
  cursor: not-allowed;
}

.comment-composer__toolbar .comment-composer__cancel {
  border: 1px solid #e8e9ec;
  background: #fff;
  color: #6c7179;
}

.image-fullscreen {
  position: fixed;
  z-index: 3000;
  inset: 0;
  display: grid;
  place-items: center;
  background: rgba(0, 0, 0, 0.94);
}

.image-fullscreen img {
  max-width: 96vw;
  max-height: 96vh;
  object-fit: contain;
}

.image-fullscreen__close {
  position: absolute;
  top: 18px;
  right: 24px;
  color: #fff;
  font-size: 34px;
}

@media (max-width: 900px) {
  .image-work-popup {
    width: 100vw;
    height: 100vh;
    max-width: none;
    border-radius: 0;
  }

  .image-work {
    width: 100vw;
    max-width: none;
    grid-template-columns: 1fr;
    grid-template-rows: minmax(0, 58vh) minmax(0, 42vh);
  }

  .image-work__side {
    min-height: 0;
  }

  .image-work__viewer {
    width: 100%;
  }

  .image-work__close {
    display: inline-flex;
    align-items: center;
    justify-content: center;
  }

  .image-work__author {
    padding-top: 16px;
  }
}
</style>
