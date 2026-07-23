<template>
  <section class="trend-card">
    <div class="trend-card__header">
      <div>
        <h2>笔记数据趋势</h2>
        <p>近 7 天内容浏览与互动变化</p>
      </div>
      <div class="trend-card__legend">
        <span><i class="is-view" /> 浏览量</span>
        <span><i class="is-like" /> 点赞量</span>
      </div>
    </div>

    <div ref="chartContainer" class="trend-chart">
      <canvas
        ref="canvas"
        aria-label="近七天笔记浏览量和点赞量趋势图"
        @mousemove="handlePointer"
        @mouseleave="hoverIndex = -1; drawChart()"
      />
      <div
        v-if="hoverIndex >= 0"
        class="trend-chart__tooltip"
        :style="{ left: `${tooltipX}px`, top: `${tooltipY}px` }"
      >
        <strong>{{ labels[hoverIndex] }}</strong>
        <span>浏览 {{ views[hoverIndex].toLocaleString() }}</span>
        <span>点赞 {{ likes[hoverIndex].toLocaleString() }}</span>
      </div>
    </div>

    <div class="trend-card__summary">
      <span>本周总浏览 <strong>89,268</strong></span>
      <span>平均互动率 <strong>12.8%</strong></span>
      <span>表现最佳 <strong>周六</strong></span>
    </div>
  </section>
</template>

<script setup>
import { nextTick, onBeforeUnmount, onMounted, ref } from 'vue'

const canvas = ref(null)
const chartContainer = ref(null)
const hoverIndex = ref(-1)
const tooltipX = ref(0)
const tooltipY = ref(0)

const labels = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
const views = [8600, 11200, 9800, 14300, 12600, 18200, 14568]
const likes = [920, 1280, 1100, 1850, 1520, 2360, 1940]
let resizeObserver = null
let points = []

function roundedLine(ctx, data, color, width, chart) {
  const maxValue = 20000
  const stepX = chart.width / (data.length - 1)
  const linePoints = data.map((value, index) => ({
    x: chart.left + stepX * index,
    y: chart.top + chart.height - (value / maxValue) * chart.height
  }))

  ctx.beginPath()
  linePoints.forEach((point, index) => {
    if (index === 0) {
      ctx.moveTo(point.x, point.y)
      return
    }

    const previous = linePoints[index - 1]
    const centerX = (previous.x + point.x) / 2
    ctx.bezierCurveTo(centerX, previous.y, centerX, point.y, point.x, point.y)
  })
  ctx.strokeStyle = color
  ctx.lineWidth = width
  ctx.lineCap = 'round'
  ctx.lineJoin = 'round'
  ctx.stroke()

  return linePoints
}

function drawChart() {
  if (!canvas.value || !chartContainer.value) {
    return
  }

  const width = chartContainer.value.clientWidth
  const height = 250
  const ratio = window.devicePixelRatio || 1
  const ctx = canvas.value.getContext('2d')

  canvas.value.width = width * ratio
  canvas.value.height = height * ratio
  canvas.value.style.width = `${width}px`
  canvas.value.style.height = `${height}px`
  ctx.scale(ratio, ratio)
  ctx.clearRect(0, 0, width, height)

  const chart = {
    left: 42,
    top: 16,
    width: Math.max(width - 58, 20),
    height: 185
  }

  ctx.font = '10px -apple-system, BlinkMacSystemFont, sans-serif'
  ctx.textAlign = 'right'
  ctx.textBaseline = 'middle'

  for (let index = 0; index <= 4; index += 1) {
    const y = chart.top + (chart.height / 4) * index
    const value = 20 - index * 5
    ctx.beginPath()
    ctx.moveTo(chart.left, y)
    ctx.lineTo(chart.left + chart.width, y)
    ctx.strokeStyle = '#eef0f3'
    ctx.lineWidth = 1
    ctx.stroke()
    ctx.fillStyle = '#a7aeba'
    ctx.fillText(`${value}k`, chart.left - 9, y)
  }

  const viewPoints = roundedLine(ctx, views, '#ef4d58', 2.5, chart)
  roundedLine(ctx, likes.map((value) => value * 5), '#9aa4b2', 1.8, chart)
  points = viewPoints

  if (hoverIndex.value >= 0) {
    const point = viewPoints[hoverIndex.value]
    ctx.beginPath()
    ctx.moveTo(point.x, chart.top)
    ctx.lineTo(point.x, chart.top + chart.height)
    ctx.strokeStyle = '#d7dae0'
    ctx.setLineDash([4, 4])
    ctx.stroke()
    ctx.setLineDash([])
    ctx.beginPath()
    ctx.arc(point.x, point.y, 4, 0, Math.PI * 2)
    ctx.fillStyle = '#fff'
    ctx.fill()
    ctx.strokeStyle = '#ef4d58'
    ctx.lineWidth = 2
    ctx.stroke()
  }

  ctx.textAlign = 'center'
  ctx.textBaseline = 'top'
  labels.forEach((label, index) => {
    const x = chart.left + (chart.width / (labels.length - 1)) * index
    ctx.fillStyle = '#98a2b3'
    ctx.fillText(label, x, chart.top + chart.height + 16)
  })
}

