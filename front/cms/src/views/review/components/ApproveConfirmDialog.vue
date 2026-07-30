<template>
  <div class="approve-dialog" role="dialog" aria-modal="true" aria-label="确认通过并发布" @click.self="$emit('close')">
    <section class="approve-dialog__panel">
      <header>
        <span>✓</span>
        <div>
          <h2>通过并发布</h2>
          <p>审核通过后，作品会按照作者设置立即发布或进入定时发布队列。</p>
        </div>
      </header>

      <div class="approve-dialog__work">
        <small>待发布作品</small>
        <strong>{{ title || '无标题作品' }}</strong>
      </div>

      <p class="approve-dialog__hint">
        请确认封面、标题、正文和素材内容均符合平台规范。该操作会写入审核记录。
      </p>

      <footer>
        <button type="button" :disabled="submitting" @click="$emit('close')">再看看</button>
        <button class="is-approve" type="button" :disabled="submitting" @click="$emit('confirm')">
          {{ submitting ? '正在提交…' : '确认通过' }}
        </button>
      </footer>
    </section>
  </div>
</template>

<script setup>
defineProps({
  title: { type: String, default: '' },
  submitting: { type: Boolean, default: false }
})

defineEmits(['close', 'confirm'])
</script>
