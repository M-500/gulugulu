<script setup>
import BaseIcon from '@/components/common/BaseIcon.vue'

defineProps({
  modelValue: {
    type: String,
    default: ''
  }
})

defineEmits(['update:modelValue', 'submit'])
</script>

<template>
  <header class="search-header">
    <form
      class="search-box"
      role="search"
      @submit.prevent="$emit('submit')"
    >
      <label
        class="visually-hidden"
        for="home-search"
      >搜索笔记</label>
      <input
        id="home-search"
        :value="modelValue"
        type="search"
        placeholder="极品满天星核桃"
        autocomplete="off"
        @input="$emit('update:modelValue', $event.target.value)"
      >
      <button
        type="submit"
        class="search-box__button"
        aria-label="搜索"
      >
        <BaseIcon
          name="search"
          size="20"
        />
      </button>
    </form>
  </header>
</template>

<style scoped>
.search-header {
  position: sticky;
  top: 0;
  z-index: 15;
  display: flex;
  justify-content: center;
  padding: 28px 24px 22px;
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(18px);
}

.search-box {
  display: flex;
  width: min(100%, 900px);
  height: 54px;
  align-items: center;
  padding: 0 10px 0 22px;
  border: 1px solid #d9dce2;
  border-radius: 999px;
  background: #fff;
  box-shadow: var(--shadow-search);
}

.search-box input {
  min-width: 0;
  flex: 1;
  border: 0;
  outline: 0;
  color: var(--color-text);
  font-size: 16px;
  background: transparent;
}

.search-box input::placeholder {
  color: #babdc4;
}

.search-box__button {
  display: inline-flex;
  width: 36px;
  height: 36px;
  flex: 0 0 36px;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: #fff;
  background: #2e3036;
}

@media (max-width: 900px) {
  .search-header {
    padding: 22px 12px 16px;
  }

  .search-box {
    height: 40px;
    padding-left: 18px;
    box-shadow: none;
  }

  .search-box input {
    font-size: 14px;
  }

  .search-box__button {
    width: 30px;
    height: 30px;
    flex-basis: 30px;
  }
}
</style>
