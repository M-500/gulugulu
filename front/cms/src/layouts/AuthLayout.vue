<template>
  <div class="auth-layout">
    <header class="auth-layout__header">
      <AppLogo />
      <nav class="auth-layout__nav" aria-label="顶部导航">
        <a href="/">官网首页</a>
        <a href="/app">前台站点</a>
        <a href="mailto:1978992154@qq.com">联系我们</a>
      </nav>
    </header>

    <main class="auth-layout__main">
      <section class="auth-layout__intro" aria-label="产品简介">
        <p class="auth-layout__eyebrow">Notes · Video · Article</p>
        <h1>把内容创作、审核与发布放在一个清爽的工作台里。</h1>
        <p>面向笔记、视频和长文社区的后台管理系统，先从安全登录开始。</p>
      </section>

      <section class="auth-layout__panel" aria-label="账号入口">
        <slot />
      </section>
    </main>

    <footer class="auth-layout__footer">
      <span>© {{ currentYear }} 咕噜咕噜</span>
      <a href="https://beian.miit.gov.cn/" target="_blank" rel="noreferrer">粤ICP备XXXXXXXX号-1</a>
    </footer>
  </div>
</template>

<script setup>
import AppLogo from '@/components/AppLogo.vue'

const currentYear = new Date().getFullYear()
</script>

<style scoped>
.auth-layout {
  display: grid;
  min-height: 100vh;
  grid-template-rows: auto 1fr auto;
  overflow: hidden;
  background:
    radial-gradient(circle at 10% 18%, rgba(239, 77, 88, 0.12), transparent 28%),
    linear-gradient(135deg, #fff8f4 0%, #f7fbff 48%, #fff 100%);
}

.auth-layout__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  width: min(1180px, calc(100% - 48px));
  margin: 0 auto;
  padding: 22px 0;
}

.auth-layout__nav {
  display: flex;
  align-items: center;
  gap: 24px;
  font-size: 14px;
}

.auth-layout__nav a {
  color: var(--color-text-muted);
  text-decoration: none;
}

.auth-layout__nav a:hover {
  color: var(--color-primary);
}

.auth-layout__main {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(360px, 430px);
  align-items: start;
  gap: 72px;
  width: min(1180px, calc(100% - 48px));
  margin: 0 auto;
  padding: 70px 0 74px;
}

.auth-layout__intro {
  max-width: 560px;
  padding-top: 34px;
}

.auth-layout__eyebrow {
  margin: 0 0 18px;
  color: var(--color-primary);
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 0;
  text-transform: uppercase;
}

.auth-layout__intro h1 {
  margin: 0;
  color: var(--color-text);
  font-size: 48px;
  font-weight: 800;
  line-height: 1.15;
}

.auth-layout__intro p:last-child {
  max-width: 440px;
  margin: 20px 0 0;
  color: var(--color-text-muted);
  font-size: 17px;
  line-height: 1.8;
}

.auth-layout__panel {
  align-self: start;
}

.auth-layout__footer {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 14px;
  min-height: 58px;
  border-top: 1px solid rgba(17, 24, 39, 0.06);
  color: var(--color-text-muted);
  font-size: 13px;
}

.auth-layout__footer a {
  color: inherit;
  text-decoration: none;
}

.auth-layout__footer a:hover {
  color: var(--color-primary);
}

@media (max-width: 900px) {
  .auth-layout__main {
    grid-template-columns: 1fr;
    gap: 32px;
    padding-top: 32px;
  }

  .auth-layout__intro {
    max-width: none;
    padding-top: 0;
  }

  .auth-layout__intro h1 {
    font-size: 34px;
  }

  .auth-layout__panel {
    align-self: auto;
    margin-bottom: 0;
  }
}

@media (max-width: 640px) {
  .auth-layout__header {
    width: min(100% - 32px, 1180px);
  }

  .auth-layout__nav {
    display: none;
  }

  .auth-layout__main {
    width: min(100% - 32px, 1180px);
    padding-bottom: 42px;
  }

  .auth-layout__intro h1 {
    font-size: 28px;
  }

  .auth-layout__footer {
    flex-wrap: wrap;
    padding: 14px 16px;
  }
}
</style>
