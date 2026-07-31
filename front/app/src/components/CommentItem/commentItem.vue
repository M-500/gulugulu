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
          @click="replying = !replying"
        >
          回复
        </button>
      </footer>
      <form
        v-if="replying"
        class="comment-item__reply-form"
        @submit.prevent="replying = false"
      >
        <input
          v-model.trim="replyText"
          :placeholder="`回复 ${comment.author?.nickname || '用户'}`"
        >
        <button type="submit">
          发送
        </button>
      </form>
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
          :avatar-size="28"
          @open-author="$emit('open-author', $event)"
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
  }
})

defineEmits(['open-author'])

const expanded = ref(false)
const replying = ref(false)
const replyText = ref('')
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
  color: #333841;
  font-size: 13px;
  line-height: 1.7;
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

.comment-item__reply-form {
  display: flex;
  gap: 8px;
  margin-top: 10px;
}

.comment-item__reply-form input {
  min-width: 0;
  height: 32px;
  flex: 1;
  border: 0;
  border-radius: 999px;
  outline: 0;
  background: #f6f7f8;
  color: #30333a;
  padding: 0 13px;
}

.comment-item__reply-form button {
  color: var(--color-primary);
  font-weight: 700;
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
