import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
import { login, getCurrentUser, type LoginRequest, type User } from '@/services/api'
import Storage, { STORAGE_KEYS } from '@/utils/storage'

export const useAuthStore = defineStore('auth', () => {
  // 从Storage工具类初始化（带项目前缀）
  const token = ref<string>(Storage.getItem(STORAGE_KEYS.AUTH_TOKEN))
  const user = ref<User | null>(Storage.getJSON<User>(STORAGE_KEYS.AUTH_USER))

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.is_admin || false)

  const setToken = (newToken: string) => {
    token.value = newToken
    Storage.setItem(STORAGE_KEYS.AUTH_TOKEN, newToken)
  }

  const setUser = (newUser: User) => {
    user.value = newUser
    Storage.setJSON(STORAGE_KEYS.AUTH_USER, newUser)
  }

  const clearAuth = () => {
    token.value = ''
    user.value = null
    Storage.removeItem(STORAGE_KEYS.AUTH_TOKEN)
    Storage.removeItem(STORAGE_KEYS.AUTH_USER)
  }

  const loginUser = async (credentials: LoginRequest) => {
    try {
      const response = await login(credentials)
      setToken(response.token)
      
      return response
    } catch (error) {
      clearAuth()
      throw error
    }
  }

  const loadUserInfo = async () => {
    if (token.value) {
      try {
        const userInfo = await getCurrentUser()
        setUser(userInfo)
      } catch (error) {
        clearAuth()
        throw error
      }
    }
  }

  // 暴露token给API拦截器使用
  const getToken = () => token.value

  const logout = () => {
    clearAuth()
  }

  return {
    token: readonly(token),
    user: readonly(user),
    isAuthenticated,
    isAdmin,
    loginUser,
    loadUserInfo,
    logout,
    clearAuth,
    getToken
  }
})