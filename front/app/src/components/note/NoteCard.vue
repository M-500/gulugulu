<script setup>
import Avatar from '@/components/Avatar/avatar.vue'
import AuthorWrapper from '@/components/AuthorWrapper/authWrapper.vue'
import BaseIcon from '@/components/common/BaseIcon.vue'
import LikeAction from '@/components/LikeAction/likeAction.vue'

defineProps({
  note: {
    type: Object,
    required: true
  }
})
const emit = defineEmits(['open', 'open-author'])

</script>

<template>
  <article class="note-card">
    <div
      v-if="note.image"
      class="note-card__media"
      role="button"
      tabindex="0"
      @click="!note.isVideo && emit('open', note)"
      @keydown.enter.prevent="!note.isVideo && emit('open', note)"
    >
      <img
        :src="note.image"
        :alt="note.title"
        loading="lazy"
      >
      <span
        v-if="note.isVideo"
        class="note-card__video"
        aria-label="视频"
      >
        <BaseIcon
          name="play"
          size="14"
          fill="currentColor"
          stroke-width="0"
        />
      </span>
    </div>

    <div
      v-else
      class="note-card__quote"
      :class="`note-card__quote--${note.tone || 'paper'}`"
    >
      <span class="note-card__quote-mark">“</span>
      <p>{{ note.quote }}</p>
    </div>

    <div class="note-card__body">
      <h2>{{ note.title }}</h2>
      <div class="note-card__meta">
        <div class="note-card__author">
          <button
            type="button"
            class="note-card__author-avatar"
            :aria-label="`查看${note.author}的主页`"
            @click="emit('open-author', note)"
          >
            <Avatar
              :avatar-url="note.avatar"
              :username="note.author"
              :size="20"
            />
          </button>
          <AuthorWrapper
            :nickname="note.author"
            inline
            @click="emit('open-author', note)"
          />
        </div>
        <LikeAction
          class="note-card__like"
          :count="note.likes"
          :label="`点赞 ${note.title}`"
          icon-size="16"
        />
      </div>
    </div>
  </article>
</template>

<style scoped>
.note-card {
  display: inline-block;
  width: 100%;
  margin: 0 0 24px;
  break-inside: avoid;
  overflow: hidden;
  border-radius: 14px;
  background: var(--color-surface);
}

.note-card__media {
  position: relative;
  max-height: 430px;
  overflow: hidden;
  border-radius: 14px;
  background: var(--color-fill);
}

.note-card__media img {
  width: 100%;
  max-height: 430px;
  object-fit: cover;
  object-position: center top;
  transition: transform 0.24s ease;
}

.note-card:hover .note-card__media img {
  transform: scale(1.015);
}

.note-card__video {
  position: absolute;
  top: 10px;
  right: 10px;
  display: inline-flex;
  width: 24px;
  height: 24px;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: #fff;
  background: rgba(32, 35, 42, 0.72);
}

.note-card__quote {
  min-height: 280px;
  padding: 32px 34px;
  border: 1px solid #eee8df;
  border-radius: 10px;
  color: #46484e;
  background: #fffaf4;
}

.note-card__quote--dark {
  display: flex;
  min-height: 310px;
  flex-direction: column;
  justify-content: center;
  color: #fff;
  background: #686868;
}

.note-card__quote--note {
  min-height: 330px;
  border: 12px solid #d9fb77;
  background:
    repeating-linear-gradient(#fffdf8 0 39px, #eee6dc 40px),
    #fffdf8;
}

.note-card__quote-mark {
  display: block;
  height: 42px;
  color: rgba(128, 114, 82, 0.25);
  font-size: 72px;
  font-weight: 900;
  line-height: 0.8;
}

.note-card__quote p {
  margin: 18px 0 0;
  white-space: pre-line;
  font-size: 28px;
  font-weight: 800;
  line-height: 1.45;
  letter-spacing: 0;
}

.note-card__body {
  padding: 10px 12px 2px;
}

.note-card h2 {
  display: -webkit-box;
  margin: 0 0 10px;
  overflow: hidden;
  color: #1f2329;
  font-size: 15px;
  font-weight: 700;
  line-height: 1.45;
  letter-spacing: 0;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

.note-card__meta {
  display: flex;
  min-width: 0;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  color: var(--color-text-muted);
  font-size: 13px;
}

.note-card__author {
  display: inline-flex;
  min-width: 0;
  align-items: center;
  gap: 7px;
}

.note-card__author-avatar {
  display: inline-flex;
  border: 0;
  background: transparent;
  cursor: pointer;
  padding: 0;
}

.note-card__author :deep(.author-wrapper__name) {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.note-card__like {
  color: var(--color-text-muted);
}

@media (max-width: 900px) {
  .note-card {
    margin-bottom: 18px;
    border-radius: 12px;
  }

  .note-card__media,
  .note-card__quote {
    border-radius: 12px;
  }

  .note-card__media,
  .note-card__media img {
    max-height: 330px;
  }

  .note-card__body {
    padding: 9px 10px 1px;
  }

  .note-card h2 {
    font-size: 14px;
  }

  .note-card__meta {
    font-size: 12px;
  }

  .note-card__quote {
    min-height: 220px;
    padding: 24px 20px;
  }

  .note-card__quote p {
    font-size: 22px;
  }
}
</style>
