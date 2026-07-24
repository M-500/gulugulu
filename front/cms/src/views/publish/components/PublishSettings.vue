<template>
  <div class="publish-settings">
    <section class="editor-card">
      <div class="editor-card__header"><h2>活动话题</h2><button class="text-button" type="button">更多 ›</button></div>
      <div class="activity-grid"><div v-for="activity in activities" :key="activity.title" class="activity-item"><span>{{ activity.image }}</span><strong>{{ activity.title }}</strong><small>{{ activity.action }}</small><button type="button">活动详情 ›</button></div></div>
    </section>
    <section class="editor-card settings-card">
      <div class="editor-card__header"><h2>内容设置</h2><button class="text-button" type="button">收起⌃</button></div>
      <div class="setting-list">
        <button class="setting-line is-disabled" type="button"><span>☷</span><strong>添加章节</strong><small>视频时长不足15秒，不支持添加章节</small></button>
        <button class="setting-line" type="button" @click="$emit('open-collection')"><span>▱</span><strong>{{ form.collection || '加入合集' }}</strong><small>汇集系列视频，有利于连续观看</small></button>
        <button class="setting-line" type="button"><span>☏</span><strong>引用笔记</strong><small>关联已有作品</small></button>
        <label class="setting-line"><span>♢</span><strong>原创声明</strong><input v-model="original" type="checkbox" @change="$emit('persist')"></label>
      </div>
      <h3>添加组件</h3><div class="setting-grid"><button type="button">⌖ 添加地点</button><button type="button">♙ 选择群聊</button><button type="button">◇ 标记地点或标记朋友</button><button type="button">↝ 添加路线</button></div>
    </section>
    <section class="editor-card more-card">
      <div class="editor-card__header"><h2>更多设置</h2><button class="text-button" type="button">收起⌃</button></div>
      <label class="select-line"><span>公开可见</span><select v-model="visibility" @change="$emit('persist')"><option value="public">公开可见</option><option value="fans">粉丝可见</option><option value="private">仅自己可见</option></select></label>
      <label class="select-line"><span>定时发布</span><input v-model="scheduled" type="checkbox" @change="$emit('persist')"></label>
    </section>
  </div>
</template>
<script setup>
import { computed } from 'vue'
const props = defineProps({ form: { type: Object, required: true }, activities: { type: Array, default: () => [] } })
const emit = defineEmits(['persist', 'open-collection', 'update-form'])
const field = (name) => computed({ get: () => props.form[name], set: (value) => emit('update-form', { [name]: value }) })
const original = field('original')
const visibility = field('visibility')
const scheduled = field('scheduled')
</script>
