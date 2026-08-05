<template>
  <article class="comment-item">
    <Avatar
      :avatar-url="comment.author?.avatar"
      :username="comment.author?.nickname || comment.author?.username"
      :size="avatarSize"
      interactive
      @click="$emit('open-author', comment)"
    />
    <div class="comment-item__body">
      <p class="comment-item__author">
        <AuthorWrapper
          :nickname="comment.author?.nickname || comment.author?.username"
          :is-author="comment.author?.isAuthor"
          inline
          @click="$emit('open-author', comment)"
        />
      </p>
      <div
        v-if="comment.content"
        class="comment-item__content"
      >
        <span
          v-if="comment.replyToName"
          class="comment-item__reply-to"
        >回复 {{ comment.replyToName }}：</span>{{ comment.content }}
      </div>
      <div
        v-if="comment.image"
        class="comment-item__image"
      >
        <img
          :src="comment.image"
          :alt="comment.content"
        >
      </div>
      <footer class="comment-item__footer">
        <span>{{ comment.meta }}</span>
        <LikeAction
          :count="comment.likes"
          label="评论点赞"
          icon-size="14"
          :resource-id="comment.id"
          resource-type="comment"
          :model-value="Boolean(comment.liked)"
          @login-request="$emit('login-request')"
        />
        <button
          type="button"
          @click="$emit('reply', { target: comment, root: rootComment || comment })"
        >
          回复
        </button>
      </footer>
      <button
        v-if="replyTotal > 0"
        type="button"
        class="comment-item__expand"
        @click="toggleReplies"
      >
        {{ expanded ? '收起回复' : `展开 ${replyTotal} 条回复` }}
      </button>
      <div
        v-if="expanded"
        class="comment-item__replies"
      >
        <CommentItem
          v-for="reply in replies"
          :key="reply.id"
          :comment="reply"
          :root-comment="rootComment || comment"
          :avatar-size="28"
          @open-author="$emit('open-author', $event)"
          @reply="$emit('reply', $event)"
          @login-request="$emit('login-request')"
        />
        <p
          v-if="repliesLoading"
          class="comment-item__reply-state"
        >
          正在加载回复...
        </p>
        <button
          v-else-if="repliesHasMore"
          type="button"
          class="comment-item__load-more"
          @click="loadMoreReplies"
        >
          加载更多回复
        </button>
      </div>
    </div>
  </article>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { showToast } from 'vant'

import Avatar from '@/components/Avatar/avatar.vue'
import AuthorWrapper from '@/components/AuthorWrapper/authWrapper.vue'
import LikeAction from '@/components/LikeAction/likeAction.vue'
import { getCommentReplies } from '@/services/commentService'

const props = defineProps({
  comment: {
    type: Object,
    required: true
  },
  avatarSize: {
    type: [Number, String],
    default: 34
  },
  rootComment: {
    type: Object,
    default: null
  }
})

defineEmits(['open-author', 'reply', 'login-request'])

const expanded = ref(false)
const replies = ref([...(props.comment.replies || [])])
const repliesHasMore = ref(Boolean(props.comment.replyHasMore))
const repliesLoading = ref(false)
const repliesPage = ref(0)
const replyTotal = computed(() => Number(props.comment.replyCount || replies.value.length))

watch(() => props.comment.replies, (value) => {
	const merged = [...replies.value, ...(value || [])]
	replies.value = [...new Map(merged.map((item) => [String(item.id), item])).values()]
	if (repliesPage.value === 0) repliesHasMore.value = Boolean(props.comment.replyHasMore)
}, { deep: true })

async function toggleReplies() {
  expanded.value = !expanded.value
  if (!expanded.value || !repliesHasMore.value) return
  await loadMoreReplies(true)
}

async function loadMoreReplies(reset = false) {
  if (repliesLoading.value) return
  repliesLoading.value = true
  try {
    const page = reset ? 1 : repliesPage.value + 1
    const data = await getCommentReplies(props.comment.id, { page, pageSize: 10 })
    const merged = reset ? data.list : [...replies.value, ...data.list]
    replies.value = [...new Map(merged.map((item) => [String(item.id), item])).values()]
    repliesPage.value = data.page
    repliesHasMore.value = data.hasMore
  } catch (error) {
    showToast(error.message || '回复加载失败')
  } finally {
    repliesLoading.value = false
  }
}
</script>

<style scoped>
.comment-item {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr);
  gap: 11px;
  margin-bottom: 20px;
}

.comment-item__author {
  display: flex;
  align-items: center;
  gap: 6px;
  margin: 0 0 6px;
  color: #8a8f99;
  font-size: 13px;
}

.comment-item__content {
  overflow-wrap: anywhere;
  color: #333841;
  font-size: 13px;
  line-height: 1.7;
  word-break: break-word;
}

.comment-item__reply-to {
  color: #315b91;
}

.comment-item__image {
  width: min(170px, 70%);
  margin-top: 9px;
  overflow: hidden;
  border-radius: 8px;
  background: #f4f5f6;
}

.comment-item__image img {
  width: 100%;
  max-height: 160px;
  object-fit: cover;
}

.comment-item__footer {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 7px;
  color: #9ba1aa;
  font-size: 12px;
}

.comment-item__footer button {
  color: inherit;
}

.comment-item__expand {
  margin-top: 10px;
  color: #315b91;
  font-size: 13px;
  font-weight: 700;
}

.comment-item__replies {
  margin-top: 12px;
}

.comment-item__replies :deep(.comment-item) {
  margin-bottom: 14px;
}

.comment-item__reply-state {
  margin: 8px 0;
  color: #9ba1aa;
  font-size: 12px;
}

.comment-item__load-more {
  margin: 2px 0 10px 39px;
  color: #315b91;
  font-size: 12px;
}
</style>
