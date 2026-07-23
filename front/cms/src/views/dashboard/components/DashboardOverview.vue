<template>
  <div class="dashboard-overview">
    <div class="dashboard-overview__heading">
      <div>
        <p>工作台</p>
        <h1>下午好，咕噜创作者</h1>
        <span>这里是你的内容数据概览，继续保持稳定更新吧。</span>
      </div>
      <button type="button" @click="lastUpdated = '刚刚'">
        <span>↻</span> 更新数据
      </button>
    </div>

    <ProfileSummary />

    <section class="metric-grid" aria-label="核心数据">
      <MetricCard
        v-for="metric in metrics"
        :key="metric.label"
        v-bind="metric"
      />
    </section>

    <section class="dashboard-grid">
      <NoteTrendChart />
      <ContentPerformance />
    </section>

    <p class="dashboard-overview__updated">数据更新于 {{ lastUpdated }}</p>
  </div>
</template>

<script setup>
import { ref } from 'vue'

import ContentPerformance from './ContentPerformance.vue'
import MetricCard from './MetricCard.vue'
import NoteTrendChart from './NoteTrendChart.vue'
import ProfileSummary from './ProfileSummary.vue'

const lastUpdated = ref('今天 14:30')

const metrics = [
  {
    label: '粉丝数量',
    value: '12,680',
    change: '+8.2%',
    hint: '较上周',
    icon: '粉',
    tone: 'rose'
  },
  {
    label: '笔记数量',
    value: '128',
    change: '+12',
    hint: '本月新增',
    icon: '记',
    tone: 'amber'
  },
  {
    label: '关注数量',
    value: '342',
    change: '+3.1%',
    hint: '较上周',
    icon: '关',
    tone: 'blue'
  },
  {
    label: '获赞数量',
    value: '86,429',
    change: '+16.4%',
    hint: '较上周',
    icon: '赞',
    tone: 'violet'
  }
]
</script>

<style scoped>
.dashboard-overview {
  width: min(1380px, 100%);
  margin: 0 auto;
}

.dashboard-overview__heading {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 24px;
  margin-bottom: 24px;
}

.dashboard-overview__heading p {
  margin: 0 0 7px;
  color: var(--color-primary);
  font-size: 12px;
  font-weight: 800;
  letter-spacing: 0.08em;
}

.dashboard-overview__heading h1 {
  margin: 0;
  color: #182230;
  font-size: 27px;
  line-height: 1.3;
}

.dashboard-overview__heading div > span {
  display: block;
  margin-top: 7px;
  color: #98a2b3;
  font-size: 13px;
}

.dashboard-overview__heading button {
  display: flex;
  height: 38px;
  align-items: center;
  gap: 7px;
  border: 1px solid #e4e7ec;
  border-radius: 9px;
  background: #fff;
  color: #475467;
  cursor: pointer;
  font-size: 12px;
  font-weight: 600;
  padding: 0 13px;
}

.dashboard-overview__heading button:hover {
  border-color: #f2a5aa;
  color: var(--color-primary);
}

.metric-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
  margin-top: 18px;
}

.dashboard-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.75fr) minmax(270px, 0.75fr);
  gap: 18px;
  margin-top: 18px;
}

.dashboard-overview__updated {
  margin: 16px 2px 0;
  color: #b0b6c0;
  font-size: 11px;
  text-align: right;
}

@media (max-width: 1180px) {
  .metric-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .dashboard-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 620px) {
  .dashboard-overview__heading {
    align-items: flex-start;
  }

  .dashboard-overview__heading h1 {
    font-size: 23px;
  }

  .dashboard-overview__heading button {
    flex: 0 0 auto;
    padding: 0 10px;
  }

  .dashboard-overview__heading button span + * {
    display: none;
  }

  .metric-grid {
    grid-template-columns: 1fr;
  }
}
</style>
