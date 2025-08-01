import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/Login.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/',
      component: () => import('@/layouts/MainLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'Dashboard',
          component: () => import('@/views/Dashboard.vue')
        },
        {
          path: '/users',
          name: 'Users',
          component: () => import('@/views/Users.vue')
        },
        {
          path: '/profile',
          name: 'Profile',
          component: () => import('@/views/Profile.vue')
        },
        {
          path: '/auth-logs',
          name: 'AuthLogs',
          component: () => import('@/views/AuthLogs.vue')
        }
      ]
    }
  ]
})

router.beforeEach((to) => {
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth === false) {
    return true
  }
  
  if (!authStore.isAuthenticated) {
    return { name: 'Login' }
  }
  
  return true
})

export default router