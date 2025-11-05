import axios from 'axios'
import { Message } from '@arco-design/web-vue'

// 全局token获取函数，由main.ts设置
let getGlobalToken: () => string = () => ''

export const setTokenGetter = (getter: () => string) => {
  getGlobalToken = getter
}

// 类型定义
export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  code: number
  token: string
  expire: string
}

export interface User {
  id: number
  username: string
  email: string
  is_admin: boolean
  banned: boolean
  created_at: string
  updated_at: string
}

export interface CreateUserRequest {
  username: string
  email: string
  password: string
  is_admin?: boolean
}

export interface ChangePasswordRequest {
  old_password: string
  new_password: string
}

export interface AdminChangePasswordRequest {
  new_password: string
}

export interface UsersResponse {
  code: number
  data: {
    users: User[]
    pagination: {
      page: number
      limit: number
      total: number
    }
  }
}

export interface StatsResponse {
  code: number
  data: {
    total_users: number
    active_users: number
    banned_users: number
    auth_count: number
  }
}

export interface AuthLog {
  id: number
  username: string
  auth_type: string
  success: boolean
  ip_address: string
  user_agent: string
  created_at: number
}

export interface AuthLogsResponse {
  code: number
  data: {
    logs: AuthLog[]
    pagination: {
      page: number
      limit: number
      total: number
    }
  }
}

// 获取 API 基础路径
const getApiBasePath = () => {
  // 生产环境下，从当前 URL 构建 API 路径
  if (import.meta.env.PROD) {
    const path = window.location.pathname
    const segments = path.split('/').filter(Boolean)
    
    // 移除 index.html 如果存在
    if (segments.length > 0 && segments[segments.length - 1] === 'index.html') {
      segments.pop()
    }
    
    // 构建 API 路径
    const basePath = segments.length > 0 ? `/${segments.join('/')}/` : '/'
    return `${basePath}api/v1`.replace(/\/+/g, '/') // 清理多余的斜杠
  }
  
  // 开发环境使用相对路径
  return '/api/v1'
}

// 创建axios实例
const api = axios.create({
  baseURL: getApiBasePath(),
  timeout: 10000,
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const token = getGlobalToken()
    if (token) {
      config.headers = config.headers || {}
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    if (error.response) {
      const { status, data } = error.response
      if (status === 401) {
        // 清理认证状态，由main.ts中的认证清理函数处理
        window.location.href = '/login'
        return Promise.reject(new Error('登录已过期，请重新登录'))
      }
      const message = data?.message || '请求失败'
      Message.error(message)
      return Promise.reject(new Error(message))
    } else if (error.request) {
      Message.error('网络错误，请检查网络连接')
      return Promise.reject(new Error('网络错误'))
    } else {
      Message.error('请求配置错误')
      return Promise.reject(error)
    }
  }
)

// 认证相关API
export const login = async (credentials: LoginRequest): Promise<LoginResponse> => {
  return api.post('/auth/login', credentials)
}

export const refreshToken = async (): Promise<LoginResponse> => {
  return api.post('/auth/refresh')
}

// 用户相关API
export const getCurrentUser = async (): Promise<User> => {
  const response = await api.get('/user/profile')
  return response.data
}

export const changePassword = async (data: ChangePasswordRequest): Promise<void> => {
  return api.put('/user/change-password', data)
}

// 管理员API
export const getUsers = async (page = 1, limit = 20): Promise<UsersResponse> => {
  return api.get('/admin/users', { params: { page, limit } })
}

export const createUser = async (data: CreateUserRequest): Promise<User> => {
  const response = await api.post('/admin/users', data)
  return response.data
}

export const adminChangePassword = async (userId: number, data: AdminChangePasswordRequest): Promise<void> => {
  return api.put(`/admin/users/${userId}/password`, data)
}


export const toggleUserBan = async (userId: number): Promise<User> => {
  const response = await api.put(`/admin/users/${userId}/ban`)
  return response.data
}

export const deleteUser = async (userId: number): Promise<void> => {
  return api.delete(`/admin/users/${userId}`)
}

export const getStats = async (): Promise<StatsResponse> => {
  return api.get('/user/stats')
}

export const getAdminStats = async (): Promise<StatsResponse> => {
  return api.get('/admin/stats')
}

export const getAuthLogs = async (page = 1, limit = 20, username?: string): Promise<AuthLogsResponse> => {
  const params: any = { page, limit }
  if (username) {
    params.username = username
  }
  return api.get('/admin/auth-logs', { params })
}