<template>
  <div class="publish-settings">
    <section class="editor-card settings-card">
      <div class="editor-card__header"><h2>内容设置</h2><button class="text-button" type="button">收起⌃</button></div>
      <div class="setting-list">
        <button class="setting-line setting-line--wide is-disabled" type="button">
          <span class="setting-icon">☷</span><span class="setting-copy"><strong>添加章节</strong><small>清晰的视频结构能提升观看体验，有助于完播率</small></span><em>视频时长不足15秒，不支持添加章节</em>
        </button>
        <button class="setting-line setting-line--wide is-disabled" type="button" disabled>
          <span class="setting-icon">▱</span><span class="setting-copy"><strong>加入合集</strong><small>合集查询与创建接口接入后开放</small></span><em>暂不可用</em>
        </button>
        <div class="setting-pair">
          <button class="setting-line setting-line--compact" type="button"><span class="setting-icon">☏</span><strong>引用笔记</strong><em>›</em></button>
          <label class="setting-line setting-line--compact"><span class="setting-icon">♢</span><strong>原创声明</strong><span class="toggle-switch"><input v-model="original" type="checkbox" @change="$emit('persist')"><i /></span></label>
        </div>
      </div>
    </section>
    <section class="editor-card more-card">
      <div class="editor-card__header"><h2>更多设置</h2><button class="text-button" type="button">收起⌃</button></div>
      <div class="visibility-field">
        <button class="select-line" type="button" :aria-expanded="visibilityOpen" @click="visibilityOpen = !visibilityOpen"><span>♧&nbsp; {{ visibilityLabel }}</span><em>⌄</em></button>
        <div v-if="visibilityOpen" class="visibility-menu">
          <button v-for="option in visibilityOptions" :key="option.value" type="button" :class="{ 'is-selected': visibility === option.value }" @click="selectVisibility(option.value)"><i>{{ visibility === option.value ? '✓' : '' }}</i>{{ option.label }}</button>
        </div>
      </div>
      <label class="select-line"><span>◷&nbsp; 定时发布</span><span class="toggle-switch"><input v-model="scheduled" type="checkbox" @change="toggleScheduled"><i /></span></label>
      <label v-if="scheduled" class="scheduled-time-field">
        <span>发布时间</span>
        <input v-model="scheduledAt" type="datetime-local" :min="minimumScheduleTime" @change="$emit('persist')">
      </label>
    </section>
  </div>
</template>
<script setup>
import { computed, ref } from 'vue'
const props = defineProps({ form: { type: Object, required: true } })
const emit = defineEmits(['persist', 'open-collection', 'update-form'])
const field = (name) => computed({ get: () => props.form[name], set: (value) => emit('update-form', { [name]: value }) })
const original = field('original')
const visibility = field('visibility')
const scheduled = field('scheduled')
const scheduledAt = field('scheduledAt')
const visibilityOpen = ref(false)
const visibilityOptions = [
  { value: 'public', label: '公开可见' },
  { value: 'private', label: '仅自己可见' },
  { value: 'mutual', label: '仅互关好友可见' }
]
const minimumScheduleTime = computed(() => {
  const value = new Date(Date.now() + 5 * 60 * 1000)
  const offset = value.getTimezoneOffset() * 60 * 1000
  return new Date(value.getTime() - offset).toISOString().slice(0, 16)
})
const visibilityLabel = computed(() => visibilityOptions.find((option) => option.value === visibility.value)?.label || '公开可见')
function selectVisibility(value) {
  visibility.value = value
  visibilityOpen.value = false
  emit('persist')
}
function toggleScheduled() {
  if (!scheduled.value) {
    scheduledAt.value = ''
  }
  emit('persist')
}
</script>
