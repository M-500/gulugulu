<script setup>
import { computed, ref } from 'vue'

import Avatar from '@/components/Avatar/avatar.vue'
import LoginDialog from '@/components/auth/LoginDialog.vue'
import BaseIcon from '@/components/common/BaseIcon.vue'
import AppSidebar from '@/components/layout/AppSidebar.vue'
import MobileTabBar from '@/components/layout/MobileTabBar.vue'
import SearchHeader from '@/components/layout/SearchHeader.vue'
import { notificationMock } from './notificationMock'

const tabs = [
  { key: 'comments', label: '评论和@' },
  { key: 'likes', label: '赞和收藏' },
  { key: 'follows', label: '新增关注' }
]

const activeTab = ref('comments')
const keyword = ref('')
const showLoginDialog = ref(false)
const likedCommentIds = ref(new Set())
const followRelations = ref(Object.fromEntries(notificationMock.follows.map((item) => [item.id, item.relation])))

const visibleNotifications = computed(() => {
  const list = notificationMock[activeTab.value] || []
  const value = keyword.value.trim().toLowerCase()
  if (!value) return list
  return list.filter((item) => [item.actor, item.action, item.content, item.quote]
    .some((text) => String(text || '').toLowerCase().includes(value)))
})

function toggleCommentLike(id) {
  const next = new Set(likedCommentIds.value)
  next.has(id) ? next.delete(id) : next.add(id)
  likedCommentIds.value = next
}

function toggleFollow(item) {
  followRelations.value[item.id] = followRelations.value[item.id] === 'mutual' ? 'none' : 'mutual'
}

function relationLabel(item) {
  const relation = followRelations.value[item.id]
  if (relation === 'mutual') return '互相关注'
  if (relation === 'follow-back') return '回关'
  return '关注'
}
</script>

<template>
  <div class="notification-shell">
    <AppSidebar @login-request="showLoginDialog = true" />

    <main class="notification-main">
      <SearchHeader v-model="keyword" />
      <section class="notification-panel">
        <div
          class="notification-tabs"
          role="tablist"
          aria-label="通知分类"
        >
          <button
            v-for="tab in tabs"
            :key="tab.key"
            type="button"
            role="tab"
            class="notification-tab"
            :class="{ active: activeTab === tab.key }"
            :aria-selected="activeTab === tab.key"
            @click="activeTab = tab.key"
          >
            {{ tab.label }}
          </button>
        </div>

        <div
          v-if="visibleNotifications.length"
          class="notification-list"
        >
          <article
            v-for="item in visibleNotifications"
            :key="item.id"
            class="notification-item"
          >
            <Avatar
              class="notification-avatar"
              :avatar-url="item.avatar"
              :username="item.actor"
              :size="48"
            />

            <div class="notification-body">
              <div class="notification-heading">
                <strong>{{ item.actor }}</strong>
                <span
                  v-if="activeTab === 'comments'"
                  class="notification-author-tag"
                >作者</span>
              </div>
              <p class="notification-meta">
                {{ item.action }} <time>{{ item.time }}</time>
              </p>
              <p
                v-if="item.content"
                class="notification-content"
              >
                {{ item.content }}
              </p>
              <p
                v-if="item.quote"
                class="notification-quote"
              >
                {{ item.quote }}
              </p>
              <div
                v-if="activeTab === 'comments' && item.content !== '该评论已删除'"
                class="notification-actions"
              >
                <button
                  type="button"
                  class="notification-action"
                >
                  <BaseIcon
                    name="message"
                    size="18"
                  />
                  回复
                </button>
                <button
                  type="button"
                  class="notification-action notification-like"
                  :class="{ active: likedCommentIds.has(item.id) }"
                  :aria-label="likedCommentIds.has(item.id) ? '取消点赞' : '点赞'"
                  @click="toggleCommentLike(item.id)"
                >
                  <BaseIcon
                    name="heart"
                    size="19"
                  />
                </button>
              </div>
              <p
                v-if="item.unsupported"
                class="notification-unsupported"
              >
                暂未支持，可在 App 内查看。
              </p>
            </div>

            <img
              v-if="item.thumbnail"
              class="notification-thumbnail"
              :src="item.thumbnail"
              alt="相关作品缩略图"
            >
            <button
              v-if="activeTab === 'follows'"
              type="button"
              class="notification-follow"
              :class="{ followed: followRelations[item.id] === 'mutual' }"
              @click="toggleFollow(item)"
            >
              {{ relationLabel(item) }}
            </button>
          </article>
        </div>
        <div
          v-else
          class="notification-empty"
        >
          没有找到相关通知
        </div>
      </section>
    </main>

    <MobileTabBar @login-request="showLoginDialog = true" />
    <LoginDialog v-model:show="showLoginDialog" />
  </div>
