<template>
  <section class="performance-card">
    <div class="performance-card__header">
      <div>
        <h2>内容表现</h2>
        <p>近 7 天数据分布</p>
      </div>
      <button type="button">查看详情</button>
    </div>

    <div class="performance-card__score">
      <div class="score-ring">
        <span>82</span>
        <small>综合分</small>
      </div>
      <div>
        <strong>表现优秀</strong>
        <p>超过 86% 的同类创作者</p>
      </div>
    </div>

    <div class="performance-list">
      <div v-for="item in performance" :key="item.label">
        <div class="performance-list__label">
          <span>{{ item.label }}</span>
          <strong>{{ item.value }}</strong>
        </div>
        <div class="performance-list__track">
          <span :style="{ width: `${item.progress}%` }" />
        </div>
      </div>
    </div>

    <div class="performance-tip">
      <span>!</span>
      <p><strong>创作建议</strong>周末 19:00 后发布，互动表现通常更好。</p>
    </div>
  </section>
</template>

<script setup>
const performance = [
  { label: '内容质量', value: '92%', progress: 92 },
  { label: '粉丝互动', value: '78%', progress: 78 },
  { label: '更新活跃', value: '84%', progress: 84 }
]
</script>

<style scoped>
.performance-card {
  border: 1px solid #ebecef;
  border-radius: 14px;
  background: #fff;
  padding: 21px 22px;
}

.performance-card__header,
.performance-card__score,
.performance-list__label {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.performance-card__header h2 {
  margin: 0;
  color: #182230;
  font-size: 15px;
}

.performance-card__header p {
  margin: 5px 0 0;
  color: #98a2b3;
  font-size: 11px;
}

.performance-card__header button {
  border: 0;
  background: transparent;
  color: var(--color-primary);
  cursor: pointer;
  font-size: 10px;
}

.performance-card__score {
  justify-content: flex-start;
  gap: 16px;
  margin: 20px 0;
}

.score-ring {
  display: grid;
  width: 70px;
  height: 70px;
  flex: 0 0 70px;
  place-content: center;
  border-radius: 50%;
  background:
    radial-gradient(circle, #fff 57%, transparent 59%),
    conic-gradient(var(--color-primary) 0 82%, #f1f2f4 82%);
  text-align: center;
}

.score-ring span {
  color: #182230;
  font-size: 20px;
  font-weight: 800;
}

.score-ring small {
  color: #98a2b3;
  font-size: 8px;
}

.performance-card__score > div:last-child strong {
  color: #344054;
  font-size: 13px;
}

.performance-card__score > div:last-child p {
  margin: 5px 0 0;
  color: #98a2b3;
  font-size: 10px;
}

.performance-list {
  display: grid;
  gap: 14px;
}

.performance-list__label {
  color: #667085;
  font-size: 10px;
}

.performance-list__label strong {
  color: #475467;
}

.performance-list__track {
  height: 5px;
  margin-top: 7px;
  overflow: hidden;
  border-radius: 4px;
  background: #f0f1f3;
}

.performance-list__track span {
  display: block;
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, #ef4d58, #ff8188);
}

.performance-tip {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  margin-top: 20px;
  border-radius: 9px;
  background: #fff8f0;
  color: #98704c;
  padding: 10px;
}

.performance-tip > span {
  display: grid;
  width: 16px;
  height: 16px;
  flex: 0 0 16px;
  place-items: center;
  border-radius: 50%;
  background: #f5c47d;
  color: #fff;
  font-size: 9px;
  font-weight: 800;
}

.performance-tip p {
  margin: 0;
  font-size: 9px;
  line-height: 1.5;
}

.performance-tip strong {
  display: block;
  color: #7d5837;
  font-size: 10px;
}
</style>
