<template>
  <section class="profile-empty">
    <div class="profile-empty__tabs">
      <span class="active">笔记 · 0</span>
      <span>专辑 · 0</span>
      <span>文件 · 0</span>
    </div>
    <div class="profile-empty__content">
      <span>☁</span>
      <p>TA 还没有收藏任何内容哦</p>
    </div>
  </section>
</template>

<style scoped>
.profile-empty {
  border-top: 1px solid #eceef2;
  padding-top: 20px;
}

.profile-empty__tabs {
  display: flex;
  gap: 26px;
  color: #777c85;
  font-size: 15px;
}

.profile-empty__tabs span {
  border-radius: 999px;
  padding: 11px 18px;
}

.profile-empty__tabs .active {
  color: #20242c;
  background: #f5f6f7;
  font-weight: 800;
}

.profile-empty__content {
  display: grid;
  min-height: 360px;
  place-content: center;
  justify-items: center;
  color: #a0a5ad;
}

.profile-empty__content span {
  display: grid;
  width: 92px;
  height: 92px;
  place-items: center;
  border: 1px solid #e8ebef;
  border-radius: 50%;
  font-size: 38px;
}

.profile-empty__content p {
  margin: 18px 0 0;
  font-size: 14px;
}
</style>
