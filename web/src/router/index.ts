import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

// 获取基础路径，支持反向代理部署
const getBasePath = () => {
  // 生产环境下，从当前 URL 获取基础路径
  if (import.meta.env.PROD) {
    const path = window.location.pathname
    const segments = path.split('/').filter(Boolean)
    
    // 如果 URL 包含 index.html，移除它
    if (segments.length > 0 && segments[segments.length - 1] === 'index.html') {
      segments.pop()
    }
    
    // 如果有路径段，返回基础路径；否则返回根路径
    return segments.length > 0 ? `/${segments.join('/')}/` : '/'
  }
  
  // 开发环境使用根路径
  return '/'
}

const router = createRouter({
  history: createWebHistory(getBasePath()),
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