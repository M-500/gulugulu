<template>
  <article class="metric-card">
    <div class="metric-card__top">
      <span class="metric-card__icon" :class="`is-${tone}`">{{ icon }}</span>
      <span class="metric-card__change">{{ change }}</span>
    </div>
    <p>{{ label }}</p>
    <div class="metric-card__value">
      <strong>{{ value }}</strong>
      <small>{{ hint }}</small>
    </div>
  </article>
</template>

<script setup>
defineProps({
  label: { type: String, required: true },
  value: { type: String, required: true },
  change: { type: String, required: true },
  hint: { type: String, required: true },
  icon: { type: String, required: true },
  tone: { type: String, default: 'rose' }
})
</script>

<style scoped>
.metric-card {
  min-width: 0;
  border: 1px solid #ebecef;
  border-radius: 13px;
  background: #fff;
  padding: 18px;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.metric-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 30px rgba(16, 24, 40, 0.06);
}

.metric-card__top,
.metric-card__value {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.metric-card__icon {
  display: grid;
  width: 34px;
  height: 34px;
  place-items: center;
  border-radius: 10px;
  font-size: 12px;
  font-weight: 800;
}

.metric-card__icon.is-rose {
  background: #fff0f1;
  color: #e1434e;
}

.metric-card__icon.is-amber {
  background: #fff6e7;
  color: #d98b19;
}

.metric-card__icon.is-blue {
  background: #edf5ff;
  color: #377ac7;
}

.metric-card__icon.is-violet {
  background: #f3efff;
  color: #7457c9;
}

.metric-card__change {
  border-radius: 12px;
  background: #ecfdf3;
  color: #039855;
  font-size: 11px;
  font-weight: 700;
  padding: 4px 7px;
}

.metric-card > p {
  margin: 15px 0 5px;
  color: #7b8494;
  font-size: 12px;
}

.metric-card__value {
  align-items: baseline;
  gap: 10px;
}

.metric-card__value strong {
  color: #182230;
  font-size: 25px;
  letter-spacing: -0.02em;
}

.metric-card__value small {
  color: #b0b6c0;
  font-size: 10px;
  white-space: nowrap;
}
</style>
