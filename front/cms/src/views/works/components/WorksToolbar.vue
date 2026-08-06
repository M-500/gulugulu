<template>
  <header class="works-toolbar">
    <nav class="works-tabs" aria-label="作品状态筛选">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        type="button"
        :class="{ 'is-active': activeStatus === tab.key }"
        @click="$emit('select-status', tab.key)"
      >
        {{ tab.label }}<span>{{ tab.count }}</span>
      </button>
    </nav>

    <form class="works-search" role="search" @submit.prevent="$emit('search')">
      <SvgIcon
        class="works-search__icon"
        name="search"
        :size="21"
        color="currentColor"
      />
      <input
        :value="keyword"
        type="search"
        maxlength="50"
        placeholder="搜索作品标题"
        aria-label="搜索作品标题"
        @input="$emit('update:keyword', $event.target.value)"
      >
      <button type="submit">搜索</button>
    </form>
  </header>
</template>

<script setup>
defineProps({
  tabs: { type: Array, default: () => [] },
  activeStatus: { type: String, required: true },
  keyword: { type: String, default: '' }
})

defineEmits(['select-status', 'update:keyword', 'search'])
</script>
