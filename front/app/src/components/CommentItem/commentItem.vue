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
      <div class="comment-item__content">
        {{ comment.content }}
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
        />
        <button
          type="button"
          @click="$emit('reply', { target: comment, root: rootComment || comment })"
        >
          回复
        </button>
      </footer>
      <button
        v-if="comment.replies?.length"
        type="button"
        class="comment-item__expand"
        @click="expanded = !expanded"
      >
        {{ expanded ? '收起回复' : `展开 ${comment.replies.length} 条回复` }}
      </button>
      <div
        v-if="expanded && comment.replies?.length"
        class="comment-item__replies"
      >
        <CommentItem
          v-for="reply in comment.replies"
          :key="reply.id"
          :comment="reply"
          :root-comment="rootComment || comment"
          :avatar-size="28"
          @open-author="$emit('open-author', $event)"
          @reply="$emit('reply', $event)"
        />
      </div>
    </div>
  </article>
</template>

<script setup>
import { ref } from 'vue'

import Avatar from '@/components/Avatar/avatar.vue'
import AuthorWrapper from '@/components/AuthorWrapper/authWrapper.vue'
import LikeAction from '@/components/LikeAction/likeAction.vue'

defineProps({
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

defineEmits(['open-author', 'reply'])

const expanded = ref(false)
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
</style>
