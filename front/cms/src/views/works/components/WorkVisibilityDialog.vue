<template>
  <el-dialog v-model="visible" class="work-visibility-dialog" width="560px" :show-close="false" align-center>
    <template #header>
      <div class="visibility-dialog__header"><strong>权限设置</strong><button type="button" @click="close">×</button></div>
    </template>
    <div v-if="work" class="visibility-dialog__work">
      <img v-if="work.coverUrl" :src="work.coverUrl" :alt="work.title">
      <div><strong>{{ work.title || '无笔记标题' }}</strong><span>{{ work.type === 'video' ? '视频作品' : '图片作品' }}</span></div>
    </div>
    <label class="visibility-dialog__field">
      <span>谁可以看到这篇作品</span>
      <el-select v-model="selected" size="large">
        <el-option v-for="option in options" :key="option.value" :label="option.label" :value="option.value" />
      </el-select>
    </label>
    <template #footer>
      <el-button round @click="close">取消</el-button>
      <el-button type="primary" round :loading="saving" @click="confirm">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({ modelValue: Boolean, work: { type: Object, default: null }, saving: Boolean })
const emit = defineEmits(['update:modelValue', 'confirm'])
const visible = ref(false)
const selected = ref('public')
const options = [
  { value: 'public', label: '公开可见' },
  { value: 'private', label: '仅自己可见' },
  { value: 'mutual', label: '仅互关好友可见' }
]

watch(() => props.modelValue, (value) => {
  visible.value = value
  if (value) selected.value = props.work?.visibility || 'public'
}, { immediate: true })
watch(visible, (value) => emit('update:modelValue', value))

function close() { visible.value = false }
function confirm() { emit('confirm', selected.value) }
</script>
