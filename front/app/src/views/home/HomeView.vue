<script setup>
import { computed, onMounted, ref } from 'vue'

import FloatingActions from '@/components/layout/FloatingActions.vue'
import AppSidebar from '@/components/layout/AppSidebar.vue'
import MobileTabBar from '@/components/layout/MobileTabBar.vue'
import SearchHeader from '@/components/layout/SearchHeader.vue'
import WaterfallFeed from '@/components/note/WaterfallFeed.vue'
import { useFeedStore } from '@/stores/feed'

const keyword = ref('')
const feedStore = useFeedStore()

const notes = computed(() => {
  const value = keyword.value.trim().toLowerCase()
  if (!value) return feedStore.items

  return feedStore.items.filter((note) => {
    return [note.title, note.author, note.quote].some((text) => String(text || '').toLowerCase().includes(value))
  })
})

onMounted(() => {
  feedStore.fetchHomeFeed()
})
</script>

<template>
  <div class="home-shell">
    <AppSidebar />

    <main class="home-main">
      <SearchHeader v-model="keyword" />
      <div class="home-content">
        <WaterfallFeed :notes="notes" />
      </div>
    </main>

    <FloatingActions />
    <MobileTabBar />
  </div>
</template>

<style scoped>
.home-shell {
  display: flex;
  min-height: 100vh;
  background: #fff;
}

.home-main {
  min-width: 0;
  flex: 1;
}

.home-content {
  width: min(100%, var(--layout-max));
  margin: 0 auto;
  padding: 0 36px 64px;
}

@media (max-width: 900px) {
  .home-shell {
    display: block;
  }

  .home-content {
    padding: 0 12px 86px;
  }
}
</style>
