<script setup>
import { computed, ref, watch } from 'vue'
import { Button as VanButton, Popup as VanPopup } from 'vant'

import Avatar from '@/components/Avatar/avatar.vue'
import AuthorWrapper from '@/components/AuthorWrapper/authWrapper.vue'
import CommentItem from '@/components/CommentItem/commentItem.vue'
import LikeAction from '@/components/LikeAction/likeAction.vue'

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

const emit = defineEmits(['update:show', 'login-request', 'open-author'])
const activeIndex = ref(0)
const fullscreen = ref(false)
const imageRatio = ref(0.75)

const images = computed(() => {
  if (!props.note?.image) return []
  return [
    props.note.image,
    props.note.image,
    props.note.image,
    props.note.image
  ]
})
const comments = computed(() => [
  {
    id: 1,
    author: {
      nickname: props.note?.author || '地主页硬笔草书',
      avatar: props.note?.avatar || '',
      isAuthor: true
    },
    content: '莫等闲，白了少年头，空悲切。这里先放一条置顶评论，后续接评论接口即可替换。',
    meta: '5天前 辽宁',
    likes: 12,
    replies: [
      {
        id: '1-1',
        author: { nickname: '问一问', avatar: '', isAuthor: false },
        content: '回复一下作者：这段写得很有感觉。',
        meta: '4小时前 上海',
        likes: 1,
        replies: []
      },
      {
        id: '1-2',
        author: { nickname: '麦田里的倾听者', avatar: props.note?.avatar || '', isAuthor: true },
        content: '谢谢喜欢，后面会继续整理。',
        meta: '3小时前 上海',
        likes: 2,
        replies: []
      }
    ]
  },
  {
    id: 2,
    author: {
      nickname: '那年鲜衣怒马',
      avatar: '',
      isAuthor: false
    },
    content: '老师厉害，和我理解的有出入吗？',
    meta: '昨天 17:57 湖北',
    likes: 0,
    image: 'https://picsum.photos/id/64/180/180',
    replies: [
      {
        id: '2-1',
        author: { nickname: '问一问', avatar: '', isAuthor: false },
        content: '感觉主要是表达方式不一样。',
        meta: '48分钟前 北京',
        likes: 0,
        replies: []
      }
    ]
  },
  {
    id: 3,
    author: {
      nickname: '生活记录者',
      avatar: '',
      isAuthor: false
    },
    content: '确实，图片细节放大看更有味道。',
    meta: '昨天 07:43 江西',
    likes: 2,
    replies: []
  },
  {
    id: 4,
    author: {
      nickname: '口口',
      avatar: '',
      isAuthor: false
    },
    content: '文字是知识和文化的载体，这种内容值得慢慢看。',
    meta: '昨天 01:22 广东',
    likes: 2,
    replies: []
  }
])

watch(() => props.note?.id, () => {
  activeIndex.value = 0
  fullscreen.value = false
  imageRatio.value = 0.75
})

watch(activeIndex, () => {
  imageRatio.value = 0.75
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

function updateImageRatio(event) {
  const image = event.target
  if (!image?.naturalWidth || !image?.naturalHeight) return
  imageRatio.value = image.naturalWidth / image.naturalHeight
}

function openAuthor(author) {
  emit('open-author', {
    authorId: props.note?.authorId || props.note?.id || 'mock',
    author: author?.nickname || props.note?.author || '咕噜用户',
    avatar: author?.avatar || props.note?.avatar || ''
  })
}
</script>

<template>
  <VanPopup
    :show="props.show"
    class="image-work-popup"
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
        <span class="image-work__count">{{ activeIndex + 1 }}/{{ images.length }}</span>
        <button
          class="image-work__arrow is-left"
          type="button"
          aria-label="上一张"
          @click="prevImage"
        >
          ‹
        </button>
        <img
          :src="images[activeIndex]"
          :alt="note.title"
          @load="updateImageRatio"
          @dblclick="fullscreen = true"
        >
        <button
          class="image-work__arrow is-right"
          type="button"
          aria-label="下一张"
          @click="nextImage"
        >
          ›
        </button>
        <div class="image-work__dots">
          <button
            v-for="(_, index) in images"
            :key="index"
            type="button"
            :class="{ active: index === activeIndex }"
            :aria-label="`切换到第${index + 1}张`"
            @click="activeIndex = index"
          />
        </div>
      </section>

      <aside class="image-work__side">
        <header class="image-work__author">
          <Avatar
            :avatar-url="note.avatar"
            :username="note.author"
            :size="38"
            interactive
            @click="openAuthor({ nickname: note.author, avatar: note.avatar })"
          />
          <AuthorWrapper
            class="image-work__author-name"
            :nickname="note.author"
            @click="openAuthor({ nickname: note.author, avatar: note.avatar })"
          />
          <VanButton
            type="primary"
            round
            size="small"
            @click="$emit('login-request')"
          >
            关注
          </VanButton>
        </header>

        <section class="image-work__content">
          <h2>{{ note.title }}</h2>
          <p>
            什么样的内容会让人停下来认真看？先把作品详情和评论交互样式搭好，后续接接口时替换真实正文即可。
          </p>
          <div class="image-work__topics">
            <span>#硬笔书法</span>
            <span>#笔记</span>
            <span>#生活记录</span>
            <span>#灵感收藏</span>
          </div>
          <div class="image-work__meta">
            编辑于 5天前 辽宁
            <button type="button">
              ···
            </button>
          </div>
        </section>

        <section class="image-work__comments">
          <p class="image-work__comment-count">
            共 191 条评论
          </p>
          <CommentItem
            v-for="comment in comments"
            :key="comment.id"
            :comment="comment"
            :avatar-size="34"
            @open-author="openAuthor($event.author)"
          />
        </section>

        <footer class="image-work__actions">
          <button type="button">
            说点什么...
          </button>
          <span>
            <LikeAction
              :count="note.likes"
              label="点赞作品"
              icon-size="21"
            />
          </span>
          <span>☆ 306</span>
          <span>💬 191</span>
          <span>↗</span>
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
      :alt="note?.title"
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

.image-work__viewer img {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: contain;
  cursor: zoom-in;
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
  grid-template-columns: 38px minmax(0, 1fr) auto;
  align-items: center;
  gap: 12px;
  padding: 24px 22px 16px;
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
  scrollbar-gutter: stable;
  padding: 14px 22px 84px;
}

.image-work__comment-count {
  margin: 0 0 16px;
  color: #8a8f99;
  font-size: 13px;
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

.image-work__actions > button {
  height: 36px;
  border-radius: 999px;
  color: #9da3ad;
  background: #f6f7f8;
  text-align: left;
  padding: 0 18px;
}

.image-work__actions span {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  color: #424751;
  white-space: nowrap;
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
