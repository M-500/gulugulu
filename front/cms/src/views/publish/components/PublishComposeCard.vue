<template>
  <section class="editor-card compose-card">
    <input v-model.trim="title" class="title-input" maxlength="50" placeholder="填写标题会有更多赞哦" @blur="$emit('persist')">
    <textarea v-model="body" maxlength="1000" placeholder="输入正文描述，输入 #话题 可自动生成话题标签" @blur="$emit('persist')" />
    <div v-if="selectedTopics.length" class="recognized-topics" aria-label="已识别话题">
      <span v-for="topic in selectedTopics" :key="topic">#{{ topic }}</span>
    </div>
    <div class="topic-row"><button v-for="topic in topics" :key="topic" type="button" :class="{ 'is-selected': selectedTopics.includes(topic) }" @click="$emit('toggle-topic', topic)">#{{ topic }}</button><button type="button">更多⌄</button></div>
    <div class="compose-actions"><button type="button"># 话题</button><button type="button">@ 用户</button><button type="button">☺ 表情</button><span>{{ body.length }}/1000</span></div>
  </section>
</template>
<script setup>
import { computed } from 'vue'
const props = defineProps({
  form: { type: Object, required: true },
  topics: { type: Array, default: () => [] },
  selectedTopics: { type: Array, default: () => [] }
})
const emit = defineEmits(['persist', 'toggle-topic', 'update-form', 'sync-body-topics'])
const title = computed({ get: () => props.form.title, set: (value) => emit('update-form', { title: value }) })
const body = computed({
  get: () => props.form.body,
  set: (value) => {
    emit('update-form', { body: value })
    emit('sync-body-topics', extractTopics(value))
  }
})

function extractTopics(value) {
  const matches = value.matchAll(/#([^\s#，。！？、,.!?；;：:]+)/g)
  return [...new Set(Array.from(matches, (match) => match[1]).filter(Boolean))]
}
</script>
