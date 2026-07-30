import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import('@/views/home/HomeView.vue'),
    meta: {
      title: '首页'
    }
  },
  {
    path: '/users/:userId',
    name: 'user-profile',
    component: () => import('@/views/profile/UserProfileView.vue'),
    meta: {
      title: '用户主页'
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

router.beforeEach((to) => {
  document.title = `${to.meta.title || '发现'} - 咕噜咕噜`
})

export default router
