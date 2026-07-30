<script setup>
import { Button as VanButton } from 'vant'

defineProps({
  user: {
    type: Object,
    required: true
  }
})
</script>

<template>
  <section class="profile-header">
    <img
      class="profile-header__avatar"
      :src="user.avatar"
      :alt="user.nickname"
    >
    <div class="profile-header__info">
      <div class="profile-header__top">
        <div>
          <h1>{{ user.nickname }}</h1>
          <p>
            小红书号：{{ user.redId }} ｜ IP属地：{{ user.location }}
          </p>
        </div>
        <div class="profile-header__actions">
          <VanButton
            type="primary"
            round
          >
            关注
          </VanButton>
          <button
            type="button"
            aria-label="私信"
          >
            💬
          </button>
          <button
            type="button"
            aria-label="更多"
          >
            ···
          </button>
        </div>
      </div>
      <p class="profile-header__bio">
        {{ user.bio }}
      </p>
      <span class="profile-header__gender">♀</span>
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
  display: grid;
  grid-template-columns: 180px minmax(0, 560px);
  justify-content: center;
  gap: 72px;
  padding: 72px 0 58px;
}

.profile-header__avatar {
  width: 176px;
  height: 176px;
  border-radius: 50%;
  object-fit: cover;
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

.profile-header__top p {
  margin: 0;
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
  margin: 18px 0 12px;
  color: #30333a;
  white-space: pre-line;
  line-height: 1.6;
}

.profile-header__gender {
  display: inline-grid;
  width: 24px;
  height: 24px;
  place-items: center;
  border-radius: 50%;
  color: #ff6d8d;
  background: #fff2f5;
  font-size: 13px;
}

.profile-header__stats {
  display: flex;
  gap: 22px;
  margin-top: 24px;
  color: #6f747d;
  font-size: 14px;
}

.profile-header__stats strong {
  color: #23262d;
}

@media (max-width: 900px) {
  .profile-header {
    grid-template-columns: 88px minmax(0, 1fr);
    gap: 18px;
    padding: 24px 14px 28px;
  }

  .profile-header__avatar {
    width: 88px;
    height: 88px;
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