</template>

<style scoped>
.notification-shell {
  display: flex;
  min-height: 100vh;
  background: #fff;
}

.notification-main {
  min-width: 0;
  flex: 1;
}

.notification-panel {
  width: min(100%, 980px);
  margin: 0 auto;
  padding: 0 30px 72px;
}

.notification-tabs {
  display: flex;
  gap: 8px;
  padding: 8px 0 16px;
  border-bottom: 1px solid var(--color-border);
}

.notification-tab {
  min-height: 40px;
  padding: 0 18px;
  border-radius: 999px;
  color: #5f6269;
  font-size: 16px;
  transition: color 0.2s ease, background 0.2s ease;
}

.notification-tab:hover,
.notification-tab.active {
  color: #25272c;
  font-weight: 700;
  background: #f6f6f7;
}

.notification-item {
  position: relative;
  display: grid;
  grid-template-columns: 48px minmax(0, 1fr) auto;
  gap: 18px;
  min-height: 112px;
  padding: 24px 0;
  border-bottom: 1px solid var(--color-border);
}

.notification-avatar {
  border: 1px solid #eceef2;
}

.notification-body {
  min-width: 0;
}

.notification-heading {
  display: flex;
  align-items: center;
  gap: 8px;
  line-height: 22px;
}

.notification-heading strong {
  overflow: hidden;
  color: #2e3035;
  font-size: 16px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.notification-author-tag {
  padding: 1px 5px;
  border-radius: 4px;
  color: #999ca3;
  font-size: 11px;
  background: #f5f5f6;
}

.notification-meta,
.notification-content,
.notification-quote {
  margin: 0;
}

.notification-meta {
  margin-top: 2px;
  color: #989ba2;
  font-size: 14px;
}

.notification-meta time {
  margin-left: 6px;
}

.notification-content {
  margin-top: 9px;
  color: #3f4147;
  font-size: 15px;
  line-height: 1.55;
}

.notification-quote {
  margin-top: 9px;
  padding-left: 9px;
  border-left: 3px solid #eeeeef;
  color: #a0a2a8;
  font-size: 13px;
  line-height: 1.5;
}

.notification-actions {
  display: flex;
  gap: 10px;
  margin-top: 14px;
}

.notification-action {
  display: inline-flex;
  min-width: 42px;
  height: 38px;
  align-items: center;
  justify-content: center;
  gap: 5px;
  padding: 0 15px;
  border: 1px solid #e7e8eb;
  border-radius: 999px;
  color: #61646b;
  background: #fff;
}

.notification-like {
  min-width: 38px;
  padding: 0;
}

.notification-like.active {
  border-color: #ffd4dc;
  color: var(--color-primary);
  background: #fff4f6;
}

.notification-thumbnail {
  width: 48px;
  height: 48px;
  border-radius: 7px;
  object-fit: cover;
}

.notification-unsupported {
  margin: 12px 0 0;
  padding: 14px;
  border-radius: 8px;
  color: #999ca2;
  font-size: 13px;
  text-align: center;
  background: #f7f7f8;
}

.notification-follow {
  align-self: center;
  min-width: 96px;
  height: 40px;
  padding: 0 18px;
  border-radius: 999px;
  color: #fff;
  font-weight: 700;
  background: var(--color-primary);
}

.notification-follow.followed {
  border: 1px solid #e5e6e9;
  color: #6a6d73;
  background: #fff;
}

.notification-empty {
  padding: 100px 20px;
  color: #9a9da4;
  text-align: center;
}

@media (max-width: 900px) {
  .notification-shell {
    display: block;
  }

  .notification-panel {
    padding: 0 16px 86px;
  }

  .notification-tabs {
    position: sticky;
    top: 78px;
    z-index: 12;
    background: #fff;
  }

  .notification-tab {
    flex: 1;
    padding: 0 8px;
    font-size: 14px;
  }

  .notification-item {
    grid-template-columns: 42px minmax(0, 1fr) auto;
    gap: 12px;
    padding: 20px 0;
  }

  .notification-avatar {
    width: 42px !important;
    height: 42px !important;
    flex-basis: 42px !important;
  }

  .notification-thumbnail {
    width: 42px;
    height: 42px;
  }

  .notification-follow {
    min-width: 72px;
    height: 36px;
    padding: 0 12px;
    font-size: 13px;
  }
}
</style>
