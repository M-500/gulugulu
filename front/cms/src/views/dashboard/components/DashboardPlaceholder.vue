<template>
  <section class="dashboard-placeholder">
    <span>{{ menu?.icon }}</span>
    <h1>{{ menu?.label }}</h1>
    <p>{{ menu?.description }}</p>
    <small>该模块正在规划中，后续将在这里接入完整功能。</small>
  </section>
</template>

<script setup>
defineProps({
  menu: {
    type: Object,
    default: () => ({})
  }
})
</script>

<style scoped>
.dashboard-placeholder {
  display: grid;
  min-height: calc(100vh - 140px);
  place-content: center;
  justify-items: center;
  border: 1px dashed #dfe2e7;
  border-radius: 16px;
  background: #fff;
  padding: 40px;
  text-align: center;
}

.dashboard-placeholder > span {
  display: grid;
  width: 58px;
  height: 58px;
  place-items: center;
  border-radius: 16px;
  background: #fff0f1;
  color: var(--color-primary);
  font-size: 26px;
}

.dashboard-placeholder h1 {
  margin: 18px 0 7px;
  color: #182230;
  font-size: 24px;
}

.dashboard-placeholder p {
  margin: 0;
  color: #667085;
  font-size: 14px;
}

.dashboard-placeholder small {
  margin-top: 14px;
  color: #b0b6c0;
}
</style>
