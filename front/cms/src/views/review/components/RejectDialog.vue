<template>
  <div class="reject-dialog" role="dialog" aria-modal="true" aria-label="驳回作品" @click.self="$emit('close')">
    <form @submit.prevent="submit">
      <header>
        <div>
          <span>!</span>
          <div>
            <h2>驳回作品</h2>
            <p>请给创作者一个清晰、可执行的修改建议。</p>
          </div>
        </div>
        <button type="button" aria-label="关闭" @click="$emit('close')">×</button>
      </header>
      <label>
        驳回原因
        <textarea v-model.trim="reason" maxlength="1000" rows="5" placeholder="例如：封面含有不适宜内容，请更换后重新提交。" autofocus />
      </label>
      <small>{{ reason.length }}/1000</small>
      <footer>
        <button type="button" @click="$emit('close')">取消</button>
        <button class="is-danger" type="submit" :disabled="!reason || submitting">
          {{ submitting ? '正在提交…' : '确认驳回' }}
        </button>
      </footer>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'

defineProps({
  submitting: { type: Boolean, default: false }
})

const emit = defineEmits(['close', 'submit'])
const reason = ref('')

function submit() {
  if (reason.value) emit('submit', reason.value)
}
</script>
