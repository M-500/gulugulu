<template>
  <header class="review-toolbar">
    <div>
      <p class="review-toolbar__eyebrow">CONTENT MODERATION</p>
      <h1>审核中心</h1>
      <p>检查已完成媒体处理的作品，并给出明确的审核结果。</p>
    </div>

    <form class="review-toolbar__filters" @submit.prevent="$emit('search')">
      <select :value="type" aria-label="作品类型" @change="$emit('update:type', $event.target.value)">
        <option value="">全部类型</option>
        <option value="image">图片作品</option>
        <option value="video">视频作品</option>
      </select>
      <label>
        <SvgIcon
          class="review-toolbar__search-icon"
          name="search"
          :size="21"
          color="currentColor"
        />
        <input
          :value="keyword"
          type="search"
          placeholder="搜索标题、作者或邮箱"
          @input="$emit('update:keyword', $event.target.value)"
        >
      </label>
      <button type="submit">搜索</button>
    </form>
  </header>

  <nav class="review-tabs" aria-label="审核状态">
    <button
      v-for="item in tabs"
      :key="item.key"
      type="button"
      :class="{ 'is-active': item.key === activeStatus }"
      @click="$emit('select-status', item.key)"
    >
      {{ item.label }} <span>{{ item.count }}</span>
    </button>
  </nav>
</template>

<script setup>
defineProps({
  activeStatus: { type: String, required: true },
  keyword: { type: String, default: '' },
  tabs: { type: Array, default: () => [] },
  type: { type: String, default: '' }
})

defineEmits(['search', 'select-status', 'update:keyword', 'update:type'])
</script>
