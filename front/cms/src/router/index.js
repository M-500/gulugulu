import { createRouter, createWebHistory } from 'vue-router'

import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/',
    redirect: '/auth'
  },
  {
    path: '/auth',
    name: 'Auth',
    component: () => import('@/views/auth/AuthView.vue'),
    meta: {
      guestOnly: true,
      title: '登录注册'
    }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/views/dashboard/DashboardView.vue'),
    meta: {
      requiresAuth: true,
      title: '工作台'
    }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

router.beforeEach((to) => {
  const authStore = useAuthStore()

  document.title = to.meta.title ? `${to.meta.title} - 咕噜咕噜 CMS` : '咕噜咕噜 CMS'

  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    return { name: 'Auth' }
  }

  if (to.meta.guestOnly && authStore.isLoggedIn) {
    return { name: 'Dashboard' }
  }

  return true
})

export default router
