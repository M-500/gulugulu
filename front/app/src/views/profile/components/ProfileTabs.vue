<script setup>
import { Tabs as VanTabs, Tab as VanTab } from 'vant'

defineProps({
  modelValue: {
    type: String,
    default: 'notes'
  }
})

defineEmits(['update:modelValue'])
</script>

<template>
  <div class="profile-tabs">
    <VanTabs
      :active="modelValue"
      shrink
      line-width="0"
      @update:active="$emit('update:modelValue', $event)"
    >
      <VanTab
        name="notes"
        title="笔记"
      />
      <VanTab
        name="collections"
        title="收藏"
      />
    </VanTabs>
  </div>
</template>

<style scoped>
.profile-tabs {
  display: flex;
  justify-content: center;
  margin-bottom: 18px;
}

.profile-tabs :deep(.van-tabs__wrap) {
  height: 46px;
}

.profile-tabs :deep(.van-tab) {
  flex: 0 0 auto;
  padding: 0 12px;
  color: #747982;
  font-size: 16px;
}

.profile-tabs :deep(.van-tab--active .van-tab__text) {
  border-radius: 999px;
  background: #f5f6f7;
  color: #20242c;
  font-weight: 800;
  padding: 10px 18px;
}
</style>
