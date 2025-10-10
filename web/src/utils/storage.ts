/**
 * LocalStorage 管理工具
 * 统一管理项目的本地存储，避免命名冲突
 */

const STORAGE_PREFIX = 'radius_mgnt'

/**
 * 生成带前缀的key
 * @param key 原始key
 * @returns 带前缀的key
 */
const getStorageKey = (key: string): string => {
  return `${STORAGE_PREFIX}.${key}`
}

/**
 * LocalStorage 工具类
 */
export class Storage {
  /**
   * 设置localStorage项
   * @param key 键名
   * @param value 值
   */
  static setItem(key: string, value: string): void {
    try {
      localStorage.setItem(getStorageKey(key), value)
    } catch (error) {
      console.error('Failed to set localStorage item:', error)
    }
  }

  /**
   * 获取localStorage项
   * @param key 键名
   * @param defaultValue 默认值
   * @returns 存储的值或默认值
   */
  static getItem(key: string, defaultValue: string = ''): string {
    try {
      return localStorage.getItem(getStorageKey(key)) || defaultValue
    } catch (error) {
      console.error('Failed to get localStorage item:', error)
      return defaultValue
    }
  }

  /**
   * 移除localStorage项
   * @param key 键名
   */
  static removeItem(key: string): void {
    try {
      localStorage.removeItem(getStorageKey(key))
    } catch (error) {
      console.error('Failed to remove localStorage item:', error)
    }
  }

  /**
   * 清除所有项目相关的localStorage项
   */
  static clear(): void {
    try {
      const keysToRemove: string[] = []
      for (let i = 0; i < localStorage.length; i++) {
        const key = localStorage.key(i)
        if (key && key.startsWith(`${STORAGE_PREFIX}.`)) {
          keysToRemove.push(key)
        }
      }
      keysToRemove.forEach(key => localStorage.removeItem(key))
    } catch (error) {
      console.error('Failed to clear project localStorage items:', error)
    }
  }

  /**
   * 设置JSON对象
   * @param key 键名
   * @param value 对象值
   */
  static setJSON<T>(key: string, value: T): void {
    try {
      // 如果值为null或undefined，直接移除该项
      if (value === null || value === undefined) {
        this.removeItem(key)
        return
      }
      const jsonString = JSON.stringify(value)
      this.setItem(key, jsonString)
    } catch (error) {
      console.error('Failed to set JSON to localStorage:', error, 'Value:', value)
    }
  }

  /**
   * 获取JSON对象
   * @param key 键名
   * @param defaultValue 默认值
   * @returns 解析后的对象或默认值
   */
  static getJSON<T>(key: string, defaultValue: T | null = null): T | null {
    try {
      const jsonString = this.getItem(key)
      // 检查是否为空字符串、null、undefined或"undefined"字符串
      if (!jsonString || jsonString === 'undefined' || jsonString === 'null') {
        return defaultValue
      }
      return JSON.parse(jsonString)
    } catch (error) {
      console.error('Failed to parse JSON from localStorage:', error, 'Raw value:', this.getItem(key))
      return defaultValue
    }
  }

  /**
   * 检查是否存在某个key
   * @param key 键名
   * @returns 是否存在
   */
  static hasItem(key: string): boolean {
    try {
      return localStorage.getItem(getStorageKey(key)) !== null
    } catch (error) {
      console.error('Failed to check localStorage item:', error)
      return false
    }
  }

  /**
   * 获取所有项目相关的key
   * @returns 项目相关的所有key数组
   */
  static getProjectKeys(): string[] {
    try {
      const projectKeys: string[] = []
      for (let i = 0; i < localStorage.length; i++) {
        const key = localStorage.key(i)
        if (key && key.startsWith(`${STORAGE_PREFIX}.`)) {
          // 移除前缀返回原始key
          projectKeys.push(key.substring(`${STORAGE_PREFIX}.`.length))
        }
      }
      return projectKeys
    } catch (error) {
      console.error('Failed to get project keys:', error)
      return []
    }
  }

  /**
   * 获取存储使用情况统计
   * @returns 存储统计信息
   */
  static getStorageStats(): {
    totalItems: number
    projectItems: number
    projectKeys: string[]
  } {
    try {
      const totalItems = localStorage.length
      const projectKeys = this.getProjectKeys()
      const projectItems = projectKeys.length

      return {
        totalItems,
        projectItems,
        projectKeys
      }
    } catch (error) {
      console.error('Failed to get storage stats:', error)
      return {
        totalItems: 0,
        projectItems: 0,
        projectKeys: []
      }
    }
  }
}

/**
 * 预定义的存储key常量
 * 统一管理项目中使用的所有key
 */
export const STORAGE_KEYS = {
  AUTH_TOKEN: 'auth_token',
  AUTH_USER: 'auth_user',
  USER_PREFERENCES: 'user_preferences',
  THEME: 'theme',
  LANGUAGE: 'language'
} as const

/**
 * 类型安全的存储key类型
 */
export type StorageKey = typeof STORAGE_KEYS[keyof typeof STORAGE_KEYS]

// 默认导出Storage类
export default Storage