<script setup>
import { Button as VanButton } from 'vant'

import Avatar from '@/components/Avatar/avatar.vue'

defineProps({
  user: {
    type: Object,
    required: true
  },
  isOwnProfile: {
    type: Boolean,
    default: false
  }
})
</script>

<template>
  <section class="profile-header">
    <Avatar class="profile-header__avatar"
            :avatar-url="user.avatar"
            :username="user.nickname"
            size="var(--profile-avatar-size)" />
    <div class="profile-header__info">
      <div class="profile-header__top">
        <div>
          <h1>{{ user.nickname }}</h1>
          <div class="profile-header__meta">
            <span>咕噜号：{{ user.redId }}</span>
            <span>IP属地：{{ user.ipAddress || '未知' }}</span>
          </div>
        </div>
        <div class="profile-header__actions">
          <VanButton v-if="!isOwnProfile"
                     type="primary"
                     round>
            关注
          </VanButton>
          <button type="button"
                  aria-label="私信">
            💬
          </button>
          <button type="button"
                  aria-label="更多">
            ···
          </button>
        </div>
      </div>
      <div class="profile-header__bio">
        <span>{{ user.bio || '这个人很神秘，还没有填写个性签名' }}</span>
        <!-- <p></p> -->
      </div>
      <div class="profile-header__stats">
        <span><strong>{{ user.following }}</strong> 关注</span>
        <span><strong>{{ user.followers }}</strong> 粉丝</span>
        <span><strong>{{ user.receives }}</strong> 获赞与收藏</span>
      </div>
    </div>
  </section>
</template>

<style scoped>
.profile-header {
  --profile-avatar-size: 176px;
  display: grid;
  grid-template-columns: 180px minmax(0, 560px);
  justify-content: center;
  gap: 72px;
  padding: 72px 0 58px;
}

.profile-header__avatar {
  font-size: 48px;
}

.profile-header__top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 24px;
}

.profile-header h1 {
  margin: 4px 0 6px;
  color: #202124;
  font-size: 26px;
}

.profile-header__meta {
  display: flex;
  flex-wrap: wrap;
  gap: 4px 14px;
  color: #a0a4ab;
  font-size: 13px;
}

.profile-header__actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.profile-header__actions :deep(.van-button) {
  width: 96px;
}

.profile-header__actions button:not(.van-button) {
  width: 42px;
  height: 42px;
  border: 1px solid #eef0f3;
  border-radius: 50%;
  background: #fff;
}

.profile-header__bio {
  margin-top: 18px;
}

.profile-header__bio span {
  color: #333;
  font-size: 14px;
}

.profile-header__bio p {
  margin: 4px 0 0;
  color: #30333a;
  white-space: pre-line;
  line-height: 1.6;
}

.profile-header__stats {
  display: flex;
  gap: 22px;
  margin-top: 28px;
  color: #6f747d;
  font-size: 14px;
}

.profile-header__stats strong {
  color: #23262d;
}

@media (max-width: 900px) {
  .profile-header {
    --profile-avatar-size: 88px;
    grid-template-columns: 88px minmax(0, 1fr);
    gap: 18px;
    padding: 24px 14px 28px;
  }

  .profile-header__avatar {
    font-size: 28px;
  }

  .profile-header__top {
    display: block;
  }

  .profile-header__actions {
    margin-top: 14px;
  }

  .profile-header h1 {
    font-size: 20px;
  }

  .profile-header__stats {
    flex-wrap: wrap;
    gap: 12px;
  }
}
</style>
