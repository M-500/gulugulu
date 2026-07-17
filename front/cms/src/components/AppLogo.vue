<template>
  <div class="app-logo" aria-label="咕噜咕噜 CMS">
    <span class="app-logo__mark">G</span>
    <span class="app-logo__text">
      <strong>咕噜咕噜</strong>
      <small>内容管理后台</small>
    </span>
  </div>
</template>

<style scoped>
.app-logo {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  color: var(--color-text);
}

.app-logo__mark {
  display: inline-grid;
  width: 40px;
  height: 40px;
  place-items: center;
  border-radius: 10px;
  background: #ef4d58;
  color: #fff;
  font-size: 22px;
  font-weight: 800;
  line-height: 1;
}

.app-logo__text {
  display: grid;
  gap: 2px;
}

.app-logo__text strong {
  font-size: 17px;
  line-height: 1.1;
}

.app-logo__text small {
  color: var(--color-text-muted);
  font-size: 12px;
}
</style>
