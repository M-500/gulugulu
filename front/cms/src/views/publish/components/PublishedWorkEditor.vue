<template>
  <section class="published-editor">
    <header class="published-editor__header">
      <div><button type="button" @click="$emit('close')">← 返回作品管理</button><h1>编辑作品</h1><p>修改已提交作品的标题和可见范围。</p></div>
      <el-tag v-if="detail" effect="light">作品 #{{ detail.workId }}</el-tag>
    </header>
    <div v-if="loading" class="published-editor__state">正在加载作品信息…</div>
    <div v-else-if="errorMessage" class="published-editor__state is-error">{{ errorMessage }}<el-button @click="load">重新加载</el-button></div>
    <div v-else-if="detail" class="published-editor__layout">
      <el-card class="published-editor__media" shadow="never">
        <img v-if="detail.coverUrl" :src="detail.coverUrl" :alt="detail.title">
        <div class="published-editor__media-copy"><strong>{{ detail.type === 'video' ? '视频作品' : `图片作品 · ${detail.assets?.length || 0} 张` }}</strong><span>媒体文件已完成处理，本次编辑不会重新上传素材。</span></div>
      </el-card>
      <el-card class="published-editor__form" shadow="never">
        <el-form label-position="top" @submit.prevent="save">
          <el-form-item label="作品标题"><el-input v-model.trim="form.title" maxlength="50" show-word-limit /></el-form-item>
          <el-form-item label="可见范围">
            <el-select v-model="form.visibility" style="width: 100%">
              <el-option label="公开可见" value="public" /><el-option label="仅自己可见" value="private" /><el-option label="仅互关好友可见" value="mutual" />
            </el-select>
          </el-form-item>
          <div class="published-editor__actions"><el-button @click="$emit('close')">取消</el-button><el-button type="primary" :loading="saving" @click="save">保存修改</el-button></div>
        </el-form>
      </el-card>
    </div>
  </section>
</template>

<script setup>
import { onMounted, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { getWork, updateCreatorWorkTitle, updateCreatorWorkVisibility } from '@/api/works'

const props = defineProps({ workId: { type: [Number, String], required: true } })
const emit = defineEmits(['close', 'saved'])
const loading = ref(false)
const saving = ref(false)
const errorMessage = ref('')
const detail = ref(null)
const form = reactive({ title: '', visibility: 'public' })

onMounted(load)
watch(() => props.workId, load)

async function load() {
  loading.value = true
  errorMessage.value = ''
  try {
    detail.value = await getWork(props.workId)
    form.title = detail.value.title || ''
    form.visibility = detail.value.visibility || 'public'
  } catch (error) { errorMessage.value = error.message || '作品信息加载失败' } finally { loading.value = false }
}

async function save() {
  if (!form.title) return ElMessage.warning('请输入作品标题')
  saving.value = true
  try {
    if (form.title !== detail.value.title) await updateCreatorWorkTitle(props.workId, form.title)
    if (form.visibility !== detail.value.visibility) await updateCreatorWorkVisibility(props.workId, form.visibility)
    ElMessage.success('作品修改成功')
    emit('saved')
  } catch (error) { ElMessage.error(error.message || '作品修改失败') } finally { saving.value = false }
}
</script>

<style scoped>
.published-editor { max-width: 1060px; margin: 0 auto; }
.published-editor__header { display: flex; align-items: flex-end; justify-content: space-between; margin-bottom: 22px; }
.published-editor__header button { border: 0; background: transparent; color: var(--color-primary); cursor: pointer; padding: 0; }
.published-editor__header h1 { margin: 10px 0 4px; color: #252d39; font-size: 26px; }
.published-editor__header p { margin: 0; color: #9aa1ab; font-size: 12px; }
.published-editor__layout { display: grid; grid-template-columns: minmax(300px, .85fr) minmax(400px, 1.15fr); gap: 20px; }
.published-editor__media, .published-editor__form { border: 1px solid #e8eaee; border-radius: 16px; }
.published-editor__media img { width: 100%; max-height: 390px; border-radius: 12px; object-fit: contain; background: #f5f6f8; }
.published-editor__media-copy { display: grid; gap: 6px; margin-top: 14px; }
.published-editor__media-copy strong { color: #363d48; font-size: 13px; }
.published-editor__media-copy span { color: #9ba1aa; font-size: 10px; }
.published-editor__actions { display: flex; justify-content: flex-end; gap: 8px; border-top: 1px solid #eff0f2; margin-top: 26px; padding-top: 18px; }
.published-editor__actions .el-button--primary { border-color: var(--color-primary); background: var(--color-primary); }
.published-editor__state { display: grid; min-height: 320px; place-content: center; gap: 14px; border-radius: 16px; background: #fff; color: #9299a3; }
.published-editor__state.is-error { color: #d84c58; }
@media (max-width: 800px) { .published-editor__layout { grid-template-columns: 1fr; } }
</style>
