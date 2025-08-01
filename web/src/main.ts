import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ArcoVue from '@arco-design/web-vue'
import App from './App.vue'
import router from './router'
import '@arco-design/web-vue/dist/arco.css'
import './style/global.css'
import { setTokenGetter } from './services/api'
import { useAuthStore } from './stores/auth'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(ArcoVue)
app.use(router)

// 设置全局token获取函数
const authStore = useAuthStore()
setTokenGetter(() => authStore.getToken())

// 应用启动时检查认证状态
const initializeAuth = async () => {
  if (authStore.token && !authStore.user) {
    try {
      await authStore.loadUserInfo()
    } catch (error) {
      console.warn('Failed to load user info on app start:', error)
      authStore.clearAuth()
    }
  }
}

// 初始化认证状态后再挂载应用
initializeAuth().then(() => {
  app.mount('#app')
})