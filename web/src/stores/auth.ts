import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
import { login, getCurrentUser, type LoginRequest, type User } from '@/services/api'

export const useAuthStore = defineStore('auth', () => {
  // 从localStorage初始化
  const token = ref<string>(localStorage.getItem('auth_token') || '')
  const user = ref<User | null>((() => {
    const stored = localStorage.getItem('auth_user')
    return stored ? JSON.parse(stored) : null
  })())

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.is_admin || false)

  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('auth_token', newToken)
  }

  const setUser = (newUser: User) => {
    user.value = newUser
    localStorage.setItem('auth_user', JSON.stringify(newUser))
  }

  const clearAuth = () => {
    token.value = ''
    user.value = null
    localStorage.removeItem('auth_token')
    localStorage.removeItem('auth_user')
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