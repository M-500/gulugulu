<script setup>
import { Tabs as VanTabs, Tab as VanTab } from 'vant'

defineProps({
  modelValue: {
    type: String,
    default: 'recommend'
  },
  tabs: {
    type: Array,
    required: true
  }
})

defineEmits(['update:modelValue', 'change'])
</script>

<template>
  <div class="recommend-tabs">
    <VanTabs
      :active="modelValue"
      shrink
      line-width="0"
      color="#ff2442"
      title-active-color="#202124"
      title-inactive-color="#666a72"
      @update:active="$emit('update:modelValue', $event)"
      @change="$emit('change', $event)"
    >
      <VanTab
        v-for="item in tabs"
        :key="item.key"
        :name="item.key"
        :title="item.label"
      />
    </VanTabs>
  </div>
</template>

<style scoped>
.recommend-tabs {
  position: sticky;
  top: 104px;
  z-index: 12;
  margin: -6px 0 20px;
  background: rgba(255, 255, 255, 0.94);
  backdrop-filter: blur(18px);
}

.recommend-tabs :deep(.van-tabs__wrap) {
  height: 44px;
}

.recommend-tabs :deep(.van-tab) {
  flex: 0 0 auto;
  padding: 0 13px;
  font-size: 15px;
}

.recommend-tabs :deep(.van-tab--active .van-tab__text) {
  border-radius: 999px;
  background: #f4f5f6;
  font-weight: 800;
  padding: 8px 15px;
}

@media (max-width: 900px) {
  .recommend-tabs {
    top: 78px;
    margin-bottom: 14px;
  }
}
</style>