function handlePointer(event) {
  if (!points.length) {
    return
  }

  const rect = canvas.value.getBoundingClientRect()
  const pointerX = event.clientX - rect.left
  const nearest = points.reduce((best, point, index) => (
    Math.abs(point.x - pointerX) < Math.abs(points[best].x - pointerX) ? index : best
  ), 0)

  hoverIndex.value = nearest
  tooltipX.value = Math.min(Math.max(points[nearest].x - 45, 8), rect.width - 100)
  tooltipY.value = Math.max(points[nearest].y - 80, 4)
  drawChart()
}

onMounted(async () => {
  await nextTick()
  drawChart()
  resizeObserver = new ResizeObserver(drawChart)
  resizeObserver.observe(chartContainer.value)
})

onBeforeUnmount(() => resizeObserver?.disconnect())
</script>

<style scoped>
.trend-card {
  min-width: 0;
  border: 1px solid #ebecef;
  border-radius: 14px;
  background: #fff;
  padding: 21px 22px 18px;
}

.trend-card__header,
.trend-card__legend,
.trend-card__summary {
  display: flex;
  align-items: center;
}

.trend-card__header {
  justify-content: space-between;
  gap: 20px;
}

.trend-card__header h2 {
  margin: 0;
  color: #182230;
  font-size: 15px;
}

.trend-card__header p {
  margin: 5px 0 0;
  color: #98a2b3;
  font-size: 11px;
}

.trend-card__legend {
  gap: 14px;
  color: #7b8494;
  font-size: 10px;
}

.trend-card__legend span {
  display: flex;
  align-items: center;
  gap: 5px;
}

.trend-card__legend i {
  width: 7px;
  height: 7px;
  border-radius: 50%;
}

.trend-card__legend i.is-view {
  background: var(--color-primary);
}

.trend-card__legend i.is-like {
  background: #9aa4b2;
}

.trend-chart {
  position: relative;
  width: 100%;
  margin-top: 14px;
}

.trend-chart canvas {
  display: block;
  width: 100%;
}

.trend-chart__tooltip {
  position: absolute;
  display: grid;
  width: 100px;
  gap: 3px;
  pointer-events: none;
  border: 1px solid #e9eaee;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.96);
  box-shadow: 0 7px 18px rgba(16, 24, 40, 0.09);
  color: #667085;
  font-size: 9px;
  padding: 7px 9px;
}

.trend-chart__tooltip strong {
  color: #344054;
  font-size: 10px;
}

.trend-card__summary {
  justify-content: flex-end;
  gap: 22px;
  border-top: 1px solid #f0f1f3;
  color: #98a2b3;
  font-size: 10px;
  padding-top: 14px;
}

.trend-card__summary strong {
  margin-left: 3px;
  color: #475467;
}

@media (max-width: 560px) {
  .trend-card {
    padding: 18px 14px;
  }

  .trend-card__header {
    align-items: flex-start;
    flex-direction: column;
    gap: 10px;
  }

  .trend-card__summary {
    align-items: flex-start;
    flex-direction: column;
    gap: 7px;
  }
}
</style>
