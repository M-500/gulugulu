<template>
  <section class="editor-card compose-card">
    <input v-model.trim="title" class="title-input" maxlength="50" placeholder="填写标题会有更多赞哦" @blur="$emit('persist')">
    <textarea v-model.trim="body" maxlength="1000" placeholder="输入正文描述，真诚有价值的分享予人温暖" @blur="$emit('persist')" />
    <div class="topic-row"><button v-for="topic in topics" :key="topic" type="button" @click="$emit('toggle-topic', topic)">#{{ topic }}</button><button type="button">更多⌄</button></div>
    <div class="compose-actions"><button type="button"># 话题</button><button type="button">@ 用户</button><button type="button">☺ 表情</button><span>{{ body.length }}/1000</span></div>
  </section>
</template>
<script setup>
import { computed } from 'vue'
const props = defineProps({ form: { type: Object, required: true }, topics: { type: Array, default: () => [] } })
const emit = defineEmits(['persist', 'toggle-topic', 'update-form'])
const title = computed({ get: () => props.form.title, set: (value) => emit('update-form', { title: value }) })
const body = computed({ get: () => props.form.body, set: (value) => emit('update-form', { body: value }) })
</script>
